-- =============================================================================
-- PATITAS AL RESCATE — Script unificado de base de datos (MariaDB)
-- =============================================================================
--
-- Este script crea desde cero la base de datos completa:
--   1. Base de datos y usuario de aplicación
--   2. Tablas de usuarios (CLIENTES, MASCOTAS, VETERINARIOS, USUARIOS)
--   3. Tablas de catálogo (SERVICIOS, PRODUCTOS, SERVICIOS_PRODUCTOS)
--   4. Tablas operativas (CITAS, FACTURAS, CITAS_SERVICIOS)
--   5. Funciones y procedimientos almacenados
--   6. Datos semilla (veterinarios, servicios, productos, admin)
--
-- Uso:
--   mysql -u root -p < patitas_rescate.sql
--
-- Requisitos:
--   MariaDB 10.5+ | MySQL 8.0+
--   Ejecutar con usuario root o con privilegios CREATE DATABASE / CREATE USER
-- =============================================================================


-- ─────────────────────────────────────────────────────────────────────────────
-- 1. BASE DE DATOS Y USUARIO DE APLICACIÓN
-- ─────────────────────────────────────────────────────────────────────────────

CREATE DATABASE IF NOT EXISTS patitas_rescate
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_general_ci;

-- Usuario de aplicación (el backend Go se conecta con estas credenciales)
CREATE USER IF NOT EXISTS 'Progra_PAR'@'%' IDENTIFIED BY 'PrograPAR_2026';
GRANT SELECT, INSERT, UPDATE, DELETE ON patitas_rescate.* TO 'Progra_PAR'@'%';
FLUSH PRIVILEGES;

USE patitas_rescate;


-- ─────────────────────────────────────────────────────────────────────────────
-- 2. TABLAS DE USUARIOS
-- ─────────────────────────────────────────────────────────────────────────────

-- Clientes: datos personales de los dueños de mascotas.
CREATE TABLE IF NOT EXISTS CLIENTES (
    ID_CLIENTE          INT AUTO_INCREMENT PRIMARY KEY,
    DIDENTIDAD_CLIENTE  VARCHAR(20)  NOT NULL,
    NOMBRE              VARCHAR(100) NOT NULL,
    APELLIDO            VARCHAR(100) NOT NULL,
    EMAIL               VARCHAR(100) NOT NULL UNIQUE,
    TELEFONO            VARCHAR(15)  NOT NULL,
    DIRECCION           VARCHAR(255),
    FECHA_REGISTRO      DATETIME     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clientes_didentidad ON CLIENTES(DIDENTIDAD_CLIENTE);

-- Veterinarios: profesionales que atienden citas.
-- ROL aquí es el rol veterinario (e.g. 'Estético'), NO el rol de acceso al sistema.
CREATE TABLE IF NOT EXISTS VETERINARIOS (
    ID_VETERINARIO          INT AUTO_INCREMENT PRIMARY KEY,
    DIDENTIDAD_VETERINARIO  VARCHAR(20)  NOT NULL,
    NOMBRE                  VARCHAR(100) NOT NULL,
    ESPECIALIDAD            VARCHAR(100),
    TELEFONO                VARCHAR(15),
    CORREO                  VARCHAR(100) UNIQUE,
    ROL                     VARCHAR(20)
);

-- Mascotas: asociadas a un cliente (relación 1:N).
CREATE TABLE IF NOT EXISTS MASCOTAS (
    ID_MASCOTA  INT AUTO_INCREMENT PRIMARY KEY,
    ID_CLIENTE  INT          NOT NULL,
    NOMBRE      VARCHAR(100) NOT NULL,
    ESPECIE     VARCHAR(50)  NOT NULL,
    RAZA        VARCHAR(50),
    MESES       INT CHECK (MESES >= 0),
    CONSTRAINT fk_mascotas_clientes
        FOREIGN KEY (ID_CLIENTE) REFERENCES CLIENTES(ID_CLIENTE)
);

CREATE INDEX idx_mascotas_cliente ON MASCOTAS(ID_CLIENTE);

-- Usuarios: credenciales de autenticación (uno por cliente).
-- ROL: 0 = Admin, 1 = Cliente (default), 2 = Veterinario.
-- ID_VETERINARIO: nullable — solo se usa si ROL = 2. Vincula con VETERINARIOS
-- para que el backend sepa cuál veterinario es. Los admins y clientes tienen NULL.
CREATE TABLE IF NOT EXISTS USUARIOS (
    ID_USUARIO      INT AUTO_INCREMENT PRIMARY KEY,
    ID_CLIENTE      INT          NOT NULL,
    CORREO          VARCHAR(100) NOT NULL UNIQUE,
    CONTRASENA      VARCHAR(255) NOT NULL,
    ROL             TINYINT      NOT NULL DEFAULT 1,
    ID_VETERINARIO  INT          NULL,
    CONSTRAINT fk_usuarios_clientes
        FOREIGN KEY (ID_CLIENTE) REFERENCES CLIENTES(ID_CLIENTE),
    CONSTRAINT fk_usuarios_veterinarios
        FOREIGN KEY (ID_VETERINARIO) REFERENCES VETERINARIOS(ID_VETERINARIO)
);

CREATE INDEX idx_usuarios_cliente ON USUARIOS(ID_CLIENTE);


-- ─────────────────────────────────────────────────────────────────────────────
-- 3. TABLAS DE CATÁLOGO (SERVICIOS Y PRODUCTOS)
-- ─────────────────────────────────────────────────────────────────────────────

-- Servicios: catálogo de lo que ofrece la veterinaria.
-- CATEGORIA agrupa servicios en el frontend (e.g. 'Estética', 'Salud', 'Recreación').
CREATE TABLE IF NOT EXISTS SERVICIOS (
    ID_SERVICIO       INT AUTO_INCREMENT PRIMARY KEY,
    NOMBRE_SERVICIO   VARCHAR(100) NOT NULL,
    DESCRIPCION       TEXT,
    PRECIO            DECIMAL(10,2) NOT NULL CHECK (PRECIO >= 0),
    DURACION_MINUTOS  INT CHECK (DURACION_MINUTOS > 0),
    CATEGORIA         VARCHAR(100)
);

-- Productos: consumibles usados durante los servicios.
CREATE TABLE IF NOT EXISTS PRODUCTOS (
    ID_PRODUCTO      INT AUTO_INCREMENT PRIMARY KEY,
    NOMBRE_PRODUCTO  VARCHAR(100)  NOT NULL,
    CATEGORIA        VARCHAR(100),
    PRECIO           DECIMAL(10,2) NOT NULL CHECK (PRECIO >= 0),
    STOCK            INT           NOT NULL CHECK (STOCK >= 0)
);

-- Relación N:N entre servicios y productos.
-- Indica qué productos consume cada servicio y cuántas unidades.
CREATE TABLE IF NOT EXISTS SERVICIOS_PRODUCTOS (
    ID_PRODUCTO_SERVICIO  INT AUTO_INCREMENT PRIMARY KEY,
    ID_PRODUCTO           INT NOT NULL,
    ID_SERVICIO           INT NOT NULL,
    UNIDADES_PRODUCTO     INT NOT NULL CHECK (UNIDADES_PRODUCTO > 0),
    CANTIDAD_CONSUMIDA    INT NOT NULL CHECK (CANTIDAD_CONSUMIDA > 0),
    CONSTRAINT fk_servprod_producto
        FOREIGN KEY (ID_PRODUCTO) REFERENCES PRODUCTOS(ID_PRODUCTO),
    CONSTRAINT fk_servprod_servicio
        FOREIGN KEY (ID_SERVICIO) REFERENCES SERVICIOS(ID_SERVICIO)
);

CREATE INDEX idx_servprod_producto  ON SERVICIOS_PRODUCTOS(ID_PRODUCTO);
CREATE INDEX idx_servprod_servicio  ON SERVICIOS_PRODUCTOS(ID_SERVICIO);


-- ─────────────────────────────────────────────────────────────────────────────
-- 4. TABLAS OPERATIVAS (CITAS, FACTURAS, RELACIONES)
-- ─────────────────────────────────────────────────────────────────────────────

-- Citas: agendadas por clientes, asignadas a un veterinario.
CREATE TABLE IF NOT EXISTS CITAS (
    ID_CITA         INT AUTO_INCREMENT PRIMARY KEY,
    ID_MASCOTA      INT         NOT NULL,
    ID_VETERINARIO  INT         NOT NULL,
    FECHA_CITA      DATETIME    NOT NULL,
    ESTADO          VARCHAR(20) NOT NULL CHECK (ESTADO IN ('Activa', 'Completada', 'Cancelada')),
    CONSTRAINT fk_citas_mascotas
        FOREIGN KEY (ID_MASCOTA) REFERENCES MASCOTAS(ID_MASCOTA),
    CONSTRAINT fk_citas_veterinarios
        FOREIGN KEY (ID_VETERINARIO) REFERENCES VETERINARIOS(ID_VETERINARIO)
);

CREATE INDEX idx_citas_mascota      ON CITAS(ID_MASCOTA);
CREATE INDEX idx_citas_veterinario  ON CITAS(ID_VETERINARIO);
CREATE INDEX idx_citas_estado_fecha ON CITAS(ESTADO, FECHA_CITA);

-- Facturas: generadas al agendar una cita.
-- ESTADO: 'Pendiente' (recién creada) → 'Pagada' (después de Stripe checkout).
-- STRIPE_SESSION_ID: almacena el ID de la sesión de Stripe para verificar el pago.
CREATE TABLE IF NOT EXISTS FACTURAS (
    ID_FACTURA         INT AUTO_INCREMENT PRIMARY KEY,
    FECHA_FACTURA      DATETIME      NOT NULL DEFAULT NOW(),
    TOTAL              DECIMAL(10,2) NOT NULL CHECK (TOTAL >= 0),
    ESTADO             VARCHAR(20)   NOT NULL DEFAULT 'Pendiente',
    STRIPE_SESSION_ID  VARCHAR(255)  NULL
);

-- Tabla puente N:N entre CITAS y SERVICIOS.
-- FACTURAS_ID_FACTURA vincula la cita-servicio con su factura.
CREATE TABLE IF NOT EXISTS CITAS_SERVICIOS (
    ID_CITA_SERVICIO    INT AUTO_INCREMENT PRIMARY KEY,
    ID_CITA             INT NOT NULL,
    ID_SERVICIO         INT NOT NULL,
    FACTURAS_ID_FACTURA INT,
    CONSTRAINT fk_citaserv_citas
        FOREIGN KEY (ID_CITA) REFERENCES CITAS(ID_CITA),
    CONSTRAINT fk_citaserv_servicios
        FOREIGN KEY (ID_SERVICIO) REFERENCES SERVICIOS(ID_SERVICIO),
    CONSTRAINT fk_citaserv_facturas
        FOREIGN KEY (FACTURAS_ID_FACTURA) REFERENCES FACTURAS(ID_FACTURA)
);

CREATE INDEX idx_citaserv_cita     ON CITAS_SERVICIOS(ID_CITA);
CREATE INDEX idx_citaserv_servicio ON CITAS_SERVICIOS(ID_SERVICIO);
CREATE INDEX idx_citaserv_factura  ON CITAS_SERVICIOS(FACTURAS_ID_FACTURA);


-- ─────────────────────────────────────────────────────────────────────────────
-- 5. FUNCIONES Y PROCEDIMIENTOS ALMACENADOS
-- ─────────────────────────────────────────────────────────────────────────────

DELIMITER //

-- Verificar si un cliente existe
CREATE OR REPLACE FUNCTION existeCliente(p_ID_Cliente INT) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM CLIENTES WHERE ID_CLIENTE = p_ID_Cliente;
    RETURN v_count > 0;
END //

-- Verificar si una mascota existe
CREATE OR REPLACE FUNCTION existeMascota(p_ID_Mascota INT) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM MASCOTAS WHERE ID_MASCOTA = p_ID_Mascota;
    RETURN v_count > 0;
END //

-- Verificar si una mascota pertenece a un cliente
CREATE OR REPLACE FUNCTION mascotaPerteneceACliente(p_ID_Mascota INT, p_ID_Cliente INT) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM MASCOTAS
    WHERE ID_MASCOTA = p_ID_Mascota AND ID_CLIENTE = p_ID_Cliente;
    RETURN v_count > 0;
END //

-- Verificar si una mascota tiene cita activa en la misma fecha
CREATE OR REPLACE FUNCTION mascotaTieneCitaActivaMismaFecha(p_ID_Mascota INT, p_Fecha_Cita DATETIME) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM CITAS
    WHERE ID_MASCOTA = p_ID_Mascota
      AND DATE(FECHA_CITA) = DATE(p_Fecha_Cita)
      AND ESTADO = 'Activa';
    RETURN v_count > 0;
END //

-- Verificar si un veterinario existe
CREATE OR REPLACE FUNCTION existeVeterinario(p_ID_Veterinario INT) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM VETERINARIOS WHERE ID_VETERINARIO = p_ID_Veterinario;
    RETURN v_count > 0;
END //

-- Verificar disponibilidad de un veterinario
CREATE OR REPLACE FUNCTION veterinarioDisponible(p_ID_Veterinario INT, p_Fecha_Cita DATETIME) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM CITAS
    WHERE ID_VETERINARIO = p_ID_Veterinario
      AND FECHA_CITA = p_Fecha_Cita
      AND ESTADO = 'Activa';
    RETURN v_count = 0;
END //

-- Verificar si un servicio existe
CREATE OR REPLACE FUNCTION existeServicio(p_ID_Servicio INT) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM SERVICIOS WHERE ID_SERVICIO = p_ID_Servicio;
    RETURN v_count > 0;
END //

-- Verificar si un producto existe
CREATE OR REPLACE FUNCTION existeProducto(p_ID_Producto INT) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM PRODUCTOS WHERE ID_PRODUCTO = p_ID_Producto;
    RETURN v_count > 0;
END //

-- Calcular el total de servicios de una cita
CREATE OR REPLACE FUNCTION calcularTotalServicios(p_ID_Cita INT) RETURNS DECIMAL(10,2)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_Total DECIMAL(10,2);
    SELECT IFNULL(SUM(s.PRECIO), 0) INTO v_Total
    FROM CITAS_SERVICIOS cs
    JOIN SERVICIOS s ON cs.ID_SERVICIO = s.ID_SERVICIO
    WHERE cs.ID_CITA = p_ID_Cita;
    RETURN v_Total;
END //

-- Verificar si una cédula ya está registrada
CREATE OR REPLACE FUNCTION existeCedula(p_didentidad VARCHAR(20)) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM CLIENTES WHERE DIDENTIDAD_CLIENTE = p_didentidad;
    RETURN v_count > 0;
END //

-- Verificar si un correo ya está registrado
CREATE OR REPLACE FUNCTION existeCorreo(p_email VARCHAR(100)) RETURNS TINYINT(1)
DETERMINISTIC READS SQL DATA
BEGIN
    DECLARE v_count INT;
    SELECT COUNT(*) INTO v_count FROM USUARIOS WHERE CORREO = p_email;
    RETURN v_count > 0;
END //

-- Procedimiento para actualizar stock de un producto
CREATE OR REPLACE PROCEDURE actualizarStock(
    IN p_ID_Producto INT,
    IN p_Cantidad    INT
)
BEGIN
    UPDATE PRODUCTOS SET STOCK = STOCK - p_Cantidad
    WHERE ID_PRODUCTO = p_ID_Producto;

    IF ROW_COUNT() = 0 THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'No se pudo actualizar el stock para el producto';
    END IF;
END //

-- Procedimiento para registrar un cliente nuevo (transaccional)
CREATE OR REPLACE PROCEDURE registrarCliente(
    IN  p_didentidad_cliente VARCHAR(20),
    IN  p_nombre             VARCHAR(100),
    IN  p_apellido           VARCHAR(100),
    IN  p_email              VARCHAR(100),
    IN  p_telefono           VARCHAR(15),
    IN  p_direccion          VARCHAR(255),
    IN  p_contrasena         VARCHAR(255),
    OUT p_ID_Cliente         INT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    IF existeCedula(p_didentidad_cliente) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'La cédula o documento de identidad ya está registrada';
    END IF;

    IF existeCorreo(p_email) THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'El correo electrónico ya está registrado';
    END IF;

    INSERT INTO CLIENTES (DIDENTIDAD_CLIENTE, NOMBRE, APELLIDO, EMAIL, TELEFONO, DIRECCION, FECHA_REGISTRO)
    VALUES (p_didentidad_cliente, p_nombre, p_apellido, p_email, p_telefono, p_direccion, NOW());

    SET p_ID_Cliente = LAST_INSERT_ID();

    INSERT INTO USUARIOS (ID_CLIENTE, CORREO, CONTRASENA)
    VALUES (p_ID_Cliente, p_email, p_contrasena);

    COMMIT;
END //

-- Procedimiento para agendar una cita (transaccional, con validaciones y stock)
CREATE OR REPLACE PROCEDURE agendarCita(
    IN p_ID_Cliente      INT,
    IN p_ID_Mascota      INT,
    IN p_ID_Veterinario  INT,
    IN p_Fecha_Cita      DATETIME,
    IN p_Servicios       VARCHAR(500)
)
BEGIN
    DECLARE v_ID_Cita       INT;
    DECLARE v_ID_Factura    INT;
    DECLARE v_Total         DECIMAL(10,2);
    DECLARE v_svcID         INT;
    DECLARE v_remaining     VARCHAR(500);
    DECLARE v_token         VARCHAR(20);
    DECLARE v_prod_id       INT;
    DECLARE v_prod_unidades INT;
    DECLARE v_prod_stock    INT;
    DECLARE v_done          INT DEFAULT 0;

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

    -- Validaciones
    IF NOT existeCliente(p_ID_Cliente) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'El cliente no existe';
    END IF;

    IF NOT existeMascota(p_ID_Mascota) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'La mascota no existe';
    END IF;

    IF NOT mascotaPerteneceACliente(p_ID_Mascota, p_ID_Cliente) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'La mascota no pertenece al cliente';
    END IF;

    IF mascotaTieneCitaActivaMismaFecha(p_ID_Mascota, p_Fecha_Cita) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'La mascota ya tiene una cita activa a la misma hora';
    END IF;

    IF NOT existeVeterinario(p_ID_Veterinario) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'El veterinario no existe';
    END IF;

    IF NOT veterinarioDisponible(p_ID_Veterinario, p_Fecha_Cita) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'El veterinario no está disponible en esa fecha y hora';
    END IF;

    -- Crear la cita
    INSERT INTO CITAS (ID_MASCOTA, ID_VETERINARIO, FECHA_CITA, ESTADO)
    VALUES (p_ID_Mascota, p_ID_Veterinario, p_Fecha_Cita, 'Activa');

    SET v_ID_Cita = LAST_INSERT_ID();

    -- Recorrer servicios (cadena separada por comas: "1,3,5")
    SET v_remaining = p_Servicios;

    WHILE LENGTH(v_remaining) > 0 DO
        IF LOCATE(',', v_remaining) > 0 THEN
            SET v_token    = TRIM(SUBSTRING_INDEX(v_remaining, ',', 1));
            SET v_remaining = TRIM(SUBSTRING(v_remaining, LOCATE(',', v_remaining) + 1));
        ELSE
            SET v_token    = TRIM(v_remaining);
            SET v_remaining = '';
        END IF;

        SET v_svcID = CAST(v_token AS UNSIGNED);

        IF NOT existeServicio(v_svcID) THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'El servicio solicitado no es válido';
        END IF;

        INSERT INTO CITAS_SERVICIOS (ID_CITA, ID_SERVICIO)
        VALUES (v_ID_Cita, v_svcID);

        -- Descontar stock de productos asociados
        SET v_done = 0;
        OPEN cur_productos;
        loop_productos: LOOP
            FETCH cur_productos INTO v_prod_id, v_prod_unidades, v_prod_stock;
            IF v_done THEN LEAVE loop_productos; END IF;

            IF v_prod_stock < v_prod_unidades THEN
                CLOSE cur_productos;
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'No hay suficiente stock para el producto';
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

    UPDATE CITAS_SERVICIOS
    SET FACTURAS_ID_FACTURA = v_ID_Factura
    WHERE ID_CITA = v_ID_Cita;

    COMMIT;
END //

DELIMITER ;


-- ─────────────────────────────────────────────────────────────────────────────
-- 6. DATOS SEMILLA
-- ─────────────────────────────────────────────────────────────────────────────

-- Veterinarios de ejemplo
INSERT INTO VETERINARIOS (ID_VETERINARIO, DIDENTIDAD_VETERINARIO, NOMBRE, ESPECIALIDAD, TELEFONO, CORREO, ROL) VALUES
    (1, 'VET123456789', 'Dr. Juan Pérez',     'Estética Canina',  '88881234', 'juan.perez@veterinaria.com',    'Estético'),
    (2, 'VET987654321', 'Dra. María López',    'Estética Felina',  '88885678', 'maria.lopez@veterinaria.com',   'Estético'),
    (3, 'VET456789123', 'Dr. Carlos Ramírez',  'Estética General', '88889012', 'carlos.ramirez@veterinaria.com','Estético');

-- Servicios con categorías (usadas por el frontend para filtrar)
INSERT INTO SERVICIOS (ID_SERVICIO, NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA) VALUES
    (1, 'Corte Mini',     'Corte de pelo para perros pequeños',             15000.00, 30, 'Estética'),
    (2, 'Corte Mediano',  'Corte de pelo para perros medianos',             25000.00, 45, 'Estética'),
    (3, 'Corte Grande',   'Corte de pelo para perros grandes',              35000.00, 60, 'Estética'),
    (4, 'Corte de Uñas',  'Corte de uñas para perros de cualquier tamaño', 10000.00, 20, 'Estética');

-- Productos consumibles
INSERT INTO PRODUCTOS (ID_PRODUCTO, NOMBRE_PRODUCTO, CATEGORIA, PRECIO, STOCK) VALUES
    (1, 'Shampoo burbuja',        'Higiene', 1500.00, 100),
    (2, 'Acondicionador burbuja', 'Higiene', 1350.00, 100),
    (3, 'Desinfectante',          'Higiene',  400.00,  50);

-- Asociaciones: qué productos consume cada servicio
-- Corte Mini (servicio 1)
INSERT INTO SERVICIOS_PRODUCTOS (ID_PRODUCTO, ID_SERVICIO, UNIDADES_PRODUCTO, CANTIDAD_CONSUMIDA) VALUES
    (1, 1, 1, 1),
    (2, 1, 1, 1),
    (3, 1, 1, 1);

-- Corte Mediano (servicio 2)
INSERT INTO SERVICIOS_PRODUCTOS (ID_PRODUCTO, ID_SERVICIO, UNIDADES_PRODUCTO, CANTIDAD_CONSUMIDA) VALUES
    (1, 2, 2, 2),
    (2, 2, 2, 2),
    (3, 2, 1, 1);

-- Corte Grande (servicio 3)
INSERT INTO SERVICIOS_PRODUCTOS (ID_PRODUCTO, ID_SERVICIO, UNIDADES_PRODUCTO, CANTIDAD_CONSUMIDA) VALUES
    (1, 3, 3, 3),
    (2, 3, 3, 3),
    (3, 3, 1, 1);

-- Usuario administrador de ejemplo
-- Contraseña: Admin123! (hash bcrypt generado con cost=10)
-- Para generar otro hash: htpasswd -nbBC 10 "" 'TuPassword' | cut -d: -f2
INSERT INTO CLIENTES (ID_CLIENTE, DIDENTIDAD_CLIENTE, NOMBRE, APELLIDO, EMAIL, TELEFONO, DIRECCION)
VALUES (1, 'ADM000000000', 'Admin', 'Sistema', 'admin@patitas.com', '00000000', 'Sistema');

INSERT INTO USUARIOS (ID_CLIENTE, CORREO, CONTRASENA, ROL)
VALUES (1, 'admin@patitas.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 0);
