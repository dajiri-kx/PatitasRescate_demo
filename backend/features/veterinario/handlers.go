/*
veterinario/handlers.go — Capa HTTP del portal veterinario.

SEGURIDAD:
Todos los handlers usan shared.RequireVeterinario(w, r) que verifica:
1. Sesión activa (cookie válida).
2. Rol == 2 (Veterinario).
3. IDVeterinario != 0 (vinculado a un registro en VETERINARIOS).
Si falla, responde 403 y retorna nil.

ENDPOINTS (3 en total):
GET  /api/vet/stats       → Métricas del dashboard vet (filtradas por su ID)
GET  /api/vet/citas       → Lista de citas asignadas a este veterinario
POST /api/vet/citas/estado → Cambiar estado (solo Completada o Cancelada)

DIFERENCIA CON ADMIN:
- El admin usa RequireAdmin y puede cambiar a cualquier estado.
- El vet usa RequireVeterinario y solo puede marcar Completada o Cancelada.
- El vet no puede reactivar una cita cancelada (no hay "Activa" en su whitelist).
*/
package veterinario

import (
	"database/sql"
	"log"
	"net/http"

	"patitas-backend/shared"
)

// RegisterRoutes registra los 3 endpoints del portal veterinario.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewVetService(db)

	mux.HandleFunc("GET /api/vet/stats", vetStatsHandler(svc))
	mux.HandleFunc("GET /api/vet/citas", vetCitasHandler(svc))
	mux.HandleFunc("POST /api/vet/citas/estado", vetUpdateEstadoHandler(svc))
}

// vetStatsHandler retorna las 3 métricas filtradas por el veterinario logueado.
// c.IDVeterinario viene de la sesión (asignado durante login).
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

// vetCitasHandler retorna la agenda completa del veterinario.
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
		// nil → []vacío para consistencia JSON.
		if list == nil {
			list = []VetCita{}
		}
		shared.JSONOk(w, list)
	}
}

// vetUpdateEstadoHandler cambia el estado de una cita del veterinario.
// WHITELIST RESTRINGIDA: solo "Completada" y "Cancelada".
// El veterinario NO puede reactivar citas (a diferencia del admin).
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
