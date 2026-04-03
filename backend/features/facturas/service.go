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
`SELECT f.ID_FACTURA, TO_CHAR(f.FECHA_FACTURA, 'YYYY-MM-DD') AS FECHA_FACTURA, f.TOTAL,
        CASE WHEN c.FECHA_CITA < SYSDATE THEN 'Pagada' ELSE 'Pendiente' END AS ESTADO
 FROM CITAS_TABLAS.FACTURAS f
 JOIN CITAS_TABLAS.CITAS c ON f.ID_CITA = c.ID_CITA
 WHERE c.ID_MASCOTA IN (
     SELECT ID_MASCOTA FROM USUARIOS_TABLAS.MASCOTAS WHERE ID_CLIENTE = :id_cliente
 )
 ORDER BY f.FECHA_FACTURA DESC`,
sql.Named("id_cliente", idCliente),
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
