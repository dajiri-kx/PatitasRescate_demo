/*
facturas/service.go — Lógica de negocio para facturas de clientes.

MODELO DE DATOS (cadena de JOINs para llegar a facturas de un cliente):
CLIENTES → MASCOTAS → CITAS → CITAS_SERVICIOS → FACTURAS

No existe relación directa CLIENTES → FACTURAS. La relación es indirecta:
1. Un CLIENTE tiene MASCOTAS.
2. Las MASCOTAS tienen CITAS.
3. Las CITAS están vinculadas a SERVICIOS vía CITAS_SERVICIOS (tabla puente).
4. CITAS_SERVICIOS tiene FK a FACTURAS (FACTURAS_ID_FACTURA).

Por eso el query usa subquery: WHERE c.ID_MASCOTA IN (SELECT ... WHERE ID_CLIENTE = ?)
y GROUP BY f.ID_FACTURA para evitar duplicados (una factura puede tener varios servicios).

ESTADOS DE FACTURA:
- 'Pendiente' → recién creada al agendar cita, aún no pagada
- 'Pagada'    → después de verificar pago con Stripe
*/
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

// Factura es el struct para la página "Mis Facturas" del cliente.
// El Estado ('Pendiente'/'Pagada') determina si se muestra botón de pagar.
type Factura struct {
	IDFactura    int64   `json:"ID_FACTURA"`
	FechaFactura string  `json:"FECHA_FACTURA"`
	Total        float64 `json:"TOTAL"`
	Estado       string  `json:"ESTADO"`
}

// ObtenerPorCliente busca facturas del cliente navegando la cadena:
// FACTURAS ← CITAS_SERVICIOS ← CITAS ← MASCOTAS(ID_CLIENTE=?)
// GROUP BY evita duplicados cuando una factura tiene múltiples servicios.
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
