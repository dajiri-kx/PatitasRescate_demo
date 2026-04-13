package facturas

import (
	"context"
	"database/sql"
)

type FacturaService struct {
	db *sql.DB
}

func NewFacturaService(db *sql.DB) *FacturaService {
	return &FacturaService{db: db}
}

type Factura struct {
	IDFactura    int64   `json:"ID_FACTURA"`
	FechaFactura string  `json:"FECHA_FACTURA"`
	Total        float64 `json:"TOTAL"`
	Estado       string  `json:"ESTADO"`
}

func (s *FacturaService) ObtenerPorCliente(ctx context.Context, idCliente int64) ([]Factura, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT f.ID_FACTURA, DATE_FORMAT(f.FECHA_FACTURA, '%Y-%m-%d') AS FECHA_FACTURA, f.TOTAL, f.ESTADO
 FROM FACTURAS f
 JOIN CITAS_SERVICIOS cs ON cs.FACTURAS_ID_FACTURA = f.ID_FACTURA
 JOIN CITAS c ON cs.ID_CITA = c.ID_CITA
 WHERE c.ID_MASCOTA IN (
     SELECT ID_MASCOTA FROM MASCOTAS WHERE ID_CLIENTE = ?
 )
 GROUP BY f.ID_FACTURA
 ORDER BY f.FECHA_FACTURA DESC`,
		idCliente,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Factura
	for rows.Next() {
		var f Factura
		if err := rows.Scan(&f.IDFactura, &f.FechaFactura, &f.Total, &f.Estado); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}
