/*
middleware.go — CORS (Cross-Origin Resource Sharing) middleware.

FLUJO DE DATOS:
En main.go, el handler final es: shared.CORS(mux).
Cada request HTTP pasa primero por este middleware antes de llegar al router.

PROPÓSITO:
Cuando el frontend se sirve desde un origen diferente al backend (ej: desarrollo
con live-server en localhost:5500 y backend en localhost:8080), el navegador
bloquea las peticiones cross-origin por seguridad. CORS agrega los headers
necesarios para que el navegador permita las peticiones.

Si CORS_ORIGIN no está definida (producción con mismo origen), NO se agregan
headers CORS — el navegador ya permite requests al mismo origen automáticamente.

PREFLIGHT: El navegador envía un OPTIONS antes de POST con Content-Type JSON.
Este middleware responde 204 (sin cuerpo) para que el navegador proceda.
Allow-Credentials:true es necesario para que la cookie de sesión viaje en CORS.
*/
package shared

import (
	"net/http"
	"os"
)

// CORS envuelve un http.Handler agregando headers CORS si CORS_ORIGIN está definida.
// Se aplica como la capa más externa: handler := shared.CORS(mux)
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Preflight: el navegador pregunta "¿puedo hacer este request?"
			// Respondemos 204 y no pasamos al router.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
