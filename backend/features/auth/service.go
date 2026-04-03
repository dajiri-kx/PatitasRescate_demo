package auth

import (
"context"
"database/sql"

"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
return &AuthService{db: db}
}

type ClienteData struct {
IDCliente int64  `json:"id_cliente"`
Nombre    string `json:"nombre"`
Apellido  string `json:"apellido"`
Correo    string `json:"correo"`
Telefono  string `json:"telefono"`
}

func (s *AuthService) Login(ctx context.Context, correo, password string) (*ClienteData, error) {
row := s.db.QueryRowContext(ctx,
`SELECT c.ID_CLIENTE, c.NOMBRE, c.APELLIDO, c.TELEFONO, u.CONTRASENA
 FROM USUARIOS_TABLAS.USUARIOS u
 JOIN USUARIOS_TABLAS.CLIENTES c ON u.ID_CLIENTE = c.ID_CLIENTE
 WHERE u.CORREO = :correo`,
sql.Named("correo", correo),
)

var idCliente int64
var nombre, apellido, telefono, hash string
if err := row.Scan(&idCliente, &nombre, &apellido, &telefono, &hash); err != nil {
if err == sql.ErrNoRows {
return nil, nil
}
return nil, err
}

if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
return nil, nil
}

return &ClienteData{
IDCliente: idCliente,
Nombre:    nombre,
Apellido:  apellido,
Correo:    correo,
Telefono:  telefono,
}, nil
}

func (s *AuthService) Registrar(ctx context.Context, identificacion, nombre, apellido, correo, telefono, direccion, password string) (int64, error) {
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
return 0, err
}

var idCliente int64
_, err = s.db.ExecContext(ctx,
`BEGIN registrarCliente(:identificacion, :nombre, :apellido, :correo, :telefono, :direccion, :password, :idCliente); END;`,
sql.Named("identificacion", identificacion),
sql.Named("nombre", nombre),
sql.Named("apellido", apellido),
sql.Named("correo", correo),
sql.Named("telefono", telefono),
sql.Named("direccion", direccion),
sql.Named("password", string(hash)),
sql.Named("idCliente", sql.Out{Dest: &idCliente}),
)
if err != nil {
return 0, err
}

return idCliente, nil
}
