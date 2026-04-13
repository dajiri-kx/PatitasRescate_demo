/*
mascotas/service.go — Lógica de negocio para gestión de mascotas.

MODELO DE DATOS:
CLIENTES (1) → (N) MASCOTAS   (un cliente puede tener muchas mascotas)
MASCOTAS → (N) CITAS           (una mascota puede tener muchas citas)

FLUJO DE DATOS:
  - ObtenerPorCliente: Dashboard "Mis Mascotas" → JOIN con CLIENTES para
    mostrar datos de dueño. Retorna lista completa con especie, raza, edad.
  - ObtenerNombres: Formulario de agendar cita → solo ID + nombre para el
    <select> dropdown (subset ligero).
  - Agregar: Formulario "Agregar Mascota" → INSERT directo. El ID_CLIENTE
    viene de la sesión (seguridad: no se puede crear mascota para otro cliente).

SEGURIDAD:
Todas las operaciones están filtradas por idCliente (de la sesión).
Un cliente solo puede ver y crear sus propias mascotas.
*/
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

// Mascota es el struct completo para la página "Mis Mascotas".
// Incluye datos del dueño resueltos por JOIN (para la tabla de visualización).
type Mascota struct {
	IDMascota       int64  `json:"ID_MASCOTA"`
	NombreMascota   string `json:"NOMBRE_MASCOTA"`
	Especie         string `json:"ESPECIE"`
	Raza            string `json:"RAZA"`
	Meses           int    `json:"MESES"`
	NombreCliente   string `json:"NOMBRE_CLIENTE"`
	ApellidoCliente string `json:"APELLIDO_CLIENTE"`
}

// MascotaNombre es un subset ligero para dropdowns (agendar cita).
type MascotaNombre struct {
	IDMascota int64  `json:"ID_MASCOTA"`
	Nombre    string `json:"NOMBRE"`
}

// ObtenerPorCliente retorna todas las mascotas del cliente con su info completa.
// JOIN CLIENTES para incluir nombre/apellido del dueño en la tabla.
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

// ObtenerNombres retorna solo ID + nombre para poblar <select> en formularios.
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

// Agregar crea una nueva mascota asociada al cliente.
// ID_CLIENTE viene como parámetro (extraído de la sesión en el handler).
func (s *MascotaService) Agregar(ctx context.Context, nombre, especie, raza string, meses int, idCliente int64) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO MASCOTAS (NOMBRE, ESPECIE, RAZA, MESES, ID_CLIENTE)
 VALUES (?, ?, ?, ?, ?)`,
		nombre, especie, raza, meses, idCliente,
	)
	return err
}
