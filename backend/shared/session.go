package shared

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

const sessionName = "patitas_session"

// Roles: 0 = Admin, 1 = Cliente, 2 = Veterinario
const (
	RolAdmin       = 0
	RolCliente     = 1
	RolVeterinario = 2
)

type ClienteSession struct {
	IDCliente     int64
	IDVeterinario int64
	Nombre        string
	Apellido      string
	Correo        string
	Telefono      string
	Rol           int
}

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

func SaveCliente(w http.ResponseWriter, r *http.Request, c *ClienteSession) error {
	sess, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	sess.Values["cliente"] = *c
	return sess.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	sess, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	return sess.Save(r, w)
}

func RequireAuth(w http.ResponseWriter, r *http.Request) *ClienteSession {
	c := GetCliente(r)
	if c == nil {
		JSONErr(w, http.StatusUnauthorized, "No autenticado.")
		return nil
	}
	return c
}

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
