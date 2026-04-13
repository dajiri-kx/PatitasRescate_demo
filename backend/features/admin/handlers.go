/*
admin/handlers.go — Capa HTTP del panel de administración.

SEGURIDAD:
Cada handler llama a shared.RequireAdmin(w, r) como primera línea.
Si el usuario no tiene sesión activa con Rol==0, RequireAdmin responde
un 403 y retorna nil → el handler hace "return" y no ejecuta nada más.

ENDPOINTS (12 en total, todos bajo /api/admin/):
┌──────────────────────────────────┬────────┬──────────────────────────────┐
│ Ruta                             │ Método │ Descripción                  │
├──────────────────────────────────┼────────┼──────────────────────────────┤
│ /api/admin/stats                 │ GET    │ Métricas del dashboard       │
│ /api/admin/servicios             │ GET    │ Listar servicios             │
│ /api/admin/servicios             │ POST   │ Crear servicio               │
│ /api/admin/servicios/editar      │ POST   │ Editar servicio              │
│ /api/admin/servicios/eliminar    │ POST   │ Eliminar servicio            │
│ /api/admin/veterinarios          │ GET    │ Listar veterinarios          │
│ /api/admin/veterinarios          │ POST   │ Crear veterinario            │
│ /api/admin/veterinarios/editar   │ POST   │ Editar veterinario           │
│ /api/admin/veterinarios/eliminar │ POST   │ Eliminar veterinario         │
│ /api/admin/clientes              │ GET    │ Listar clientes (solo lect.) │
│ /api/admin/citas                 │ GET    │ Listar todas las citas       │
│ /api/admin/citas/estado          │ POST   │ Cambiar estado de cita       │
└──────────────────────────────────┴────────┴──────────────────────────────┘

PATRÓN HANDLER:
 1. RequireAdmin → corta si no es admin.
 2. Para GET: llamar service → responder JSON.
 3. Para POST: DecodeBody → validar campos → llamar service → responder JSON.
 4. nil-guard en listas: si el slice es nil, devolvemos []vacío para que
    el frontend reciba un JSON array [] en vez de null.

NOTA sobre mutaciones: Se usa POST para editar/eliminar en vez de PUT/DELETE
porque el frontend simplificado solo usa GET y POST (sin REST puro).
*/
package admin

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"patitas-backend/shared"
)

// RegisterRoutes registra los 12 endpoints del panel admin.
// Todos los handlers están protegidos internamente por RequireAdmin.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	svc := NewAdminService(db)

	// Dashboard — tarjetas de métricas
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

	// Clientes (solo lectura — no CRUD)
	mux.HandleFunc("GET /api/admin/clientes", listClientesHandler(svc))

	// Citas — ver todas + cambiar estado
	mux.HandleFunc("GET /api/admin/citas", listCitasHandler(svc))
	mux.HandleFunc("POST /api/admin/citas/estado", updateEstadoCitaHandler(svc))
}

// ---------- Dashboard ----------

// statsHandler retorna las 4 métricas del dashboard admin.
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

// ---------- Servicios CRUD ----------

// listServiciosHandler retorna todos los servicios para la tabla del panel.
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
		// nil → []vacío para que el frontend reciba un array, no null.
		if list == nil {
			list = []ServicioRow{}
		}
		shared.JSONOk(w, list)
	}
}

// createServicioHandler decodifica el body JSON y crea el servicio.
// body usa tags minúsculas ("nombre", "precio") — convención del frontend.
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

// ---------- Veterinarios CRUD ----------

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

// ---------- Clientes (solo lectura) ----------

// listClientesHandler — el admin solo puede ver clientes, no crear/editar.
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

// ---------- Citas (gestión de estado) ----------

// listCitasHandler retorna TODAS las citas del sistema con JOINs resueltos.
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

// updateEstadoCitaHandler cambia el estado de una cita.
// MÁQUINA DE ESTADOS permitida: Activa, Completada, Cancelada.
// El admin puede poner cualquiera de los 3 (el veterinario solo Completada/Cancelada).
// NOTA: body.ID llega como string desde el frontend → se parsea con strconv.
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
		// El frontend envía ID como string → parse a int64.
		id, err := strconv.ParseInt(body.ID, 10, 64)
		if err != nil {
			shared.JSONErr(w, 400, "ID inválido.")
			return
		}
		// Whitelist de estados válidos — previene inyección de valores arbitrarios.
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
