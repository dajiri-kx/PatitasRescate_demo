/*
admin/service.go — Lógica de negocio para el panel de administración.

ACCESO: Solo usuarios con ROL=0 (Admin). El middleware RequireAdmin verifica esto.

OPERACIONES DISPONIBLES:
 1. GetStats — 4 conteos para las tarjetas del dashboard admin:
    servicios totales, veterinarios, clientes registrados, citas activas.
 2. CRUD Servicios — Crear, listar, editar, eliminar servicios veterinarios.
 3. CRUD Veterinarios — Crear, listar, editar, eliminar veterinarios.
 4. ListClientes — Solo lectura (el admin no crea/edita clientes).
 5. ListCitas + UpdateEstadoCita — Ver todas las citas y cambiar su estado.

FLUJO DE DATOS:
Frontend admin (panel-layout SPA) → apiGet/apiPost('/admin/...') →
handlers.go (RequireAdmin) → service.go (queries directas a MariaDB) →
JSON response → frontend renderiza en tablas/formularios modales.

PATRÓN CRUD:
Cada entidad (Servicios, Veterinarios) tiene: List, Create, Update, Delete.
- List: SELECT con IFNULL para campos nullable (evita errores de Scan).
- Create: INSERT → retorna LastInsertId.
- Update: UPDATE WHERE ID=? (no verifica existencia — falla silenciosamente).
- Delete: DELETE WHERE ID=? (mismo comportamiento).
*/
package admin

import (
	"context"
	"database/sql"
)

type AdminService struct {
	db *sql.DB
}

func NewAdminService(db *sql.DB) *AdminService {
	return &AdminService{db: db}
}

// ---------- Dashboard stats ----------

// DashboardStats son las 4 métricas que se muestran en las tarjetas del dashboard admin.
type DashboardStats struct {
	Servicios    int `json:"servicios"`
	Veterinarios int `json:"veterinarios"`
	Clientes     int `json:"clientes"`
	CitasActivas int `json:"citas_activas"`
}

// GetStats ejecuta 4 COUNT queries para las tarjetas del dashboard.
// Usa un slice de structs para evitar repetir el patrón QueryRow+Scan.
func (s *AdminService) GetStats(ctx context.Context) (*DashboardStats, error) {
	st := &DashboardStats{}
	queries := []struct {
		dest  *int
		query string
	}{
		{&st.Servicios, `SELECT COUNT(*) FROM SERVICIOS`},
		{&st.Veterinarios, `SELECT COUNT(*) FROM VETERINARIOS`},
		{&st.Clientes, `SELECT COUNT(*) FROM CLIENTES`},
		{&st.CitasActivas, `SELECT COUNT(*) FROM CITAS WHERE ESTADO = 'Activa'`},
	}
	for _, q := range queries {
		if err := s.db.QueryRowContext(ctx, q.query).Scan(q.dest); err != nil {
			return nil, err
		}
	}
	return st, nil
}

// ---------- Servicios CRUD ----------

// ServicioRow es el struct para listar/editar servicios en el panel admin.
// Los JSON tags coinciden con los nombres de columnas de MariaDB en mayúsculas
// (convención usada en todo el proyecto para facilitar el mapeo frontend-backend).
type ServicioRow struct {
	ID          int64   `json:"ID_SERVICIO"`
	Nombre      string  `json:"NOMBRE_SERVICIO"`
	Descripcion string  `json:"DESCRIPCION"`
	Precio      float64 `json:"PRECIO"`
	Duracion    int     `json:"DURACION_MINUTOS"`
	Categoria   string  `json:"CATEGORIA"`
}

// ListServicios retorna todos los servicios ordenados por categoría y nombre.
// IFNULL protege contra columnas NULL que causarían error en Scan.
func (s *AdminService) ListServicios(ctx context.Context) ([]ServicioRow, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ID_SERVICIO, NOMBRE_SERVICIO, IFNULL(DESCRIPCION,''), PRECIO, IFNULL(DURACION_MINUTOS,0), CATEGORIA
		 FROM SERVICIOS ORDER BY CATEGORIA, NOMBRE_SERVICIO`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ServicioRow
	for rows.Next() {
		var r ServicioRow
		if err := rows.Scan(&r.ID, &r.Nombre, &r.Descripcion, &r.Precio, &r.Duracion, &r.Categoria); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

func (s *AdminService) CreateServicio(ctx context.Context, nombre, descripcion, categoria string, precio float64, duracion int) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA) VALUES (?,?,?,?,?)`,
		nombre, descripcion, precio, duracion, categoria)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *AdminService) UpdateServicio(ctx context.Context, id int64, nombre, descripcion, categoria string, precio float64, duracion int) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE SERVICIOS SET NOMBRE_SERVICIO=?, DESCRIPCION=?, PRECIO=?, DURACION_MINUTOS=?, CATEGORIA=? WHERE ID_SERVICIO=?`,
		nombre, descripcion, precio, duracion, categoria, id)
	return err
}

func (s *AdminService) DeleteServicio(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM SERVICIOS WHERE ID_SERVICIO=?`, id)
	return err
}

// ---------- Veterinarios CRUD ----------

// VeterinarioRow usa JSON tag "DIDENTIDAD_VETERINARIO" para coincidir con la
// columna de MariaDB (corregido de "DIDENTIDAD" genérico).
type VeterinarioRow struct {
	ID           int64  `json:"ID_VETERINARIO"`
	DIdentidad   string `json:"DIDENTIDAD_VETERINARIO"`
	Nombre       string `json:"NOMBRE"`
	Especialidad string `json:"ESPECIALIDAD"`
	Telefono     string `json:"TELEFONO"`
	Correo       string `json:"CORREO"`
	Rol          string `json:"ROL"`
}

func (s *AdminService) ListVeterinarios(ctx context.Context) ([]VeterinarioRow, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ID_VETERINARIO, DIDENTIDAD_VETERINARIO, NOMBRE, IFNULL(ESPECIALIDAD,''), IFNULL(TELEFONO,''), IFNULL(CORREO,''), IFNULL(ROL,'')
		 FROM VETERINARIOS ORDER BY NOMBRE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []VeterinarioRow
	for rows.Next() {
		var r VeterinarioRow
		if err := rows.Scan(&r.ID, &r.DIdentidad, &r.Nombre, &r.Especialidad, &r.Telefono, &r.Correo, &r.Rol); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

func (s *AdminService) CreateVeterinario(ctx context.Context, didentidad, nombre, especialidad, telefono, correo, rol string) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO VETERINARIOS (DIDENTIDAD_VETERINARIO, NOMBRE, ESPECIALIDAD, TELEFONO, CORREO, ROL) VALUES (?,?,?,?,?,?)`,
		didentidad, nombre, especialidad, telefono, correo, rol)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *AdminService) UpdateVeterinario(ctx context.Context, id int64, didentidad, nombre, especialidad, telefono, correo, rol string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE VETERINARIOS SET DIDENTIDAD_VETERINARIO=?, NOMBRE=?, ESPECIALIDAD=?, TELEFONO=?, CORREO=?, ROL=? WHERE ID_VETERINARIO=?`,
		didentidad, nombre, especialidad, telefono, correo, rol, id)
	return err
}

func (s *AdminService) DeleteVeterinario(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM VETERINARIOS WHERE ID_VETERINARIO=?`, id)
	return err
}

// ---------- Clientes (read-only) ----------

// ClienteRow usa JSON tag "DIDENTIDAD_CLIENTE" para coincidir con la columna real.
type ClienteRow struct {
	ID         int64  `json:"ID_CLIENTE"`
	DIdentidad string `json:"DIDENTIDAD_CLIENTE"`
	Nombre     string `json:"NOMBRE"`
	Apellido   string `json:"APELLIDO"`
	Email      string `json:"EMAIL"`
	Telefono   string `json:"TELEFONO"`
	Registro   string `json:"FECHA_REGISTRO"`
}

// ListClientes retorna todos los clientes. Solo lectura — el admin no crea clientes.
func (s *AdminService) ListClientes(ctx context.Context) ([]ClienteRow, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ID_CLIENTE, DIDENTIDAD_CLIENTE, NOMBRE, APELLIDO, EMAIL, TELEFONO,
		        DATE_FORMAT(FECHA_REGISTRO, '%Y-%m-%d') AS FECHA_REGISTRO
		 FROM CLIENTES ORDER BY NOMBRE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ClienteRow
	for rows.Next() {
		var r ClienteRow
		if err := rows.Scan(&r.ID, &r.DIdentidad, &r.Nombre, &r.Apellido, &r.Email, &r.Telefono, &r.Registro); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

// ---------- Citas management ----------

// CitaRow incluye datos de mascota, cliente y veterinario resueltos por JOINs.
// Total se calcula con subquery SUM de precios de servicios vinculados.
type CitaRow struct {
	ID          int64   `json:"ID_CITA"`
	FechaCita   string  `json:"FECHA_CITA"`
	Estado      string  `json:"ESTADO"`
	Mascota     string  `json:"MASCOTA"`
	Cliente     string  `json:"CLIENTE"`
	Veterinario string  `json:"VETERINARIO"`
	Total       float64 `json:"TOTAL"`
}

// ListCitas retorna TODAS las citas del sistema (sin filtro por cliente).
// El admin ve citas de todos los clientes con toda la información relevante.
func (s *AdminService) ListCitas(ctx context.Context) ([]CitaRow, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT c.ID_CITA, DATE_FORMAT(c.FECHA_CITA, '%Y-%m-%d %H:%i') AS FECHA_CITA,
		        c.ESTADO, m.NOMBRE AS MASCOTA,
		        CONCAT(cl.NOMBRE, ' ', cl.APELLIDO) AS CLIENTE,
		        v.NOMBRE AS VETERINARIO,
		        IFNULL((SELECT SUM(s.PRECIO) FROM CITAS_SERVICIOS cs JOIN SERVICIOS s ON cs.ID_SERVICIO=s.ID_SERVICIO WHERE cs.ID_CITA=c.ID_CITA), 0) AS TOTAL
		 FROM CITAS c
		 JOIN MASCOTAS m ON c.ID_MASCOTA = m.ID_MASCOTA
		 JOIN CLIENTES cl ON m.ID_CLIENTE = cl.ID_CLIENTE
		 JOIN VETERINARIOS v ON c.ID_VETERINARIO = v.ID_VETERINARIO
		 ORDER BY c.FECHA_CITA DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []CitaRow
	for rows.Next() {
		var r CitaRow
		if err := rows.Scan(&r.ID, &r.FechaCita, &r.Estado, &r.Mascota, &r.Cliente, &r.Veterinario, &r.Total); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

// UpdateEstadoCita cambia el estado de una cita. El admin puede poner
// Activa, Completada o Cancelada (validado en el handler).
func (s *AdminService) UpdateEstadoCita(ctx context.Context, idCita int64, estado string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE CITAS SET ESTADO=? WHERE ID_CITA=?`, estado, idCita)
	return err
}
