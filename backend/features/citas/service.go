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

func (s *CitaService) ObtenerServicios(ctx context.Context) ([]Servicio, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ID_SERVICIO, NOMBRE_SERVICIO, IFNULL(DESCRIPCION, ' ') AS DESCRIPCION
		 FROM SERVICIOS ORDER BY NOMBRE_SERVICIO`)
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
	// Parsear lista de servicios
	parts := strings.Split(servicios, ",")
	var svcIDs []int64
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return errors.New("ORA-20007: El servicio solicitado no es válido")
		}
		svcIDs = append(svcIDs, id)
	}
	if len(svcIDs) == 0 {
		return errors.New("ORA-20007: Debe seleccionar al menos un servicio")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Validar cliente
	var count int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM CLIENTES WHERE ID_CLIENTE = ?`, idCliente).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return errors.New("ORA-20001: El cliente no existe")
	}

	// Validar mascota
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM MASCOTAS WHERE ID_MASCOTA = ?`, idMascota).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return errors.New("ORA-20002: La mascota no existe")
	}

	// Validar mascota pertenece al cliente
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM MASCOTAS WHERE ID_MASCOTA = ? AND ID_CLIENTE = ?`, idMascota, idCliente).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return errors.New("ORA-20003: La mascota no pertenece al cliente")
	}

	// Validar que la mascota no tenga cita activa el mismo día
	if err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM CITAS
		 WHERE ID_MASCOTA = ? AND DATE(FECHA_CITA) = DATE(STR_TO_DATE(?, '%Y-%m-%d %H:%i')) AND ESTADO = 'Activa'`,
		idMascota, fechaCita).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("ORA-20004: La mascota ya tiene una cita activa a la misma hora")
	}

	// Validar veterinario
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM VETERINARIOS WHERE ID_VETERINARIO = ?`, idVeterinario).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return errors.New("ORA-20005: El veterinario no existe")
	}

	// Validar disponibilidad del veterinario
	if err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM CITAS
		 WHERE ID_VETERINARIO = ? AND FECHA_CITA = STR_TO_DATE(?, '%Y-%m-%d %H:%i') AND ESTADO = 'Activa'`,
		idVeterinario, fechaCita).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("ORA-20006: El veterinario no está disponible en esa fecha y hora")
	}

	// Crear la cita
	result, err := tx.ExecContext(ctx,
		`INSERT INTO CITAS (ID_MASCOTA, ID_VETERINARIO, FECHA_CITA, ESTADO)
		 VALUES (?, ?, STR_TO_DATE(?, '%Y-%m-%d %H:%i'), 'Activa')`,
		idMascota, idVeterinario, fechaCita)
	if err != nil {
		return err
	}
	idCita, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Registrar servicios y actualizar stock
	for _, svcID := range svcIDs {
		// Verificar servicio
		if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM SERVICIOS WHERE ID_SERVICIO = ?`, svcID).Scan(&count); err != nil {
			return err
		}
		if count == 0 {
			return errors.New("ORA-20007: El servicio solicitado no es válido")
		}

		// Insertar cita-servicio
		_, err = tx.ExecContext(ctx,
			`INSERT INTO CITAS_SERVICIOS (ID_CITA, ID_SERVICIO) VALUES (?, ?)`,
			idCita, svcID)
		if err != nil {
			return err
		}

		// Verificar y actualizar stock de productos asociados
		prodRows, err := tx.QueryContext(ctx,
			`SELECT sp.ID_PRODUCTO, sp.UNIDADES_PRODUCTO, p.STOCK
			 FROM SERVICIOS_PRODUCTOS sp
			 JOIN PRODUCTOS p ON sp.ID_PRODUCTO = p.ID_PRODUCTO
			 WHERE sp.ID_SERVICIO = ?`, svcID)
		if err != nil {
			return err
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
				return err
			}
			prods = append(prods, pi)
		}
		prodRows.Close()

		for _, pi := range prods {
			if pi.stock < pi.unidades {
				return errors.New("ORA-20008: No hay suficiente stock para el producto")
			}
			_, err = tx.ExecContext(ctx,
				`UPDATE PRODUCTOS SET STOCK = STOCK - ? WHERE ID_PRODUCTO = ?`,
				pi.unidades, pi.idProducto)
			if err != nil {
				return err
			}
		}
	}

	// Calcular total y crear factura
	var total float64
	if err = tx.QueryRowContext(ctx,
		`SELECT IFNULL(SUM(s.PRECIO), 0)
		 FROM CITAS_SERVICIOS cs
		 JOIN SERVICIOS s ON cs.ID_SERVICIO = s.ID_SERVICIO
		 WHERE cs.ID_CITA = ?`, idCita).Scan(&total); err != nil {
		return err
	}

	factResult, err := tx.ExecContext(ctx,
		`INSERT INTO FACTURAS (TOTAL, FECHA_FACTURA) VALUES (?, NOW())`, total)
	if err != nil {
		return err
	}
	idFactura, err := factResult.LastInsertId()
	if err != nil {
		return err
	}

	// Asociar factura con cita-servicios
	_, err = tx.ExecContext(ctx,
		`UPDATE CITAS_SERVICIOS SET FACTURAS_ID_FACTURA = ? WHERE ID_CITA = ?`,
		idFactura, idCita)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *CitaService) Cancelar(ctx context.Context, idCita, idCliente int64) (bool, error) {
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
