package citas

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewCitaService(db)
	mux.HandleFunc("GET /api/citas", obtenerHandler(svc))
	mux.HandleFunc("GET /api/citas/activas", obtenerActivasHandler(svc))
	mux.HandleFunc("GET /api/citas/veterinarios", obtenerVeterinariosHandler(svc))
	mux.HandleFunc("GET /api/citas/servicios", obtenerServiciosHandler(svc))
	mux.HandleFunc("POST /api/citas/agendar", agendarHandler(svc))
	mux.HandleFunc("POST /api/citas/cancelar", cancelarHandler(svc))
}

func obtenerHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerPorCliente(r.Context(), c.IDCliente)
		if err != nil {
			log.Printf("Error citas: %v", err)
			shared.JSONErr(w, 500, "Error al obtener citas.")
			return
		}
		if list == nil {
			list = []Cita{}
		}
		shared.JSONOk(w, list)
	}
}

func obtenerActivasHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerActivas(r.Context(), c.IDCliente)
		if err != nil {
			log.Printf("Error citas activas: %v", err)
			shared.JSONErr(w, 500, "Error al obtener citas activas.")
			return
		}
		if list == nil {
			list = []CitaActiva{}
		}
		shared.JSONOk(w, list)
	}
}

func obtenerVeterinariosHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerVeterinarios(r.Context())
		if err != nil {
			log.Printf("Error veterinarios: %v", err)
			shared.JSONErr(w, 500, "Error al obtener veterinarios.")
			return
		}
		if list == nil {
			list = []Veterinario{}
		}
		shared.JSONOk(w, list)
	}
}

func obtenerServiciosHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}
		list, err := svc.ObtenerServicios(r.Context())
		if err != nil {
			log.Printf("Error servicios: %v", err)
			shared.JSONErr(w, 500, "Error al obtener servicios.")
			return
		}
		if list == nil {
			list = []Servicio{}
		}
		shared.JSONOk(w, list)
	}
}

func agendarHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}

		var body struct {
			IDMascota   string   `json:"id_mascota"`
			Fecha       string   `json:"fecha"`
			Hora        string   `json:"hora"`
			Servicio    []string `json:"servicio"`
			Veterinario string   `json:"veterinario"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		idMascota, err := strconv.ParseInt(body.IDMascota, 10, 64)
		if err != nil {
			shared.JSONErr(w, 400, "Mascota inválida.")
			return
		}
		idVet, err := strconv.ParseInt(body.Veterinario, 10, 64)
		if err != nil {
			shared.JSONErr(w, 400, "Veterinario inválido.")
			return
		}

		fechaCita := body.Fecha + " " + body.Hora
		serviciosList := strings.Join(body.Servicio, ",")

		err = svc.Agendar(r.Context(), c.IDCliente, idMascota, idVet, fechaCita, serviciosList)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "ORA-200") {
				// Extraer mensaje después del código
				if idx := strings.Index(errMsg, ": "); idx != -1 {
					shared.JSONErr(w, 400, errMsg[idx+2:])
				} else {
					shared.JSONErr(w, 400, errMsg)
				}
				return
			}
			log.Printf("Error agendar: %v", err)
			shared.JSONErr(w, 500, "Error al agendar la cita.")
			return
		}

		shared.JSONMsg(w, "Cita agendada exitosamente.")
	}
}

func cancelarHandler(svc *CitaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireAuth(w, r)
		if c == nil {
			return
		}

		var body struct {
			IDCita string `json:"id_cita"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}

		idCita, err := strconv.ParseInt(body.IDCita, 10, 64)
		if err != nil || idCita <= 0 {
			shared.JSONErr(w, 400, "ID de cita inválido.")
			return
		}

		ok, err := svc.Cancelar(r.Context(), idCita, c.IDCliente)
		if err != nil {
			log.Printf("Error cancelar: %v", err)
			shared.JSONErr(w, 500, "Error al cancelar la cita.")
			return
		}
		if !ok {
			shared.JSONErr(w, 404, "Cita no encontrada o no autorizada.")
			return
		}

		shared.JSONMsg(w, "Cita cancelada exitosamente.")
	}
}
