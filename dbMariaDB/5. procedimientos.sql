-- Procedimientos y funciones para MariaDB
-- Ejecutar después de los scripts 1-4
USE patitas_rescate;

DELIMITER //

-- Función para verificar si un cliente existe
CREATE OR REPLACE FUNCTION existeCliente(p_ID_Cliente INT) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM CLIENTES
    WHERE ID_CLIENTE = p_ID_Cliente;
    RETURN v_count > 0;
END //

-- Función para verificar si una mascota existe
CREATE OR REPLACE FUNCTION existeMascota(p_ID_Mascota INT) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM MASCOTAS
    WHERE ID_MASCOTA = p_ID_Mascota;
    RETURN v_count > 0;
END //

-- Función para verificar si una mascota pertenece a un cliente
CREATE OR REPLACE FUNCTION mascotaPerteneceACliente(p_ID_Mascota INT, p_ID_Cliente INT) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM MASCOTAS
    WHERE ID_MASCOTA = p_ID_Mascota AND ID_CLIENTE = p_ID_Cliente;
    RETURN v_count > 0;
END //

-- Función para verificar si una mascota tiene una cita activa en la misma fecha
CREATE OR REPLACE FUNCTION mascotaTieneCitaActivaMismaFecha(p_ID_Mascota INT, p_Fecha_Cita DATETIME) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM CITAS
    WHERE ID_MASCOTA = p_ID_Mascota
      AND DATE(FECHA_CITA) = DATE(p_Fecha_Cita)
      AND ESTADO = 'Activa';
    RETURN v_count > 0;
END //

-- Función para verificar si un veterinario existe
CREATE OR REPLACE FUNCTION existeVeterinario(p_ID_Veterinario INT) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM VETERINARIOS
    WHERE ID_VETERINARIO = p_ID_Veterinario;
    RETURN v_count > 0;
END //

-- Función para verificar si un veterinario está disponible
CREATE OR REPLACE FUNCTION veterinarioDisponible(p_ID_Veterinario INT, p_Fecha_Cita DATETIME) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM CITAS
    WHERE ID_VETERINARIO = p_ID_Veterinario
      AND FECHA_CITA = p_Fecha_Cita
      AND ESTADO = 'Activa';
    RETURN v_count = 0;
END //

-- Función para verificar si un servicio existe
CREATE OR REPLACE FUNCTION existeServicio(p_ID_Servicio INT) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM SERVICIOS
    WHERE ID_SERVICIO = p_ID_Servicio;
    RETURN v_count > 0;
END //

-- Función para verificar si un producto existe
CREATE OR REPLACE FUNCTION existeProducto(p_ID_Producto INT) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM PRODUCTOS
    WHERE ID_PRODUCTO = p_ID_Producto;
    RETURN v_count > 0;
END //

-- Función para calcular el total de servicios de una cita
CREATE OR REPLACE FUNCTION calcularTotalServicios(p_ID_Cita INT) RETURNS DECIMAL(10,2)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_Total DECIMAL(10,2);
    SELECT IFNULL(SUM(s.PRECIO), 0) INTO v_Total
    FROM CITAS_SERVICIOS cs
    JOIN SERVICIOS s ON cs.ID_SERVICIO = s.ID_SERVICIO
    WHERE cs.ID_CITA = p_ID_Cita;
    RETURN v_Total;
END //

-- Procedimiento para actualizar stock
CREATE OR REPLACE PROCEDURE actualizarStock(
    IN p_ID_Producto INT,
    IN p_Cantidad INT
)
BEGIN
    UPDATE PRODUCTOS
    SET STOCK = STOCK - p_Cantidad
    WHERE ID_PRODUCTO = p_ID_Producto;

    IF ROW_COUNT() = 0 THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'No se pudo actualizar el stock para el producto';
    END IF;
END //

-- Función para verificar si existe una cédula
CREATE OR REPLACE FUNCTION existeCedula(p_didentidad VARCHAR(20)) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM CLIENTES
    WHERE DIDENTIDAD_CLIENTE = p_didentidad;
    RETURN v_count > 0;
END //

-- Función para verificar si existe un correo
CREATE OR REPLACE FUNCTION existeCorreo(p_email VARCHAR(100)) RETURNS TINYINT(1)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count
    FROM USUARIOS
    WHERE CORREO = p_email;
    RETURN v_count > 0;
END //

-- Procedimiento para registrar un cliente
CREATE OR REPLACE PROCEDURE registrarCliente(
    IN p_didentidad_cliente VARCHAR(20),
    IN p_nombre VARCHAR(100),
    IN p_apellido VARCHAR(100),
    IN p_email VARCHAR(100),
    IN p_telefono VARCHAR(15),
    IN p_direccion VARCHAR(255),
    IN p_contrasena VARCHAR(255),
    OUT p_ID_Cliente INT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    -- Validar cédula o documento de identidad
    IF existeCedula(p_didentidad_cliente) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'La cédula o documento de identidad ya está registrada';
    END IF;

    -- Validar correo electrónico
    IF existeCorreo(p_email) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'El correo electrónico ya está registrado';
    END IF;

    -- Insertar cliente
    INSERT INTO CLIENTES (DIDENTIDAD_CLIENTE, NOMBRE, APELLIDO, EMAIL, TELEFONO, DIRECCION, FECHA_REGISTRO)
    VALUES (p_didentidad_cliente, p_nombre, p_apellido, p_email, p_telefono, p_direccion, NOW());

    SET p_ID_Cliente = LAST_INSERT_ID();

    -- Insertar usuario
    INSERT INTO USUARIOS (ID_CLIENTE, CORREO, CONTRASENA)
    VALUES (p_ID_Cliente, p_email, p_contrasena);

    COMMIT;
END //

-- Procedimiento para agendar una cita
CREATE OR REPLACE PROCEDURE agendarCita(
    IN p_ID_Cliente INT,
    IN p_ID_Mascota INT,
    IN p_ID_Veterinario INT,
    IN p_Fecha_Cita DATETIME,
    IN p_Servicios VARCHAR(500)
)
BEGIN
    DECLARE v_ID_Cita INT;
    DECLARE v_ID_Factura INT;
    DECLARE v_Total DECIMAL(10,2);
    DECLARE v_svcID INT;
    DECLARE v_remaining VARCHAR(500);
    DECLARE v_token VARCHAR(20);
    DECLARE v_prod_id INT;
    DECLARE v_prod_unidades INT;
    DECLARE v_prod_stock INT;
    DECLARE v_done INT DEFAULT 0;

    DECLARE cur_productos CURSOR FOR
        SELECT sp.ID_PRODUCTO, sp.UNIDADES_PRODUCTO, p.STOCK
        FROM SERVICIOS_PRODUCTOS sp
        JOIN PRODUCTOS p ON sp.ID_PRODUCTO = p.ID_PRODUCTO
        WHERE sp.ID_SERVICIO = v_svcID;

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_done = 1;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    -- Validar cliente
    IF NOT existeCliente(p_ID_Cliente) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'El cliente no existe';
    END IF;

    -- Validar mascota
    IF NOT existeMascota(p_ID_Mascota) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'La mascota no existe';
    END IF;

    IF NOT mascotaPerteneceACliente(p_ID_Mascota, p_ID_Cliente) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'La mascota no pertenece al cliente';
    END IF;

    -- Validar si la mascota ya tiene una cita activa en la misma fecha
    IF mascotaTieneCitaActivaMismaFecha(p_ID_Mascota, p_Fecha_Cita) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'La mascota ya tiene una cita activa a la misma hora';
    END IF;

    -- Validar veterinario
    IF NOT existeVeterinario(p_ID_Veterinario) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'El veterinario no existe';
    END IF;

    -- Validar disponibilidad del veterinario
    IF NOT veterinarioDisponible(p_ID_Veterinario, p_Fecha_Cita) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'El veterinario no está disponible en esa fecha y hora';
    END IF;

    -- Crear la cita
    INSERT INTO CITAS (ID_MASCOTA, ID_VETERINARIO, FECHA_CITA, ESTADO)
    VALUES (p_ID_Mascota, p_ID_Veterinario, p_Fecha_Cita, 'Activa');

    SET v_ID_Cita = LAST_INSERT_ID();

    -- Recorrer servicios (cadena separada por comas)
    SET v_remaining = p_Servicios;

    WHILE LENGTH(v_remaining) > 0 DO
        -- Extraer el siguiente token
        IF LOCATE(',', v_remaining) > 0 THEN
            SET v_token = TRIM(SUBSTRING_INDEX(v_remaining, ',', 1));
            SET v_remaining = TRIM(SUBSTRING(v_remaining, LOCATE(',', v_remaining) + 1));
        ELSE
            SET v_token = TRIM(v_remaining);
            SET v_remaining = '';
        END IF;

        SET v_svcID = CAST(v_token AS UNSIGNED);

        -- Validar servicio
        IF NOT existeServicio(v_svcID) THEN
            SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'El servicio solicitado no es válido';
        END IF;

        -- Insertar cita-servicio
        INSERT INTO CITAS_SERVICIOS (ID_CITA, ID_SERVICIO)
        VALUES (v_ID_Cita, v_svcID);

        -- Validar y actualizar stock de productos asociados
        SET v_done = 0;
        OPEN cur_productos;

        loop_productos: LOOP
            FETCH cur_productos INTO v_prod_id, v_prod_unidades, v_prod_stock;
            IF v_done THEN
                LEAVE loop_productos;
            END IF;

            IF v_prod_stock < v_prod_unidades THEN
                CLOSE cur_productos;
                SIGNAL SQLSTATE '45000'
                    SET MESSAGE_TEXT = 'No hay suficiente stock para el producto';
            END IF;

            CALL actualizarStock(v_prod_id, v_prod_unidades);
        END LOOP;

        CLOSE cur_productos;
    END WHILE;

    -- Generar factura
    SET v_Total = calcularTotalServicios(v_ID_Cita);

    INSERT INTO FACTURAS (TOTAL, FECHA_FACTURA)
    VALUES (v_Total, NOW());

    SET v_ID_Factura = LAST_INSERT_ID();

    -- Asociar factura con los servicios de la cita
    UPDATE CITAS_SERVICIOS
    SET FACTURAS_ID_FACTURA = v_ID_Factura
    WHERE ID_CITA = v_ID_Cita;

    COMMIT;
END //

DELIMITER ;

-- Consulta diagnóstica: citas por cliente
-- SELECT c.ID_CITA, c.FECHA_CITA, c.ESTADO, m.NOMBRE AS MASCOTA, v.NOMBRE AS VETERINARIO
-- FROM CITAS c
-- JOIN MASCOTAS m ON c.ID_MASCOTA = m.ID_MASCOTA
-- JOIN VETERINARIOS v ON c.ID_VETERINARIO = v.ID_VETERINARIO
-- WHERE m.ID_CLIENTE = ?;

-- Consulta diagnóstica: facturas por cliente
-- SELECT DISTINCT f.ID_FACTURA, f.TOTAL, f.FECHA_FACTURA,
--        CASE WHEN c.FECHA_CITA < NOW() THEN 'Pagada' ELSE 'Pendiente' END AS ESTADO
-- FROM FACTURAS f
-- JOIN CITAS_SERVICIOS cs ON f.ID_FACTURA = cs.FACTURAS_ID_FACTURA
-- JOIN CITAS c ON cs.ID_CITA = c.ID_CITA
-- JOIN MASCOTAS m ON c.ID_MASCOTA = m.ID_MASCOTA
-- WHERE m.ID_CLIENTE = ?;
