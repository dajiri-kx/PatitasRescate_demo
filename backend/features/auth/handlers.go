package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewAuthService(db)
	mux.HandleFunc("POST /api/auth/login", loginHandler(svc))
	mux.HandleFunc("POST /api/auth/register", registerHandler(svc))
	mux.HandleFunc("POST /api/auth/logout", logoutHandler())
	mux.HandleFunc("GET /api/auth/check-session", checkSessionHandler())
}

func loginHandler(svc *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		username := strings.TrimSpace(body.Username)
		if username == "" || body.Password == "" {
			shared.JSONErr(w, 400, "Todos los campos son obligatorios.")
			return
		}

		cliente, err := svc.Login(r.Context(), username, body.Password)
		if err != nil {
			log.Printf("Error login: %v", err)
			shared.JSONErr(w, 500, "Error al iniciar sesión. Intente más tarde.")
			return
		}
		if cliente == nil {
			shared.JSONErr(w, 401, "Credenciales incorrectas.")
			return
		}

		sess := &shared.ClienteSession{
			IDCliente: cliente.IDCliente,
			Nombre:    cliente.Nombre,
			Apellido:  cliente.Apellido,
			Correo:    cliente.Correo,
			Telefono:  cliente.Telefono,
		}
		if err := shared.SaveCliente(w, r, sess); err != nil {
			log.Printf("Error sesión: %v", err)
			shared.JSONErr(w, 500, "Error interno.")
			return
		}

		shared.JSONOk(w, cliente)
	}
}

func registerHandler(svc *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Identificacion  string `json:"identificacion"`
			Nombre          string `json:"nombre"`
			PrimerApellido  string `json:"primerApellido"`
			Correo          string `json:"correo"`
			Telefono        string `json:"telefono"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirmPassword"`
			DireccionSennas string `json:"direccionSennas"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		if body.Identificacion == "" || body.Nombre == "" || body.PrimerApellido == "" ||
			body.Correo == "" || body.Telefono == "" || body.Password == "" || body.DireccionSennas == "" {
			shared.JSONErr(w, 400, "Todos los campos son obligatorios.")
			return
		}
		if body.Password != body.ConfirmPassword {
			shared.JSONErr(w, 400, "Las contraseñas no coinciden.")
			return
		}

		idCliente, err := svc.Registrar(r.Context(), body.Identificacion, body.Nombre, body.PrimerApellido,
			body.Correo, body.Telefono, body.DireccionSennas, body.Password)
		if err != nil {
			if strings.Contains(err.Error(), "ORA-20010") {
				shared.JSONErr(w, 409, "El correo electrónico ya está registrado.")
				return
			}
			log.Printf("Error registro: %v", err)
			shared.JSONErr(w, 500, "Error al registrar. Intente más tarde.")
			return
		}

		shared.JSONMsg(w, fmt.Sprintf("Cliente registrado exitosamente con ID %d.", idCliente))
	}
}

func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shared.ClearSession(w, r)
		shared.JSONMsg(w, "Sesión cerrada.")
	}
}

func checkSessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.GetCliente(r)
		if c == nil {
			shared.JSONErr(w, 401, "No autenticado.")
			return
		}
		shared.JSONOk(w, map[string]interface{}{
			"id_cliente": c.IDCliente,
			"nombre":     c.Nombre,
			"apellido":   c.Apellido,
			"correo":     c.Correo,
			"telefono":   c.Telefono,
		})
	}
}
