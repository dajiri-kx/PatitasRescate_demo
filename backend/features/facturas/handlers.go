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
