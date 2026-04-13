package checkout

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewCheckoutService(db)
	mux.HandleFunc("POST /api/checkout/crear-sesion", crearSesionHandler(svc))
	mux.HandleFunc("POST /api/checkout/verificar", verificarHandler(svc))
}

func crearSesionHandler(svc *CheckoutService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}

		var body struct {
			IDFactura string `json:"id_factura"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		idFactura, err := strconv.ParseInt(body.IDFactura, 10, 64)
		if err != nil {
			shared.JSONErr(w, 400, "ID de factura inválido.")
			return
		}

		url, err := svc.CrearSesion(r.Context(), idFactura, c.IDCliente)
		if err != nil {
			errMsg := err.Error()
			if errMsg == "factura no encontrada o no autorizada" || errMsg == "esta factura ya fue pagada" {
				shared.JSONErr(w, 400, errMsg)
				return
			}
			log.Printf("Error checkout crear-sesion: %v", err)
			shared.JSONErr(w, 500, "Error al crear sesión de pago.")
			return
		}

		shared.JSONOk(w, map[string]string{"url": url})
	}
}

func verificarHandler(svc *CheckoutService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}

		var body struct {
			SessionID string `json:"session_id"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		if body.SessionID == "" {
			shared.JSONErr(w, 400, "Session ID requerido.")
			return
		}

		err := svc.VerificarPago(r.Context(), body.SessionID)
		if err != nil {
			errMsg := err.Error()
			if errMsg == "sesión de pago no encontrada" || errMsg == "el pago aún no ha sido completado" {
				shared.JSONErr(w, 400, errMsg)
				return
			}
			log.Printf("Error checkout verificar: %v", err)
			shared.JSONErr(w, 500, "Error al verificar el pago.")
			return
		}

		shared.JSONMsg(w, "Pago verificado exitosamente.")
	}
}
