package contacto

import (
	"net/http"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/contacto", enviarHandler())
}

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
