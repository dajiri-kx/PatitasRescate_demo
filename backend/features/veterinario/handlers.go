package veterinario

import (
	"database/sql"
	"log"
	"net/http"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewVetService(db)

	mux.HandleFunc("GET /api/vet/stats", vetStatsHandler(svc))
	mux.HandleFunc("GET /api/vet/citas", vetCitasHandler(svc))
	mux.HandleFunc("POST /api/vet/citas/estado", vetUpdateEstadoHandler(svc))
}

func vetStatsHandler(svc *VetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireVeterinario(w, r)
		if c == nil {
			return
		}
		stats, err := svc.GetStats(r.Context(), c.IDVeterinario)
		if err != nil {
			log.Printf("Vet stats error: %v", err)
			shared.JSONErr(w, 500, "Error al obtener estadísticas.")
			return
		}
		shared.JSONOk(w, stats)
	}
}

func vetCitasHandler(svc *VetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireVeterinario(w, r)
		if c == nil {
			return
		}
		list, err := svc.ListCitas(r.Context(), c.IDVeterinario)
		if err != nil {
			log.Printf("Vet citas error: %v", err)
			shared.JSONErr(w, 500, "Error al obtener citas.")
			return
		}
		if list == nil {
			list = []VetCita{}
		}
		shared.JSONOk(w, list)
	}
}

func vetUpdateEstadoHandler(svc *VetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := shared.RequireVeterinario(w, r)
		if c == nil {
			return
		}
		var body struct {
			ID     int64  `json:"id"`
			Estado string `json:"estado"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		if body.Estado != "Completada" && body.Estado != "Cancelada" {
			shared.JSONErr(w, 400, "Estado inválido. Use 'Completada' o 'Cancelada'.")
			return
		}
		if err := svc.UpdateEstadoCita(r.Context(), c.IDVeterinario, body.ID, body.Estado); err != nil {
			log.Printf("Vet update cita: %v", err)
			shared.JSONErr(w, 500, "Error al actualizar cita.")
			return
		}
		shared.JSONMsg(w, "Cita actualizada.")
	}
}
