package veterinario

import (
	"context"
	"database/sql"
)

type VetService struct {
	db *sql.DB
}

func NewVetService(db *sql.DB) *VetService {
	return &VetService{db: db}
}

// ---------- Dashboard stats ----------

type VetStats struct {
	CitasHoy         int `json:"citas_hoy"`
	CitasPendientes  int `json:"citas_pendientes"`
	CitasCompletadas int `json:"citas_completadas"`
}

func (s *VetService) GetStats(ctx context.Context, idVet int64) (*VetStats, error) {
	st := &VetStats{}
	queries := []struct {
		dest  *int
		query string
	}{
		{&st.CitasHoy, `SELECT COUNT(*) FROM CITAS WHERE ID_VETERINARIO = ? AND ESTADO = 'Activa' AND DATE(FECHA_CITA) = CURDATE()`},
		{&st.CitasPendientes, `SELECT COUNT(*) FROM CITAS WHERE ID_VETERINARIO = ? AND ESTADO = 'Activa'`},
		{&st.CitasCompletadas, `SELECT COUNT(*) FROM CITAS WHERE ID_VETERINARIO = ? AND ESTADO = 'Completada'`},
	}
	for _, q := range queries {
		if err := s.db.QueryRowContext(ctx, q.query, idVet).Scan(q.dest); err != nil {
			return nil, err
		}
	}
	return st, nil
}

// ---------- Mis citas (agenda) ----------

type VetCita struct {
	ID        int64   `json:"ID_CITA"`
	Fecha     string  `json:"FECHA_CITA"`
	Estado    string  `json:"ESTADO"`
	Mascota   string  `json:"MASCOTA"`
	Especie   string  `json:"ESPECIE"`
	Raza      string  `json:"RAZA"`
	Cliente   string  `json:"CLIENTE"`
	Telefono  string  `json:"TELEFONO_CLIENTE"`
	Servicios string  `json:"SERVICIOS"`
	Total     float64 `json:"TOTAL"`
}

func (s *VetService) ListCitas(ctx context.Context, idVet int64) ([]VetCita, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT c.ID_CITA,
		        DATE_FORMAT(c.FECHA_CITA, '%Y-%m-%d %H:%i') AS FECHA_CITA,
		        c.ESTADO,
		        m.NOMBRE AS MASCOTA,
		        m.ESPECIE,
		        IFNULL(m.RAZA, '') AS RAZA,
		        CONCAT(cl.NOMBRE, ' ', cl.APELLIDO) AS CLIENTE,
		        cl.TELEFONO,
		        IFNULL(GROUP_CONCAT(sv.NOMBRE_SERVICIO SEPARATOR ', '), '') AS SERVICIOS,
		        IFNULL(SUM(sv.PRECIO), 0) AS TOTAL
		 FROM CITAS c
		 JOIN MASCOTAS m ON c.ID_MASCOTA = m.ID_MASCOTA
		 JOIN CLIENTES cl ON m.ID_CLIENTE = cl.ID_CLIENTE
		 LEFT JOIN CITAS_SERVICIOS cs ON cs.ID_CITA = c.ID_CITA
		 LEFT JOIN SERVICIOS sv ON cs.ID_SERVICIO = sv.ID_SERVICIO
		 WHERE c.ID_VETERINARIO = ?
		 GROUP BY c.ID_CITA
		 ORDER BY c.FECHA_CITA DESC`, idVet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []VetCita
	for rows.Next() {
		var r VetCita
		if err := rows.Scan(&r.ID, &r.Fecha, &r.Estado, &r.Mascota, &r.Especie, &r.Raza, &r.Cliente, &r.Telefono, &r.Servicios, &r.Total); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

// ---------- Cambiar estado de cita ----------

func (s *VetService) UpdateEstadoCita(ctx context.Context, idVet, idCita int64, estado string) error {
	// Solo permite cambiar sus propias citas
	_, err := s.db.ExecContext(ctx,
		`UPDATE CITAS SET ESTADO = ? WHERE ID_CITA = ? AND ID_VETERINARIO = ?`,
		estado, idCita, idVet)
	return err
}
