/*
facturas/handlers.go — Capa HTTP para facturas del cliente.

ENDPOINT ÚNICO:
GET /api/facturas → Lista de facturas del cliente autenticado.

FLUJO FRONTEND:
1. Página "Mis Facturas" llama apiGet('/facturas')
2. Recibe [{ID_FACTURA, FECHA_FACTURA, TOTAL, ESTADO}, ...]
3. Si ESTADO='Pendiente', muestra botón "Pagar" → llama a /api/checkout/crear-sesion
4. Si ESTADO='Pagada', muestra badge verde "Pagada"
*/
package facturas

import (
	"database/sql"
	"log"
	"net/http"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewFacturaService(db)
	mux.HandleFunc("GET /api/facturas", obtenerHandler(svc))
}

// obtenerHandler → GET /api/facturas
// Retorna todas las facturas del cliente. El service resuelve la cadena de JOINs
// para filtrar por ID_CLIENTE a través de MASCOTAS → CITAS → CITAS_SERVICIOS.
func obtenerHandler(svc *FacturaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerPorCliente(r.Context(), c.IDCliente)
		if err != nil {
			log.Printf("Error facturas: %v", err)
			shared.JSONErr(w, 500, "Error al obtener facturas.")
			return
		}
		if list == nil {
			list = []Factura{}
		}
		shared.JSONOk(w, list)
	}
}
