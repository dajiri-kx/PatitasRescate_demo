package mascotas

import (
"context"
"database/sql"
)

type MascotaService struct {
db *sql.DB
}

func NewMascotaService(db *sql.DB) *MascotaService {
return &MascotaService{db: db}
}

type Mascota struct {
IDMascota       int64  `json:"ID_MASCOTA"`
NombreMascota   string `json:"NOMBRE_MASCOTA"`
Especie         string `json:"ESPECIE"`
Raza            string `json:"RAZA"`
Meses           int    `json:"MESES"`
NombreCliente   string `json:"NOMBRE_CLIENTE"`
ApellidoCliente string `json:"APELLIDO_CLIENTE"`
}

type MascotaNombre struct {
IDMascota int64  `json:"ID_MASCOTA"`
Nombre    string `json:"NOMBRE"`
}

func (s *MascotaService) ObtenerPorCliente(ctx context.Context, idCliente int64) ([]Mascota, error) {
rows, err := s.db.QueryContext(ctx,
`SELECT m.ID_MASCOTA, m.NOMBRE AS NOMBRE_MASCOTA, m.ESPECIE, m.RAZA, m.MESES,
        c.NOMBRE AS NOMBRE_CLIENTE, c.APELLIDO AS APELLIDO_CLIENTE
 FROM USUARIOS_TABLAS.MASCOTAS m
 JOIN USUARIOS_TABLAS.CLIENTES c ON m.ID_CLIENTE = c.ID_CLIENTE
 WHERE m.ID_CLIENTE = :id_cliente
 ORDER BY m.ID_MASCOTA`,
sql.Named("id_cliente", idCliente),
)
if err != nil {
return nil, err
}
defer rows.Close()

var list []Mascota
for rows.Next() {
var m Mascota
if err := rows.Scan(&m.IDMascota, &m.NombreMascota, &m.Especie, &m.Raza, &m.Meses, &m.NombreCliente, &m.ApellidoCliente); err != nil {
return nil, err
}
list = append(list, m)
}
return list, rows.Err()
}

func (s *MascotaService) ObtenerNombres(ctx context.Context, idCliente int64) ([]MascotaNombre, error) {
rows, err := s.db.QueryContext(ctx,
`SELECT ID_MASCOTA, NOMBRE FROM USUARIOS_TABLAS.MASCOTAS WHERE ID_CLIENTE = :id_cliente`,
sql.Named("id_cliente", idCliente),
)
if err != nil {
return nil, err
}
defer rows.Close()

var list []MascotaNombre
for rows.Next() {
var m MascotaNombre
if err := rows.Scan(&m.IDMascota, &m.Nombre); err != nil {
return nil, err
}
list = append(list, m)
}
return list, rows.Err()
}

func (s *MascotaService) Agregar(ctx context.Context, nombre, especie, raza string, meses int, idCliente int64) error {
_, err := s.db.ExecContext(ctx,
`INSERT INTO USUARIOS_TABLAS.MASCOTAS (NOMBRE, ESPECIE, RAZA, MESES, ID_CLIENTE)
 VALUES (:nombre, :especie, :raza, :meses, :id_cliente)`,
sql.Named("nombre", nombre),
sql.Named("especie", especie),
sql.Named("raza", raza),
sql.Named("meses", meses),
sql.Named("id_cliente", idCliente),
)
return err
}
