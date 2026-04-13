/*
contacto/handlers.go — Formulario de contacto público.

ENDPOINT:
POST /api/contacto → Recibe mensaje de contacto del visitante.

NOTA: Este endpoint NO guarda en base de datos. Es un stub/placeholder
que valida los campos y responde con un mensaje de éxito.
En producción se podría integrar con un servicio de email (SendGrid, etc.)
o guardar en una tabla MENSAJES_CONTACTO.

NO requiere autenticación — es accesible para visitantes no registrados.
Es el único feature que no recibe *sql.DB (no usa base de datos).
*/
package contacto

import (
	"net/http"

	"patitas-backend/shared"
)

// RegisterRoutes registra la única ruta de contacto.
// A diferencia de otros features, NO recibe *sql.DB porque no usa base de datos.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/contacto", enviarHandler())
}

// enviarHandler → POST /api/contacto
// Valida campos requeridos y responde con mensaje de éxito (sin persistencia).
func enviarHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Nombre   string `json:"nombre"`
			Email    string `json:"email"`
			Telefono string `json:"telefono"`
			Mensaje  string `json:"mensaje"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		if body.Nombre == "" || body.Email == "" || body.Telefono == "" || body.Mensaje == "" {
			shared.JSONErr(w, 400, "Todos los campos son obligatorios.")
			return
		}

		shared.JSONMsg(w, "Mensaje recibido exitosamente. Nos pondremos en contacto pronto.")
	}
}
