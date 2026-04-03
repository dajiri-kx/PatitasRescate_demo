package shared

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

const sessionName = "patitas_session"

type ClienteSession struct {
	IDCliente int64
	Nombre    string
	Apellido  string
	Correo    string
	Telefono  string
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
