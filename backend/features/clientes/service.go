package clientes

import (
	"context"
	"database/sql"
)

type ClienteService struct {
	db *sql.DB
}

func NewClienteService(db *sql.DB) *ClienteService {
	return &ClienteService{db: db}
}

type Perfil struct {
	Nombre         string `json:"nombre"`
	Apellido       string `json:"apellido"`
	Correo         string `json:"correo"`
	Telefono       string `json:"telefono"`
	Identificacion string `json:"identificacion"`
	Direccion      string `json:"direccion"`
}

func (s *ClienteService) ObtenerPerfil(ctx context.Context, idCliente int64) (*Perfil, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT c.NOMBRE, c.APELLIDO, u.CORREO, c.TELEFONO, c.IDENTIFICACION, c.DIRECCION
		 FROM CLIENTES c
		 JOIN USUARIOS u ON c.ID_CLIENTE = u.ID_CLIENTE
		 WHERE c.ID_CLIENTE = ?`,
		idCliente,
	)

	var p Perfil
	if err := row.Scan(&p.Nombre, &p.Apellido, &p.Correo, &p.Telefono, &p.Identificacion, &p.Direccion); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (s *ClienteService) Actualizar(ctx context.Context, idCliente int64, nombre, apellido, telefono string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE CLIENTES
		 SET NOMBRE = ?, APELLIDO = ?, TELEFONO = ?
		 WHERE ID_CLIENTE = ?`,
		nombre, apellido, telefono, idCliente,
	)
	return err
}
