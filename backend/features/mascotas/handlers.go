/*
mascotas/handlers.go — Capa HTTP para gestión de mascotas del cliente.

ENDPOINTS:
GET  /api/mascotas         → Lista completa de mascotas del cliente (mis-mascotas page)
GET  /api/mascotas/nombres → Solo ID+nombre para dropdowns (agendar-cita)
POST /api/mascotas/agregar → Registrar nueva mascota

FLUJO: Todos usan RequireAuth → c.IDCliente → service → JSONOk.
Las mascotas siempre están filtradas por el cliente de la sesión.
*/
package mascotas

import (
	"database/sql"
	"log"
	"net/http"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewMascotaService(db)
	mux.HandleFunc("GET /api/mascotas", obtenerHandler(svc))
	mux.HandleFunc("GET /api/mascotas/nombres", obtenerNombresHandler(svc))
	mux.HandleFunc("POST /api/mascotas/agregar", agregarHandler(svc))
}

// obtenerHandler → GET /api/mascotas
// Retorna lista completa con especie, raza, edad y datos del dueño.
func obtenerHandler(svc *MascotaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerPorCliente(r.Context(), c.IDCliente)
		if err != nil {
			log.Printf("Error mascotas: %v", err)
			shared.JSONErr(w, 500, "Error al obtener mascotas.")
			return
		}
		if list == nil {
			list = []Mascota{}
		}
		shared.JSONOk(w, list)
	}
}

// obtenerNombresHandler → GET /api/mascotas/nombres
// Retorna subset ligero (ID + nombre) para el <select> del formulario agendar-cita.
func obtenerNombresHandler(svc *MascotaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerNombres(r.Context(), c.IDCliente)
		if err != nil {
			log.Printf("Error nombres mascotas: %v", err)
			shared.JSONErr(w, 500, "Error al obtener mascotas.")
			return
		}
		if list == nil {
			list = []MascotaNombre{}
		}
		shared.JSONOk(w, list)
	}
}

// agregarHandler → POST /api/mascotas/agregar
// Recibe {nombre, especie, raza, edad} y crea la mascota vinculada al c.IDCliente.
func agregarHandler(svc *MascotaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}

		var body struct {
			Nombre  string `json:"nombre"`
			Especie string `json:"especie"`
			Raza    string `json:"raza"`
			Edad    int    `json:"edad"` // Edad en meses, se almacena como MESES en BD
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		if body.Nombre == "" || body.Especie == "" || body.Raza == "" {
			shared.JSONErr(w, 400, "Todos los campos son obligatorios.")
			return
		}

		if err := svc.Agregar(r.Context(), body.Nombre, body.Especie, body.Raza, body.Edad, c.IDCliente); err != nil {
			log.Printf("Error agregar mascota: %v", err)
			shared.JSONErr(w, 500, "Error al registrar mascota.")
			return
		}

		shared.JSONMsg(w, "Mascota registrada exitosamente.")
	}
}
