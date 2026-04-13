/*
auth/handlers.go — Capa HTTP para autenticación.

FLUJO COMPLETO DE LOGIN (frontend → backend → frontend):
 1. Frontend: login/index.html → POST /api/auth/login con {username, password}
    (usa apiPost de api.js, que incluye credentials:'include' para cookies).
 2. Backend: loginHandler() decodifica JSON → llama svc.Login() → obtiene ClienteData.
 3. Backend: Mapea ClienteData → ClienteSession → SaveCliente() guarda en cookie.
 4. Backend: Responde JSON {ok:true, data: {id_cliente, nombre, rol, ...}}.
 5. Frontend: Recibe el JSON → Auth.save(data) guarda en localStorage.
 6. Frontend: Redirige según rol:
    - Rol 0 (Admin)  → /frontend/features/admin/
    - Rol 2 (Vet)    → /frontend/features/veterinario/
    - Rol 1 (Cliente) → /frontend/features/dashboard/

DOBLE ESTADO (cookie + localStorage):
  - Cookie (HttpOnly): la usa el backend para autenticar cada request. El frontend
    NO puede leerla — viaja automáticamente con cada fetch.
  - localStorage: la usa el frontend para saber si está logueado, mostrar el nombre
    en el navbar, y decidir qué links mostrar (sin hacer request al servidor).

VALIDACIÓN BACKEND (registerHandler):
Las validaciones se duplican en frontend y backend como buena práctica.
El backend es la última línea de defensa — nunca confiar solo en validación JS.
Regex: 9 dígitos para cédula CR, 8 dígitos para teléfono, complejidad de contraseña.
*/
package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"patitas-backend/shared"
)

// Expresiones regulares para validación en el servidor (duplicadas del frontend).
var (
	reDigits9 = regexp.MustCompile(`^\d{9}$`)      // Cédula costarricense: 9 dígitos
	reDigits8 = regexp.MustCompile(`^\d{8}$`)      // Teléfono CR: 8 dígitos
	rePwUpper = regexp.MustCompile(`[A-Z]`)        // Al menos una mayúscula
	rePwLower = regexp.MustCompile(`[a-z]`)        // Al menos una minúscula
	rePwDigit = regexp.MustCompile(`\d`)           // Al menos un dígito
	rePwSpec  = regexp.MustCompile(`[^A-Za-z0-9]`) // Al menos un carácter especial
)

// RegisterRoutes registra las 4 rutas de autenticación en el router.
// Go 1.22+ permite "METHOD /path" en HandleFunc para routing basado en método HTTP.
// Cada handler recibe svc como closure — así tiene acceso al servicio sin globales.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewAuthService(db)
	mux.HandleFunc("POST /api/auth/login", loginHandler(svc))
	mux.HandleFunc("POST /api/auth/register", registerHandler(svc))
	mux.HandleFunc("POST /api/auth/logout", logoutHandler())
	mux.HandleFunc("GET /api/auth/check-session", checkSessionHandler())
}

// loginHandler maneja POST /api/auth/login.
// Flujo: JSON body → Login service → crear sesión cookie → responder con datos usuario.
func loginHandler(svc *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string `json:"username"` // En realidad es el correo electrónico
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

		// svc.Login retorna: *ClienteData (éxito), nil (credenciales incorrectas), o error
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

		// Mapear ClienteData (del service) → ClienteSession (para la cookie).
		// La sesión se serializa con gob y se firma con HMAC en la cookie.
		sess := &shared.ClienteSession{
			IDCliente:     cliente.IDCliente,
			IDVeterinario: cliente.IDVeterinario,
			Nombre:        cliente.Nombre,
			Apellido:      cliente.Apellido,
			Correo:        cliente.Correo,
			Telefono:      cliente.Telefono,
			Rol:           cliente.Rol,
		}
		if err := shared.SaveCliente(w, r, sess); err != nil {
			log.Printf("Error sesión: %v", err)
			shared.JSONErr(w, 500, "Error interno.")
			return
		}

		// Responder con los datos del cliente — el frontend los guarda en localStorage
		// para uso inmediato (navbar, redirección por rol) sin necesidad de otro request.
		shared.JSONOk(w, cliente)
	}
}

// registerHandler maneja POST /api/auth/register.
// Flujo: validar campos → svc.Registrar (tx) → responder con mensaje de éxito.
// NO crea sesión automáticamente — el usuario debe hacer login después.
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

		// Validación de campos obligatorios
		if body.Identificacion == "" || body.Nombre == "" || body.PrimerApellido == "" ||
			body.Correo == "" || body.Telefono == "" || body.Password == "" || body.DireccionSennas == "" {
			shared.JSONErr(w, 400, "Todos los campos son obligatorios.")
			return
		}
		if body.Password != body.ConfirmPassword {
			shared.JSONErr(w, 400, "Las contraseñas no coinciden.")
			return
		}

		// Validaciones con regex (mismas reglas que el frontend)
		if !reDigits9.MatchString(body.Identificacion) {
			shared.JSONErr(w, 400, "La identificación debe ser exactamente 9 dígitos numéricos.")
			return
		}
		if !strings.Contains(body.Correo, "@") {
			shared.JSONErr(w, 400, "El correo debe contener un @.")
			return
		}
		if !reDigits8.MatchString(body.Telefono) {
			shared.JSONErr(w, 400, "El teléfono debe ser exactamente 8 dígitos numéricos.")
			return
		}
		// Complejidad de contraseña: mín 8 chars, mayúscula, minúscula, dígito, especial
		if len(body.Password) < 8 || !rePwUpper.MatchString(body.Password) ||
			!rePwLower.MatchString(body.Password) || !rePwDigit.MatchString(body.Password) ||
			!rePwSpec.MatchString(body.Password) {
			shared.JSONErr(w, 400, "La contraseña debe tener mínimo 8 caracteres, una mayúscula, una minúscula, un número y un carácter especial.")
			return
		}

		// Llamar al servicio de registro (transacción: CLIENTES + USUARIOS)
		idCliente, err := svc.Registrar(r.Context(), body.Identificacion, body.Nombre, body.PrimerApellido,
			body.Correo, body.Telefono, body.DireccionSennas, body.Password)
		if err != nil {
			// Errores ORA-20010 son errores de negocio (duplicados) — no son errores del sistema.
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

// logoutHandler maneja POST /api/auth/logout.
// Destruye la cookie de sesión (MaxAge=-1). El frontend también hace Auth.clear().
func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shared.ClearSession(w, r)
		shared.JSONMsg(w, "Sesión cerrada.")
	}
}

// checkSessionHandler maneja GET /api/auth/check-session.
// El frontend puede llamar esto al cargar una página para verificar si la cookie
// sigue válida y obtener los datos actualizados del usuario desde el servidor.
// Útil para validar que la sesión no expiró sin depender solo de localStorage.
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
			"rol":        c.Rol,
		})
	}
}
