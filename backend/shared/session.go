/*
session.go — Manejo de sesiones y middleware de autorización por roles.

FLUJO DE DATOS (LOGIN → SESIÓN):
 1. El usuario envía correo+contraseña a POST /api/auth/login.
 2. auth/handlers.go llama a auth/service.go Login() → retorna ClienteData.
 3. loginHandler() mapea ClienteData → ClienteSession (este struct).
 4. SaveCliente() serializa ClienteSession dentro de una cookie firmada
    (gorilla/sessions usa gob para serializar structs en la cookie).
 5. En cada request posterior, el navegador envía la cookie automáticamente
    (credentials: 'include' en el frontend).
 6. GetCliente() deserializa la cookie → devuelve *ClienteSession o nil.

ROLES (TINYINT en tabla USUARIOS):

	0 = Admin   → acceso al panel administrativo (/api/admin/*)
	1 = Cliente  → acceso a citas, mascotas, facturas, checkout
	2 = Veterinario → acceso al portal veterinario (/api/vet/*)

MIDDLEWARE (RequireAuth / RequireAdmin / RequireVeterinario):
  - Se llaman al inicio de cada handler. Retornan *ClienteSession si pasa
    la validación, o nil si falla (y ya enviaron la respuesta 401/403).
  - El patrón en handlers es:  c := shared.RequireAuth(w, r); if c == nil { return }
  - Esto evita necesitar middleware HTTP clásico (wrapper de handlers).
*/
package shared

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

// Store es el almacén global de cookies firmadas. Se inicializa en InitSessionStore().
var Store *sessions.CookieStore

// sessionName es el nombre de la cookie que se envía al navegador.
const sessionName = "patitas_session"

// Constantes de roles — coinciden con el TINYINT de la columna USUARIOS.ROL.
// Roles: 0 = Admin, 1 = Cliente, 2 = Veterinario
const (
	RolAdmin       = 0
	RolCliente     = 1
	RolVeterinario = 2
)

// ClienteSession contiene los datos del usuario autenticado que se guardan
// en la cookie de sesión. Se serializa/deserializa con encoding/gob.
// IDVeterinario es 0 para clientes y admins; solo tiene valor para veterinarios
// que tengan su tabla USUARIOS.ID_VETERINARIO vinculada a VETERINARIOS.
type ClienteSession struct {
	IDCliente     int64
	IDVeterinario int64
	Nombre        string
	Apellido      string
	Correo        string
	Telefono      string
	Rol           int
}

// InitSessionStore configura el almacén de cookies. Se llama una vez desde main.go.
// SESSION_KEY es la clave HMAC para firmar la cookie (si no se define, usa un default).
// Opciones de la cookie:
//   - Path "/" → la cookie viaja con todas las rutas.
//   - MaxAge 86400 → 24 horas de vida.
//   - HttpOnly → JavaScript no puede leer la cookie (protección XSS).
//   - SameSite Lax → la cookie se envía en navegación normal pero no en
//     requests cross-origin (protección CSRF básica).
//
// gob.Register() es necesario porque gorilla/sessions usa encoding/gob
// para serializar valores personalizados dentro de la cookie.
func InitSessionStore() {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		key = "patitas-al-rescate-session-secret-32b!"
	}
	Store = sessions.NewCookieStore([]byte(key))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	gob.Register(ClienteSession{})
}

// GetCliente extrae la sesión del usuario de la cookie en el request.
// Retorna nil si no hay sesión válida (cookie ausente, expirada, o corrupta).
// Esta función NO escribe al response — es seguro llamarla sin efectos secundarios.
func GetCliente(r *http.Request) *ClienteSession {
	sess, err := Store.Get(r, sessionName)
	if err != nil {
		return nil
	}
	val, ok := sess.Values["cliente"]
	if !ok {
		return nil
	}
	c, ok := val.(ClienteSession)
	if !ok {
		return nil
	}
	return &c
}

// SaveCliente guarda el struct ClienteSession en la cookie de sesión.
// Se llama después de un login exitoso. La cookie firmada se envía al
// navegador como header Set-Cookie en el response.
func SaveCliente(w http.ResponseWriter, r *http.Request, c *ClienteSession) error {
	sess, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	sess.Values["cliente"] = *c
	return sess.Save(r, w)
}

// ClearSession destruye la sesión poniendo MaxAge=-1 (cookie expirada).
// El navegador eliminará la cookie al recibir esta respuesta.
func ClearSession(w http.ResponseWriter, r *http.Request) error {
	sess, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	return sess.Save(r, w)
}

// RequireAuth verifica que el request tiene una sesión válida (cualquier rol).
// Si no hay sesión, responde 401 y retorna nil → el handler debe hacer return.
// Usado por: citas, mascotas, facturas, checkout.
func RequireAuth(w http.ResponseWriter, r *http.Request) *ClienteSession {
	c := GetCliente(r)
	if c == nil {
		JSONErr(w, http.StatusUnauthorized, "No autenticado.")
		return nil
	}
	return c
}

// RequireAdmin = RequireAuth + verificación de Rol == 0 (Admin).
// Responde 403 si el usuario está autenticado pero no es admin.
// Usado por: todos los endpoints /api/admin/*.
func RequireAdmin(w http.ResponseWriter, r *http.Request) *ClienteSession {
	c := GetCliente(r)
	if c == nil {
		JSONErr(w, http.StatusUnauthorized, "No autenticado.")
		return nil
	}
	if c.Rol != RolAdmin {
		JSONErr(w, http.StatusForbidden, "Acceso denegado.")
		return nil
	}
	return c
}

// RequireVeterinario = RequireAuth + Rol==2 + IDVeterinario!=0.
// La doble verificación asegura que el usuario veterinario tenga su cuenta
// vinculada a un registro en la tabla VETERINARIOS (via USUARIOS.ID_VETERINARIO).
// Esto es necesario porque un usuario puede tener rol 2 pero no estar vinculado
// todavía (campo ID_VETERINARIO es nullable en USUARIOS).
// Usado por: todos los endpoints /api/vet/*.
func RequireVeterinario(w http.ResponseWriter, r *http.Request) *ClienteSession {
	c := GetCliente(r)
	if c == nil {
		JSONErr(w, http.StatusUnauthorized, "No autenticado.")
		return nil
	}
	if c.Rol != RolVeterinario {
		JSONErr(w, http.StatusForbidden, "Acceso denegado.")
		return nil
	}
	if c.IDVeterinario == 0 {
		JSONErr(w, http.StatusForbidden, "Cuenta de veterinario no vinculada.")
		return nil
	}
	return c
}
