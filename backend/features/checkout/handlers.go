/*
checkout/handlers.go — Capa HTTP para pagos con Stripe Checkout.

ENDPOINTS (2):
POST /api/checkout/crear-sesion → Crea la sesión de Stripe y retorna {url}.
POST /api/checkout/verificar    → Verifica si el pago fue completado.

SEGURIDAD:
Ambos requieren RequireAuth (cualquier usuario logueado).
La validación de propiedad (que la factura sea del cliente) se hace en service.go.

CLASIFICACIÓN DE ERRORES:
Los handlers distinguen entre errores de negocio (400) y errores del sistema (500):
- "factura no encontrada o no autorizada" → 400 (el usuario hizo algo mal).
- "esta factura ya fue pagada" → 400 (no es un error del sistema).
- Cualquier otro error (DB, Stripe API) → 500 (error real del servidor).
*/
package checkout

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"patitas-backend/shared"
)

// RegisterRoutes registra los 2 endpoints de checkout.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewCheckoutService(db)
	mux.HandleFunc("POST /api/checkout/crear-sesion", crearSesionHandler(svc))
	mux.HandleFunc("POST /api/checkout/verificar", verificarHandler(svc))
}

// crearSesionHandler recibe {id_factura}, crea la sesión Stripe, retorna {url}.
// El frontend redirige al usuario a esa URL para completar el pago.
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

		// c.IDCliente viene de la sesión, NO del body — previene pagar facturas ajenas.
		url, err := svc.CrearSesion(r.Context(), idFactura, c.IDCliente)
		if err != nil {
			errMsg := err.Error()
			// Errores de negocio → 400 (no es culpa del servidor).
			if errMsg == "factura no encontrada o no autorizada" || errMsg == "esta factura ya fue pagada" {
				shared.JSONErr(w, 400, errMsg)
				return
			}
			// Error real del sistema → 500.
			log.Printf("Error checkout crear-sesion: %v", err)
			shared.JSONErr(w, 500, "Error al crear sesión de pago.")
			return
		}

		shared.JSONOk(w, map[string]string{"url": url})
	}
}

// verificarHandler recibe {session_id} de Stripe y verifica que el pago se completó.
// Se llama automáticamente desde la página /pago-felicidades/ al cargar.
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
