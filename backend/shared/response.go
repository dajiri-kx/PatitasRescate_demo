/*
response.go — Helpers para respuestas JSON estandarizadas.

FLUJO DE DATOS:
Todos los handlers de la API devuelven un JSON con la misma estructura (Response).
Esto permite que el frontend (api.js → apiFetch) maneje las respuestas de forma
uniforme:
  - Éxito:  { "ok": true, "data": {...} }       o  { "ok": true, "message": "..." }
  - Error:  { "ok": false, "error": "..." }

El frontend verifica res.ok (HTTP status) y luego data.ok / data.error.

FUNCIONES HELPER:
- JSONOk(w, data)         → 200 + ok:true + data (para listas, objetos, etc.)
- JSONMsg(w, msg)         → 200 + ok:true + message (para confirmaciones sin data)
- JSONErr(w, status, msg) → status + ok:false + error (para errores descriptivos)
- DecodeBody(r, v)        → Decodifica el body JSON del request en un struct Go
*/
package shared

import (
	"encoding/json"
	"net/http"
)

// Response es el sobre (envelope) estándar de todas las respuestas JSON de la API.
// Los campos omitempty evitan enviar claves vacías innecesariamente.
type Response struct {
	OK      bool        `json:"ok"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON escribe una respuesta con el status HTTP y el struct Response dado.
// Todos los demás helpers (JSONOk, JSONMsg, JSONErr) llaman a esta función.
func JSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// JSONOk responde 200 con ok:true y el campo data (cualquier tipo serializable).
// Ejemplo de uso: shared.JSONOk(w, listaDeServicios)
func JSONOk(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{OK: true, Data: data})
}

// JSONMsg responde 200 con ok:true y un mensaje de texto (sin data).
// Ejemplo: shared.JSONMsg(w, "Servicio creado exitosamente.")
func JSONMsg(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusOK, Response{OK: true, Message: msg})
}

// JSONErr responde con un código de error HTTP y un mensaje descriptivo.
// El frontend lee data.error para mostrar al usuario.
// Ejemplo: shared.JSONErr(w, 400, "Todos los campos son obligatorios.")
func JSONErr(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, Response{OK: false, Error: msg})
}

// DecodeBody decodifica el body JSON del request HTTP en el struct apuntado por v.
// Cierra el body después de leer. Si falla, el handler responde con JSONErr 400.
func DecodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
