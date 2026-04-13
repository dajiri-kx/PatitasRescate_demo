/*
citas/service.go — Lógica de negocio para gestión de citas veterinarias.

MODELO DE DATOS DE CITAS (relaciones entre tablas):
CLIENTES (1) → (N) MASCOTAS → (N) CITAS ← (1) VETERINARIOS
CITAS (N) ←→ (N) SERVICIOS   (vía tabla puente CITAS_SERVICIOS)
CITAS_SERVICIOS → (1) FACTURAS (la factura agrupa los servicios de una cita)

FLUJO DE DATOS — AGENDAR CITA (la operación más compleja del sistema):
Frontend envía: {id_mascota, fecha, hora, servicio[], veterinario}
→ handlers.go agendarHandler() une fecha+hora y servicios → llama svc.Agendar()
→ Agendar() ejecuta TODO en una transacción:
 1. Parsear IDs de servicios (vienen como string "1,3,5")
 2. Validar: cliente existe, mascota existe, mascota es del cliente
 3. Validar: no hay cita activa ese día para esa mascota
 4. Validar: veterinario existe y está disponible en esa fecha/hora
 5. INSERT CITA → obtener ID_CITA
 6. Para cada servicio: validar, INSERT CITAS_SERVICIOS, descontar stock productos
 7. Calcular total = SUM(precios de servicios)
 8. INSERT FACTURA con el total → obtener ID_FACTURA
 9. UPDATE CITAS_SERVICIOS para vincular con la factura
 10. COMMIT

→ Retorna ID_FACTURA al frontend (para redirigir a Stripe Checkout)

CÓDIGOS ORA-200xx:
Se usan como convención para errores de negocio (inspirados en Oracle PL/SQL).
El handler detecta "ORA-200" en el mensaje y responde 400 en vez de 500.
*/
package citas

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type CitaService struct {
	db *sql.DB
}

func NewCitaService(db *sql.DB) *CitaService {
	return &CitaService{db: db}
}

// Cita es el struct para la lista de citas del cliente (vista mis-citas).
// Incluye nombre de mascota y veterinario resueltos por JOINs.
type Cita struct {
	IDCita      int64  `json:"ID_CITA"`
	FechaCita   string `json:"FECHA_CITA"`
	Estado      string `json:"ESTADO"`
	Mascota     string `json:"MASCOTA"`
	Veterinario string `json:"VETERINARIO"`
}

// CitaActiva es un subset de Cita usado en el formulario de cancelar cita.
type CitaActiva struct {
	IDCita    int64  `json:"ID_CITA"`
	FechaCita string `json:"FECHA_CITA"`
	Mascota   string `json:"MASCOTA"`
}

// Veterinario para el dropdown del formulario de agendar cita.
type Veterinario struct {
	IDVeterinario int64  `json:"ID_VETERINARIO"`
	Nombre        string `json:"NOMBRE"`
}

// Servicio para las cards de selección en el formulario de agendar cita.
type Servicio struct {
	IDServicio     int64   `json:"ID_SERVICIO"`
	NombreServicio string  `json:"NOMBRE_SERVICIO"`
	Descripcion    string  `json:"DESCRIPCION"`
	Precio         float64 `json:"PRECIO"`
	Categoria      string  `json:"CATEGORIA"`
}

// ObtenerPorCliente devuelve TODAS las citas de un cliente (activas, completadas, canceladas).
// Se accede a las citas del cliente a través de la mascota: CITAS → MASCOTAS → ID_CLIENTE.
// Usado en: GET /api/citas (página "Mis Citas" del cliente).
func (s *CitaService) ObtenerPorCliente(ctx context.Context, idCliente int64) ([]Cita, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT C.ID_CITA, DATE_FORMAT(C.FECHA_CITA, '%Y-%m-%d %H:%i') AS FECHA_CITA,
		        C.ESTADO, M.NOMBRE AS MASCOTA, V.NOMBRE AS VETERINARIO
		 FROM CITAS C
		 JOIN MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
		 JOIN VETERINARIOS V ON C.ID_VETERINARIO = V.ID_VETERINARIO
		 WHERE M.ID_CLIENTE = ?
		 ORDER BY C.FECHA_CITA DESC`,
		idCliente,
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

// ObtenerActivas retorna solo citas con Estado='Activa' para el formulario de cancelar.
func (s *CitaService) ObtenerActivas(ctx context.Context, idCliente int64) ([]CitaActiva, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT C.ID_CITA, DATE_FORMAT(C.FECHA_CITA, '%Y-%m-%d %H:%i') AS FECHA_CITA, M.NOMBRE AS MASCOTA
		 FROM CITAS C
		 JOIN MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
		 WHERE M.ID_CLIENTE = ? AND C.ESTADO = 'Activa'
		 ORDER BY C.FECHA_CITA`,
		idCliente,
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

// ObtenerVeterinarios devuelve todos los veterinarios para el dropdown de agendar cita.
func (s *CitaService) ObtenerVeterinarios(ctx context.Context) ([]Veterinario, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ID_VETERINARIO, NOMBRE FROM VETERINARIOS ORDER BY NOMBRE`)
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

// ObtenerServicios devuelve servicios, opcionalmente filtrados por categoría.
// Si categoria == "" → devuelve todos ordenados por categoría y nombre.
// Si categoria != "" → devuelve solo los de esa categoría.
// El frontend primero muestra un dropdown de categorías, y al seleccionar una,
// recarga los servicios filtrados con ?categoria=X.
func (s *CitaService) ObtenerServicios(ctx context.Context, categoria string) ([]Servicio, error) {
	var rows *sql.Rows
	var err error
	if categoria != "" {
		rows, err = s.db.QueryContext(ctx,
			`SELECT ID_SERVICIO, NOMBRE_SERVICIO, IFNULL(DESCRIPCION, ' ') AS DESCRIPCION, PRECIO, CATEGORIA
			 FROM SERVICIOS WHERE CATEGORIA = ? ORDER BY NOMBRE_SERVICIO`, categoria)
	} else {
		rows, err = s.db.QueryContext(ctx,
			`SELECT ID_SERVICIO, NOMBRE_SERVICIO, IFNULL(DESCRIPCION, ' ') AS DESCRIPCION, PRECIO, CATEGORIA
			 FROM SERVICIOS ORDER BY CATEGORIA, NOMBRE_SERVICIO`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Servicio
	for rows.Next() {
		var sv Servicio
		if err := rows.Scan(&sv.IDServicio, &sv.NombreServicio, &sv.Descripcion, &sv.Precio, &sv.Categoria); err != nil {
			return nil, err
		}
		list = append(list, sv)
	}
	return list, rows.Err()
}

// Agendar es la operación más compleja del sistema. Crea una cita completa en
// una transacción atómica: valida todo → crea cita → registra servicios →
// descuenta stock de productos → genera factura → vincula todo.
//
// PARÁMETROS (vienen del frontend vía handlers.go):
//   - idCliente: de la sesión (cookie), no del body — seguridad
//   - idMascota, idVeterinario: del formulario
//   - fechaCita: "2025-06-15 10:00" (fecha + hora concatenadas por el handler)
//   - servicios: "1,3,5" (IDs separados por coma, concatenados por el handler)
//
// RETORNA: idFactura — el frontend lo usa para crear la sesión de Stripe.
func (s *CitaService) Agendar(ctx context.Context, idCliente, idMascota, idVeterinario int64, fechaCita, servicios string) (int64, error) {
	// === PASO 1: Parsear la lista de servicios (string "1,3,5" → []int64) ===
	parts := strings.Split(servicios, ",")
	var svcIDs []int64
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return 0, errors.New("ORA-20007: El servicio solicitado no es válido")
		}
		svcIDs = append(svcIDs, id)
	}
	if len(svcIDs) == 0 {
		return 0, errors.New("ORA-20007: Debe seleccionar al menos un servicio")
	}

	// === PASO 2: Abrir transacción — todo o nada ===
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Se ejecuta si hay error; no-op si ya se hizo Commit

	// === PASO 3: Cadena de validaciones (ORA-20001 a ORA-20006) ===
	// Cada validación verifica existencia/pertenencia antes de crear datos.
	var count int

	// 3a. ¿El cliente existe en la BD?
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM CLIENTES WHERE ID_CLIENTE = ?`, idCliente).Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("ORA-20001: El cliente no existe")
	}

	// 3b. ¿La mascota existe?
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM MASCOTAS WHERE ID_MASCOTA = ?`, idMascota).Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("ORA-20002: La mascota no existe")
	}

	// 3c. ¿La mascota pertenece a ESTE cliente? (seguridad: evita agendar mascota ajena)
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM MASCOTAS WHERE ID_MASCOTA = ? AND ID_CLIENTE = ?`, idMascota, idCliente).Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("ORA-20003: La mascota no pertenece al cliente")
	}

	// 3d. ¿La mascota ya tiene cita activa el mismo día?
	if err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM CITAS
		 WHERE ID_MASCOTA = ? AND DATE(FECHA_CITA) = DATE(STR_TO_DATE(?, '%Y-%m-%d %H:%i')) AND ESTADO = 'Activa'`,
		idMascota, fechaCita).Scan(&count); err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("ORA-20004: La mascota ya tiene una cita activa a la misma hora")
	}

	// 3e. ¿El veterinario existe?
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM VETERINARIOS WHERE ID_VETERINARIO = ?`, idVeterinario).Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("ORA-20005: El veterinario no existe")
	}

	// 3f. ¿El veterinario está libre en esa fecha/hora?
	if err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM CITAS
		 WHERE ID_VETERINARIO = ? AND FECHA_CITA = STR_TO_DATE(?, '%Y-%m-%d %H:%i') AND ESTADO = 'Activa'`,
		idVeterinario, fechaCita).Scan(&count); err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("ORA-20006: El veterinario no está disponible en esa fecha y hora")
	}

	// === PASO 4: Crear la cita (estado inicial = 'Activa') ===
	result, err := tx.ExecContext(ctx,
		`INSERT INTO CITAS (ID_MASCOTA, ID_VETERINARIO, FECHA_CITA, ESTADO)
		 VALUES (?, ?, STR_TO_DATE(?, '%Y-%m-%d %H:%i'), 'Activa')`,
		idMascota, idVeterinario, fechaCita)
	if err != nil {
		return 0, err
	}
	idCita, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// === PASO 5: Registrar cada servicio y descontar stock de productos ===
	for _, svcID := range svcIDs {
		// 5a. Verificar que el servicio existe
		if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM SERVICIOS WHERE ID_SERVICIO = ?`, svcID).Scan(&count); err != nil {
			return 0, err
		}
		if count == 0 {
			return 0, errors.New("ORA-20007: El servicio solicitado no es válido")
		}

		// 5b. Crear relación CITAS_SERVICIOS (tabla puente N:N)
		_, err = tx.ExecContext(ctx,
			`INSERT INTO CITAS_SERVICIOS (ID_CITA, ID_SERVICIO) VALUES (?, ?)`,
			idCita, svcID)
		if err != nil {
			return 0, err
		}

		// 5c. Verificar y descontar stock de productos asociados al servicio.
		// SERVICIOS_PRODUCTOS indica qué productos y cuántas unidades consume cada servicio.
		// Ejemplo: "Baño completo" consume 2 unidades de "Shampoo antipulgas".
		prodRows, err := tx.QueryContext(ctx,
			`SELECT sp.ID_PRODUCTO, sp.UNIDADES_PRODUCTO, p.STOCK
			 FROM SERVICIOS_PRODUCTOS sp
			 JOIN PRODUCTOS p ON sp.ID_PRODUCTO = p.ID_PRODUCTO
			 WHERE sp.ID_SERVICIO = ?`, svcID)
		if err != nil {
			return 0, err
		}
		type prodInfo struct {
			idProducto int64
			unidades   int
			stock      int
		}
		var prods []prodInfo
		for prodRows.Next() {
			var pi prodInfo
			if err = prodRows.Scan(&pi.idProducto, &pi.unidades, &pi.stock); err != nil {
				prodRows.Close()
				return 0, err
			}
			prods = append(prods, pi)
		}
		prodRows.Close()

		for _, pi := range prods {
			if pi.stock < pi.unidades {
				return 0, errors.New("ORA-20008: No hay suficiente stock para el producto")
			}
			_, err = tx.ExecContext(ctx,
				`UPDATE PRODUCTOS SET STOCK = STOCK - ? WHERE ID_PRODUCTO = ?`,
				pi.unidades, pi.idProducto)
			if err != nil {
				return 0, err
			}
		}
	}

	// === PASO 6: Calcular total de la cita (suma de precios de todos los servicios) ===
	var total float64
	if err = tx.QueryRowContext(ctx,
		`SELECT IFNULL(SUM(s.PRECIO), 0)
		 FROM CITAS_SERVICIOS cs
		 JOIN SERVICIOS s ON cs.ID_SERVICIO = s.ID_SERVICIO
		 WHERE cs.ID_CITA = ?`, idCita).Scan(&total); err != nil {
		return 0, err
	}

	// === PASO 7: Crear factura con el total calculado ===
	factResult, err := tx.ExecContext(ctx,
		`INSERT INTO FACTURAS (TOTAL, FECHA_FACTURA) VALUES (?, NOW())`, total)
	if err != nil {
		return 0, err
	}
	idFactura, err := factResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	// === PASO 8: Vincular la factura con los registros de CITAS_SERVICIOS ===
	// CITAS_SERVICIOS.FACTURAS_ID_FACTURA enlaza la tabla puente con la factura.
	_, err = tx.ExecContext(ctx,
		`UPDATE CITAS_SERVICIOS SET FACTURAS_ID_FACTURA = ? WHERE ID_CITA = ?`,
		idFactura, idCita)
	if err != nil {
		return 0, err
	}

	// === PASO 9: Commit — si todo fue exitoso, los cambios se persisten ===
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return idFactura, nil
}

// Cancelar elimina una cita y sus detalles. Solo permite cancelar citas propias
// (verifica que la mascota de la cita pertenezca al cliente).
// Primero borra DETALLE_CITAS (FK), luego la CITA misma.
// Retorna true si se eliminó al menos una fila (cita encontrada y autorizada).
func (s *CitaService) Cancelar(ctx context.Context, idCita, idCliente int64) (bool, error) {
	// Paso 1: Borrar detalles de cita (FK requiere borrar primero los hijos)
	// El subquery verifica que la cita pertenece al cliente vía MASCOTAS.
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM DETALLE_CITAS
		 WHERE ID_CITA = ? AND ID_CITA IN (
		     SELECT C.ID_CITA FROM CITAS C
		     JOIN MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
		     WHERE M.ID_CLIENTE = ? AND C.ID_CITA = ?
		 )`,
		idCita, idCliente, idCita,
	)
	if err != nil {
		return false, err
	}

	// Paso 2: Borrar la cita (solo si la mascota es del cliente)
	result, err := s.db.ExecContext(ctx,
		`DELETE FROM CITAS
		 WHERE ID_CITA = ? AND ID_MASCOTA IN (
		     SELECT ID_MASCOTA FROM MASCOTAS WHERE ID_CLIENTE = ?
		 )`,
		idCita, idCliente,
	)
	if err != nil {
		return false, err
	}

	n, _ := result.RowsAffected()
	return n > 0, nil
}
