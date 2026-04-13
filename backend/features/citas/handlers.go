/*
citas/handlers.go — Capa HTTP para gestión de citas veterinarias.

ENDPOINTS Y FLUJO DE DATOS:
GET  /api/citas             → Todas las citas del cliente (mis-citas page)
GET  /api/citas/activas     → Solo citas activas (formulario cancelar)
GET  /api/citas/veterinarios → Lista de vets (dropdown agendar)
GET  /api/citas/servicios   → Servicios, ?categoria=X opcional (cards agendar)
POST /api/citas/agendar     → Crear cita completa (la más compleja)
POST /api/citas/cancelar    → Eliminar cita activa

PATRÓN DE HANDLERS:
1. shared.RequireAuth → verifica sesión, obtiene c.IDCliente
2. DecodeBody → parsea JSON del request
3. svc.XXX() → llama a la lógica de negocio (service.go)
4. JSONOk / JSONErr → responde al frontend

MANEJO DE ERRORES ORA-200xx:
Los errores de negocio del service (validaciones) usan prefijo "ORA-200xx".
El handler detecta este prefijo y responde 400 (error del usuario) en vez de
500 (error del servidor). El mensaje se extrae limpio para mostrar al usuario.

NOTA DE SEGURIDAD:
El IDCliente viene de la SESIÓN (cookie), no del body del request.
Esto evita que un usuario pueda enviar un ID de cliente ajeno.
*/
package citas

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"patitas-backend/shared"
)

// RegisterRoutes registra las 6 rutas de citas en el router principal.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewCitaService(db)
	mux.HandleFunc("GET /api/citas", obtenerHandler(svc))
	mux.HandleFunc("GET /api/citas/activas", obtenerActivasHandler(svc))
	mux.HandleFunc("GET /api/citas/veterinarios", obtenerVeterinariosHandler(svc))
	mux.HandleFunc("GET /api/citas/servicios", obtenerServiciosHandler(svc))
	mux.HandleFunc("POST /api/citas/agendar", agendarHandler(svc))
	mux.HandleFunc("POST /api/citas/cancelar", cancelarHandler(svc))
}

// obtenerHandler → GET /api/citas
// Flujo: sesión → c.IDCliente → ObtenerPorCliente → lista de citas con nombre
// de mascota y veterinario resueltos por JOINs en el service.
func obtenerHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerPorCliente(r.Context(), c.IDCliente)
		if err != nil {
			log.Printf("Error citas: %v", err)
			shared.JSONErr(w, 500, "Error al obtener citas.")
			return
		}
		// Retornar array vacío en vez de null para que el frontend no falle
		if list == nil {
			list = []Cita{}
		}
		shared.JSONOk(w, list)
	}
}

// obtenerActivasHandler → GET /api/citas/activas
// Usado por el formulario de cancelar-cita para mostrar solo citas cancelables.
func obtenerActivasHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerActivas(r.Context(), c.IDCliente)
		if err != nil {
			log.Printf("Error citas activas: %v", err)
			shared.JSONErr(w, 500, "Error al obtener citas activas.")
			return
		}
		if list == nil {
			list = []CitaActiva{}
		}
		shared.JSONOk(w, list)
	}
}

// obtenerVeterinariosHandler → GET /api/citas/veterinarios
// Carga el dropdown de veterinarios en el formulario de agendar.
func obtenerVeterinariosHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerVeterinarios(r.Context())
		if err != nil {
			log.Printf("Error veterinarios: %v", err)
			shared.JSONErr(w, 500, "Error al obtener veterinarios.")
			return
		}
		if list == nil {
			list = []Veterinario{}
		}
		shared.JSONOk(w, list)
	}
}

// obtenerServiciosHandler → GET /api/citas/servicios?categoria=X
// Si ?categoria está presente, filtra por categoría; si no, devuelve todos.
// El frontend primero muestra categorías, al seleccionar una recarga servicios.
func obtenerServiciosHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		categoria := r.URL.Query().Get("categoria")
		list, err := svc.ObtenerServicios(r.Context(), categoria)
		if err != nil {
			log.Printf("Error servicios: %v", err)
			shared.JSONErr(w, 500, "Error al obtener servicios.")
			return
		}
		if list == nil {
			list = []Servicio{}
		}
		shared.JSONOk(w, list)
	}
}

// agendarHandler → POST /api/citas/agendar
// Esta es la operación más compleja del sistema.
//
// FLUJO FRONTEND → BACKEND → FRONTEND:
// 1. Frontend envía: {id_mascota, fecha, hora, servicio:["1","3","5"], veterinario}
// 2. Handler une fecha+hora ("2025-06-15" + "10:00" → "2025-06-15 10:00")
// 3. Handler une servicios (["1","3","5"] → "1,3,5")
// 4. svc.Agendar() ejecuta la transacción completa (ver service.go)
// 5. Retorna {message, id_factura} → el frontend usa id_factura para ir a pagar
//
// El idCliente viene de c.IDCliente (sesión), NO del body → seguridad.
func agendarHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}

		var body struct {
			IDMascota   string   `json:"id_mascota"`
			Fecha       string   `json:"fecha"`       // "2025-06-15"
			Hora        string   `json:"hora"`        // "10:00"
			Servicio    []string `json:"servicio"`    // ["1", "3", "5"]
			Veterinario string   `json:"veterinario"` // "2" (ID como string)
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		// Convertir strings a int64 (el frontend envía IDs como strings)
		idMascota, err := strconv.ParseInt(body.IDMascota, 10, 64)
		if err != nil {
			shared.JSONErr(w, 400, "Mascota inválida.")
			return
		}
		idVet, err := strconv.ParseInt(body.Veterinario, 10, 64)
		if err != nil {
			shared.JSONErr(w, 400, "Veterinario inválido.")
			return
		}

		// Preparar datos para el service: concatenar fecha+hora y servicios
		fechaCita := body.Fecha + " " + body.Hora         // "2025-06-15 10:00"
		serviciosList := strings.Join(body.Servicio, ",") // "1,3,5"

		idFactura, err := svc.Agendar(r.Context(), c.IDCliente, idMascota, idVet, fechaCita, serviciosList)
		if err != nil {
			errMsg := err.Error()
			// Errores ORA-200xx = errores de validación de negocio → 400
			if strings.Contains(errMsg, "ORA-200") {
				// Extraer el mensaje legible después del código "ORA-200xx: "
				if idx := strings.Index(errMsg, ": "); idx != -1 {
					shared.JSONErr(w, 400, errMsg[idx+2:])
				} else {
					shared.JSONErr(w, 400, errMsg)
				}
				return
			}
			// Errores no-ORA = errores del sistema → 500
			log.Printf("Error agendar: %v", err)
			shared.JSONErr(w, 500, "Error al agendar la cita.")
			return
		}

		// Retornar id_factura — el frontend lo usa para llamar a
		// POST /api/checkout/crear-sesion y redirigir a Stripe.
		shared.JSONOk(w, map[string]interface{}{
			"message":    "Cita agendada exitosamente.",
			"id_factura": idFactura,
		})
	}
}

// cancelarHandler → POST /api/citas/cancelar
// Elimina una cita activa. Verifica que la cita pertenece al cliente.
func cancelarHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}

		var body struct {
			IDCita string `json:"id_cita"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		idCita, err := strconv.ParseInt(body.IDCita, 10, 64)
		if err != nil || idCita <= 0 {
			shared.JSONErr(w, 400, "ID de cita inválido.")
			return
		}

		// svc.Cancelar verifica internamente que la mascota de la cita sea del cliente
		ok, err := svc.Cancelar(r.Context(), idCita, c.IDCliente)
		if err != nil {
			log.Printf("Error cancelar: %v", err)
			shared.JSONErr(w, 500, "Error al cancelar la cita.")
			return
		}
		if !ok {
			shared.JSONErr(w, 404, "Cita no encontrada o no autorizada.")
			return
		}

		shared.JSONMsg(w, "Cita cancelada exitosamente.")
	}
}
