package auth

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

type ClienteData struct {
	IDCliente     int64  `json:"id_cliente"`
	IDVeterinario int64  `json:"id_veterinario"`
	Nombre        string `json:"nombre"`
	Apellido      string `json:"apellido"`
	Correo        string `json:"correo"`
	Telefono      string `json:"telefono"`
	Rol           int    `json:"rol"`
}

func (s *AuthService) Login(ctx context.Context, correo, password string) (*ClienteData, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT c.ID_CLIENTE, c.NOMBRE, c.APELLIDO, c.TELEFONO, u.CONTRASENA, u.ROL, IFNULL(u.ID_VETERINARIO, 0)
		 FROM USUARIOS u
		 JOIN CLIENTES c ON u.ID_CLIENTE = c.ID_CLIENTE
		 WHERE u.CORREO = ?`,
		correo,
	)

	var idCliente, idVet int64
	var nombre, apellido, telefono, hash string
	var rol int
	if err := row.Scan(&idCliente, &nombre, &apellido, &telefono, &hash, &rol, &idVet); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return nil, nil
	}

	return &ClienteData{
		IDCliente:     idCliente,
		IDVeterinario: idVet,
		Nombre:        nombre,
		Apellido:      apellido,
		Correo:        correo,
		Telefono:      telefono,
		Rol:           rol,
	}, nil
}

func (s *AuthService) Registrar(ctx context.Context, identificacion, nombre, apellido, correo, telefono, direccion, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Verificar cédula duplicada
	var count int
	err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM CLIENTES WHERE DIDENTIDAD_CLIENTE = ?`, identificacion).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("ORA-20010: La cédula o documento de identidad ya está registrada")
	}

	// Verificar correo duplicado
	err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM USUARIOS WHERE CORREO = ?`, correo).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("ORA-20010: El correo electrónico ya está registrado")
	}

	// Insertar cliente
	result, err := tx.ExecContext(ctx,
		`INSERT INTO CLIENTES (DIDENTIDAD_CLIENTE, NOMBRE, APELLIDO, EMAIL, TELEFONO, DIRECCION, FECHA_REGISTRO)
		 VALUES (?, ?, ?, ?, ?, ?, NOW())`,
		identificacion, nombre, apellido, correo, telefono, direccion,
	)
	if err != nil {
		return 0, err
	}

	idCliente, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insertar usuario
	_, err = tx.ExecContext(ctx,
		`INSERT INTO USUARIOS (ID_CLIENTE, CORREO, CONTRASENA) VALUES (?, ?, ?)`,
		idCliente, correo, string(hash),
	)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return idCliente, nil
}
