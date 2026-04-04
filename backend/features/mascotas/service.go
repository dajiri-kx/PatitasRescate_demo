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
 FROM MASCOTAS m
 JOIN CLIENTES c ON m.ID_CLIENTE = c.ID_CLIENTE
 WHERE m.ID_CLIENTE = ?
 ORDER BY m.ID_MASCOTA`,
		idCliente,
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
		`SELECT ID_MASCOTA, NOMBRE FROM MASCOTAS WHERE ID_CLIENTE = ?`,
		idCliente,
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
		`INSERT INTO MASCOTAS (NOMBRE, ESPECIE, RAZA, MESES, ID_CLIENTE)
 VALUES (?, ?, ?, ?, ?)`,
		nombre, especie, raza, meses, idCliente,
	)
	return err
}
