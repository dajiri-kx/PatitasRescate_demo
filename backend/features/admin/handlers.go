package admin

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"patitas-backend/shared"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewAdminService(db)

	// Dashboard
	mux.HandleFunc("GET /api/admin/stats", statsHandler(svc))

	// Servicios CRUD
	mux.HandleFunc("GET /api/admin/servicios", listServiciosHandler(svc))
	mux.HandleFunc("POST /api/admin/servicios", createServicioHandler(svc))
	mux.HandleFunc("POST /api/admin/servicios/editar", updateServicioHandler(svc))
	mux.HandleFunc("POST /api/admin/servicios/eliminar", deleteServicioHandler(svc))

	// Veterinarios CRUD
	mux.HandleFunc("GET /api/admin/veterinarios", listVeterinariosHandler(svc))
	mux.HandleFunc("POST /api/admin/veterinarios", createVeterinarioHandler(svc))
	mux.HandleFunc("POST /api/admin/veterinarios/editar", updateVeterinarioHandler(svc))
	mux.HandleFunc("POST /api/admin/veterinarios/eliminar", deleteVeterinarioHandler(svc))

	// Clientes (read-only)
	mux.HandleFunc("GET /api/admin/clientes", listClientesHandler(svc))

	// Citas management
	mux.HandleFunc("GET /api/admin/citas", listCitasHandler(svc))
	mux.HandleFunc("POST /api/admin/citas/estado", updateEstadoCitaHandler(svc))
}

// ---------- Dashboard ----------

func statsHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		stats, err := svc.GetStats(r.Context())
		if err != nil {
			log.Printf("Admin stats error: %v", err)
			shared.JSONErr(w, 500, "Error al obtener estadísticas.")
			return
		}
		shared.JSONOk(w, stats)
	}
}

// ---------- Servicios ----------

func listServiciosHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		list, err := svc.ListServicios(r.Context())
		if err != nil {
			log.Printf("Admin servicios error: %v", err)
			shared.JSONErr(w, 500, "Error al obtener servicios.")
			return
		}
		if list == nil {
			list = []ServicioRow{}
		}
		shared.JSONOk(w, list)
	}
}

func createServicioHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		var body struct {
			Nombre      string  `json:"nombre"`
			Descripcion string  `json:"descripcion"`
			Precio      float64 `json:"precio"`
			Duracion    int     `json:"duracion"`
			Categoria   string  `json:"categoria"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		if body.Nombre == "" || body.Categoria == "" {
			shared.JSONErr(w, 400, "Nombre y categoría son obligatorios.")
			return
		}
		id, err := svc.CreateServicio(r.Context(), body.Nombre, body.Descripcion, body.Categoria, body.Precio, body.Duracion)
		if err != nil {
			log.Printf("Admin crear servicio: %v", err)
			shared.JSONErr(w, 500, "Error al crear servicio.")
			return
		}
		shared.JSONOk(w, map[string]interface{}{"id": id, "message": "Servicio creado."})
	}
}

func updateServicioHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		var body struct {
			ID          int64   `json:"id"`
			Nombre      string  `json:"nombre"`
			Descripcion string  `json:"descripcion"`
			Precio      float64 `json:"precio"`
			Duracion    int     `json:"duracion"`
			Categoria   string  `json:"categoria"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		if err := svc.UpdateServicio(r.Context(), body.ID, body.Nombre, body.Descripcion, body.Categoria, body.Precio, body.Duracion); err != nil {
			log.Printf("Admin editar servicio: %v", err)
			shared.JSONErr(w, 500, "Error al editar servicio.")
			return
		}
		shared.JSONMsg(w, "Servicio actualizado.")
	}
}

func deleteServicioHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		var body struct {
			ID int64 `json:"id"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		if err := svc.DeleteServicio(r.Context(), body.ID); err != nil {
			log.Printf("Admin eliminar servicio: %v", err)
			shared.JSONErr(w, 500, "Error al eliminar servicio.")
			return
		}
		shared.JSONMsg(w, "Servicio eliminado.")
	}
}

// ---------- Veterinarios ----------

func listVeterinariosHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		list, err := svc.ListVeterinarios(r.Context())
		if err != nil {
			log.Printf("Admin veterinarios error: %v", err)
			shared.JSONErr(w, 500, "Error al obtener veterinarios.")
			return
		}
		if list == nil {
			list = []VeterinarioRow{}
		}
		shared.JSONOk(w, list)
	}
}

func createVeterinarioHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		var body struct {
			DIdentidad   string `json:"didentidad"`
			Nombre       string `json:"nombre"`
			Especialidad string `json:"especialidad"`
			Telefono     string `json:"telefono"`
			Correo       string `json:"correo"`
			Rol          string `json:"rol"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		if body.Nombre == "" || body.DIdentidad == "" {
			shared.JSONErr(w, 400, "Nombre e identificación son obligatorios.")
			return
		}
		id, err := svc.CreateVeterinario(r.Context(), body.DIdentidad, body.Nombre, body.Especialidad, body.Telefono, body.Correo, body.Rol)
		if err != nil {
			log.Printf("Admin crear veterinario: %v", err)
			shared.JSONErr(w, 500, "Error al crear veterinario.")
			return
		}
		shared.JSONOk(w, map[string]interface{}{"id": id, "message": "Veterinario creado."})
	}
}

func updateVeterinarioHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		var body struct {
			ID           int64  `json:"id"`
			DIdentidad   string `json:"didentidad"`
			Nombre       string `json:"nombre"`
			Especialidad string `json:"especialidad"`
			Telefono     string `json:"telefono"`
			Correo       string `json:"correo"`
			Rol          string `json:"rol"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		if err := svc.UpdateVeterinario(r.Context(), body.ID, body.DIdentidad, body.Nombre, body.Especialidad, body.Telefono, body.Correo, body.Rol); err != nil {
			log.Printf("Admin editar veterinario: %v", err)
			shared.JSONErr(w, 500, "Error al editar veterinario.")
			return
		}
		shared.JSONMsg(w, "Veterinario actualizado.")
	}
}

func deleteVeterinarioHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		var body struct {
			ID int64 `json:"id"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		if err := svc.DeleteVeterinario(r.Context(), body.ID); err != nil {
			log.Printf("Admin eliminar veterinario: %v", err)
			shared.JSONErr(w, 500, "Error al eliminar veterinario.")
			return
		}
		shared.JSONMsg(w, "Veterinario eliminado.")
	}
}

// ---------- Clientes ----------

func listClientesHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		list, err := svc.ListClientes(r.Context())
		if err != nil {
			log.Printf("Admin clientes error: %v", err)
			shared.JSONErr(w, 500, "Error al obtener clientes.")
			return
		}
		if list == nil {
			list = []ClienteRow{}
		}
		shared.JSONOk(w, list)
	}
}

// ---------- Citas ----------

func listCitasHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		list, err := svc.ListCitas(r.Context())
		if err != nil {
			log.Printf("Admin citas error: %v", err)
			shared.JSONErr(w, 500, "Error al obtener citas.")
			return
		}
		if list == nil {
			list = []CitaRow{}
		}
		shared.JSONOk(w, list)
	}
}

func updateEstadoCitaHandler(svc *AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shared.RequireAdmin(w, r) == nil {
			return
		}
		var body struct {
			ID     string `json:"id"`
			Estado string `json:"estado"`
		}
		if err := shared.DecodeBody(r, &body); err != nil {
			shared.JSONErr(w, 400, "Datos inválidos.")
			return
		}
		id, err := strconv.ParseInt(body.ID, 10, 64)
		if err != nil {
			shared.JSONErr(w, 400, "ID inválido.")
			return
		}
		validStates := map[string]bool{"Activa": true, "Completada": true, "Cancelada": true}
		if !validStates[body.Estado] {
			shared.JSONErr(w, 400, "Estado inválido. Use: Activa, Completada, Cancelada.")
			return
		}
		if err := svc.UpdateEstadoCita(r.Context(), id, body.Estado); err != nil {
			log.Printf("Admin estado cita: %v", err)
			shared.JSONErr(w, 500, "Error al actualizar estado.")
			return
		}
		shared.JSONMsg(w, "Estado actualizado.")
	}
}
