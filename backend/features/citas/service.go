package citas

import (
"context"
"database/sql"
)

type CitaService struct {
db *sql.DB
}

func NewCitaService(db *sql.DB) *CitaService {
return &CitaService{db: db}
}

type Cita struct {
IDCita      int64  `json:"ID_CITA"`
FechaCita   string `json:"FECHA_CITA"`
Estado      string `json:"ESTADO"`
Mascota     string `json:"MASCOTA"`
Veterinario string `json:"VETERINARIO"`
}

type CitaActiva struct {
IDCita    int64  `json:"ID_CITA"`
FechaCita string `json:"FECHA_CITA"`
Mascota   string `json:"MASCOTA"`
}

type Veterinario struct {
IDVeterinario int64  `json:"ID_VETERINARIO"`
Nombre        string `json:"NOMBRE"`
}

type Servicio struct {
IDServicio     int64  `json:"ID_SERVICIO"`
NombreServicio string `json:"NOMBRE_SERVICIO"`
Descripcion    string `json:"DESCRIPCION"`
}

func (s *CitaService) ObtenerPorCliente(ctx context.Context, idCliente int64) ([]Cita, error) {
rows, err := s.db.QueryContext(ctx,
`SELECT C.ID_CITA, TO_CHAR(C.FECHA_CITA, 'YYYY-MM-DD HH24:MI') AS FECHA_CITA,
        C.ESTADO, M.NOMBRE AS MASCOTA, V.NOMBRE AS VETERINARIO
 FROM CITAS_TABLAS.CITAS C
 JOIN USUARIOS_TABLAS.MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
 JOIN CITAS_TABLAS.VETERINARIOS V ON C.ID_VETERINARIO = V.ID_VETERINARIO
 WHERE M.ID_CLIENTE = :id_cliente
 ORDER BY C.FECHA_CITA DESC`,
sql.Named("id_cliente", idCliente),
)
if err != nil {
return nil, err
}
defer rows.Close()

var list []Cita
for rows.Next() {
var c Cita
if err := rows.Scan(&c.IDCita, &c.FechaCita, &c.Estado, &c.Mascota, &c.Veterinario); err != nil {
return nil, err
}
list = append(list, c)
}
return list, rows.Err()
}

func (s *CitaService) ObtenerActivas(ctx context.Context, idCliente int64) ([]CitaActiva, error) {
rows, err := s.db.QueryContext(ctx,
`SELECT C.ID_CITA, TO_CHAR(C.FECHA_CITA, 'YYYY-MM-DD HH24:MI') AS FECHA_CITA, M.NOMBRE AS MASCOTA
 FROM CITAS_TABLAS.CITAS C
 JOIN USUARIOS_TABLAS.MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
 WHERE M.ID_CLIENTE = :id_cliente AND C.ESTADO = 'Activa'
 ORDER BY C.FECHA_CITA`,
sql.Named("id_cliente", idCliente),
)
if err != nil {
return nil, err
}
defer rows.Close()

var list []CitaActiva
for rows.Next() {
var c CitaActiva
if err := rows.Scan(&c.IDCita, &c.FechaCita, &c.Mascota); err != nil {
return nil, err
}
list = append(list, c)
}
return list, rows.Err()
}

func (s *CitaService) ObtenerVeterinarios(ctx context.Context) ([]Veterinario, error) {
rows, err := s.db.QueryContext(ctx,
`SELECT ID_VETERINARIO, NOMBRE FROM CITAS_TABLAS.VETERINARIOS ORDER BY NOMBRE`)
if err != nil {
return nil, err
}
defer rows.Close()

var list []Veterinario
for rows.Next() {
var v Veterinario
if err := rows.Scan(&v.IDVeterinario, &v.Nombre); err != nil {
return nil, err
}
list = append(list, v)
}
return list, rows.Err()
}

func (s *CitaService) ObtenerServicios(ctx context.Context) ([]Servicio, error) {
rows, err := s.db.QueryContext(ctx,
`SELECT ID_SERVICIO, NOMBRE_SERVICIO, NVL(DESCRIPCION, ' ') AS DESCRIPCION
 FROM SERVICIOS_TABLAS.SERVICIOS ORDER BY NOMBRE_SERVICIO`)
if err != nil {
return nil, err
}
defer rows.Close()

var list []Servicio
for rows.Next() {
var sv Servicio
if err := rows.Scan(&sv.IDServicio, &sv.NombreServicio, &sv.Descripcion); err != nil {
return nil, err
}
list = append(list, sv)
}
return list, rows.Err()
}

func (s *CitaService) Agendar(ctx context.Context, idCliente, idMascota, idVeterinario int64, fechaCita, servicios string) error {
_, err := s.db.ExecContext(ctx,
`BEGIN agendarCita(:id_cliente, :id_mascota, :id_veterinario, TO_DATE(:fecha_cita, 'YYYY-MM-DD HH24:MI'), :servicios); END;`,
sql.Named("id_cliente", idCliente),
sql.Named("id_mascota", idMascota),
sql.Named("id_veterinario", idVeterinario),
sql.Named("fecha_cita", fechaCita),
sql.Named("servicios", servicios),
)
return err
}

func (s *CitaService) Cancelar(ctx context.Context, idCita, idCliente int64) (bool, error) {
_, err := s.db.ExecContext(ctx,
`DELETE FROM CITAS_TABLAS.DETALLE_CITAS
 WHERE ID_CITA = :id_cita AND ID_CITA IN (
     SELECT C.ID_CITA FROM CITAS_TABLAS.CITAS C
     JOIN USUARIOS_TABLAS.MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
     WHERE M.ID_CLIENTE = :id_cliente AND C.ID_CITA = :id_cita2
 )`,
sql.Named("id_cita", idCita),
sql.Named("id_cliente", idCliente),
sql.Named("id_cita2", idCita),
)
if err != nil {
return false, err
}

result, err := s.db.ExecContext(ctx,
`DELETE FROM CITAS_TABLAS.CITAS
 WHERE ID_CITA = :id_cita AND ID_MASCOTA IN (
     SELECT ID_MASCOTA FROM USUARIOS_TABLAS.MASCOTAS WHERE ID_CLIENTE = :id_cliente
 )`,
sql.Named("id_cita", idCita),
sql.Named("id_cliente", idCliente),
)
if err != nil {
return false, err
}

n, _ := result.RowsAffected()
return n > 0, nil
}
