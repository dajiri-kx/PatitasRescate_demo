-- Función para verificar si un cliente existe
CREATE OR REPLACE FUNCTION existeCliente(p_ID_Cliente IN NUMBER) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM usuarios_tablas.clientes
    WHERE id_cliente = p_ID_Cliente;

    RETURN v_count > 0;
END;
/

-- Función para verificar si una mascota existe
CREATE OR REPLACE FUNCTION existeMascota(p_ID_Mascota IN NUMBER) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM usuarios_tablas.mascotas
    WHERE id_mascota = p_ID_Mascota;

    RETURN v_count > 0;
END;
/

-- Función para verificar si una mascota pertenece a un cliente
CREATE OR REPLACE FUNCTION mascotaPerteneceACliente(p_ID_Mascota IN NUMBER, p_ID_Cliente IN NUMBER) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM usuarios_tablas.mascotas
    WHERE id_mascota = p_ID_Mascota AND id_cliente = p_ID_Cliente;

    RETURN v_count > 0;
END;
/

-- Función para verificar si una mascota tiene una cita activa
CREATE OR REPLACE FUNCTION mascotaTieneCitaActivaNowMismoServicio(p_ID_Mascota IN NUMBER, p_Fecha_Cita IN DATE) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM citas_tablas.citas
    WHERE id_mascota = p_ID_Mascota 
      AND TRUNC(fecha_cita) = TRUNC(p_Fecha_Cita) -- Validar misma fecha
      AND estado = 'Activa';

    RETURN v_count > 0;
END;
/

-- Función para verificar si un veterinario existe
CREATE OR REPLACE FUNCTION existeVeterinario(p_ID_Veterinario IN NUMBER) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM usuarios_tablas.veterinarios
    WHERE id_veterinario = p_ID_Veterinario;

    RETURN v_count > 0;
END;
/

-- Función para verificar si un veterinario está disponible
CREATE OR REPLACE FUNCTION veterinarioDisponible(p_ID_Veterinario IN NUMBER, p_Fecha_Cita IN DATE) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM citas_tablas.citas
    WHERE id_veterinario = p_ID_Veterinario
      AND fecha_cita = p_Fecha_Cita
      AND estado = 'Activa';

    RETURN v_count = 0; -- Devuelve TRUE si no hay citas activas en esa fecha
END;
/

-- Función para verificar si un servicio existe
CREATE OR REPLACE FUNCTION existeServicio(p_ID_Servicio IN NUMBER) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM servicios_tablas.servicios
    WHERE id_servicio = p_ID_Servicio;

    RETURN v_count > 0;
END;
/

-- Función para verificar si un servicio requiere productos
CREATE OR REPLACE FUNCTION servicioRequiereProductos(p_ID_Servicio IN NUMBER) RETURN SYS_REFCURSOR IS
    v_cursor SYS_REFCURSOR;
BEGIN
    -- Abrir un cursor para devolver los productos asociados al servicio y sus cantidades
    OPEN v_cursor FOR
        SELECT 
            sp.id_producto, 
            sp.unidades_producto
        FROM 
            servicios_tablas.servicios_productos sp
        WHERE 
            sp.id_servicio = p_ID_Servicio;

    RETURN v_cursor; -- Retornar el cursor con los resultados
END;
/

-- Función para verificar si un producto existe
CREATE OR REPLACE FUNCTION existeProducto(p_ID_Producto IN NUMBER) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM servicios_tablas.productos
    WHERE id_producto = p_ID_Producto;

    RETURN v_count > 0;
END;
/

-- Función para insertar una cita
CREATE OR REPLACE FUNCTION insertarCita(
    p_ID_Mascota IN NUMBER,
    p_ID_Veterinario IN NUMBER,
    p_Fecha_Cita IN DATE,
    p_Estado IN VARCHAR2
) RETURN NUMBER IS
    v_ID_Cita NUMBER;
BEGIN
    INSERT INTO citas_tablas.citas (id_mascota, id_veterinario, fecha_cita, estado)
    VALUES (p_ID_Mascota, p_ID_Veterinario, p_Fecha_Cita, p_Estado)
    RETURNING id_cita INTO v_ID_Cita;

    RETURN v_ID_Cita;
END;
/

-- Función para insertar un servicio asociado a una cita
CREATE OR REPLACE FUNCTION insertarCitaServicio(
    p_ID_Cita IN NUMBER,
    p_ID_Servicio IN NUMBER
) RETURN NUMBER IS
    v_ID_Cita_Servicio NUMBER;
BEGIN
    INSERT INTO citas_tablas.citas_servicios (id_cita, id_servicio)
    VALUES (p_ID_Cita, p_ID_Servicio)
    RETURNING id_cita_servicio INTO v_ID_Cita_Servicio;

    RETURN v_ID_Cita_Servicio;
END;
/

CREATE OR REPLACE PROCEDURE actualizarStock(
    p_ID_Producto IN NUMBER,
    p_Cantidad IN NUMBER
) IS
BEGIN
    UPDATE servicios_tablas.productos
    SET stock = stock - p_Cantidad
    WHERE id_producto = p_ID_Producto;

    IF SQL%ROWCOUNT = 0 THEN
        RAISE_APPLICATION_ERROR(-20009, 'No se pudo actualizar el stock para el producto con ID ' || p_ID_Producto);
    END IF;
END;
/

CREATE OR REPLACE FUNCTION calcularTotalServicios(
    p_ID_Cita IN NUMBER
) RETURN NUMBER IS
    v_Total NUMBER;
BEGIN
    SELECT SUM(s.precio)
    INTO v_Total
    FROM citas_tablas.citas_servicios cs
    JOIN servicios_tablas.servicios s ON cs.id_servicio = s.id_servicio
    WHERE cs.id_cita = p_ID_Cita;

    RETURN v_Total;
END;
/

CREATE OR REPLACE FUNCTION crearFactura(
    p_Total IN NUMBER
) RETURN NUMBER IS
    v_ID_Factura NUMBER;
BEGIN
    INSERT INTO citas_tablas.facturas (total, fecha_factura)
    VALUES (p_Total, SYSDATE)
    RETURNING id_factura INTO v_ID_Factura;

    RETURN v_ID_Factura;
END;
/

CREATE OR REPLACE PROCEDURE asociarFacturaConCita(
    p_ID_Cita_Servicio IN NUMBER,
    p_ID_Factura IN NUMBER
) IS
BEGIN
    UPDATE citas_tablas.citas_servicios
    SET facturas_id_factura = p_ID_Factura
    WHERE id_cita_servicio = p_ID_Cita_Servicio;

    IF SQL%ROWCOUNT = 0 THEN
        RAISE_APPLICATION_ERROR(-20010, 'No se pudo asociar la factura con la cita-servicio');
    END IF;
END;
/

-- Procedimiento para agendar una cita cambio hoy 24/04/25
CREATE OR REPLACE PROCEDURE agendarCita(
    p_ID_Cliente IN NUMBER,
    p_ID_Mascota IN NUMBER,
    p_ID_Veterinario IN NUMBER,
    p_Fecha_Cita IN DATE,
    p_Servicios IN VARCHAR2 -- Cambiado a cadena separada por comas
) AS
    v_ID_Cita NUMBER;
    v_ID_Factura NUMBER;
    v_idCitaServicio NUMBER;
    v_Total NUMBER;
    v_Servicios SYS.ODCINUMBERLIST;
BEGIN
    -- Convertir la cadena separada por comas en una colección
    SELECT CAST(MULTISET(
        SELECT TO_NUMBER(REGEXP_SUBSTR(p_Servicios, '[^,]+', 1, LEVEL))
        FROM DUAL
        CONNECT BY REGEXP_SUBSTR(p_Servicios, '[^,]+', 1, LEVEL) IS NOT NULL
    ) AS SYS.ODCINUMBERLIST) INTO v_Servicios FROM DUAL;

    -- Validar cliente
    IF NOT existeCliente(p_ID_Cliente) THEN
        RAISE_APPLICATION_ERROR(-20001, 'El cliente no existe');
    END IF;

    -- Validar mascota
    IF NOT existeMascota(p_ID_Mascota) THEN
        RAISE_APPLICATION_ERROR(-20002, 'La mascota no existe');
    END IF;

    IF NOT mascotaPerteneceACliente(p_ID_Mascota, p_ID_Cliente) THEN
        RAISE_APPLICATION_ERROR(-20003, 'La mascota no pertenece al cliente');
    END IF;

    -- Validar si la mascota ya tiene una cita activa
    IF mascotaTieneCitaActivaNowMismoServicio(p_ID_Mascota, p_Fecha_Cita) THEN
        RAISE_APPLICATION_ERROR(-20004, 'La mascota ya tiene una cita activa a la misma hora');
    END IF;

    -- Validar veterinario
    IF NOT existeVeterinario(p_ID_Veterinario) THEN
        RAISE_APPLICATION_ERROR(-20005, 'El veterinario no existe');
    END IF;

    -- Validar disponibilidad del veterinario
    IF NOT veterinarioDisponible(p_ID_Veterinario, p_Fecha_Cita) THEN
        RAISE_APPLICATION_ERROR(-20006, 'El veterinario no está disponible en esa fecha y hora');
    END IF;

    -- Crear la cita
    v_ID_Cita := insertarCita(p_ID_Mascota, p_ID_Veterinario, p_Fecha_Cita, 'Activa');

    -- Registrar servicios
    FOR i IN 1..v_Servicios.COUNT LOOP
        IF NOT existeServicio(v_Servicios(i)) THEN
            RAISE_APPLICATION_ERROR(-20007, 'El servicio solicitado no es válido');
        END IF;

        v_idCitaServicio := insertarCitaServicio(v_ID_Cita, v_Servicios(i));

        -- Validar y actualizar stock de productos asociados al servicio
        FOR producto IN (
            SELECT sp.id_producto, sp.unidades_producto, p.stock
            FROM servicios_tablas.servicios_productos sp
            JOIN servicios_tablas.productos p ON sp.id_producto = p.id_producto
            WHERE sp.id_servicio = v_Servicios(i)
        ) LOOP
            IF producto.stock < producto.unidades_producto THEN
                RAISE_APPLICATION_ERROR(-20008, 'No hay suficiente stock para el producto con ID ' || producto.id_producto);
            END IF;

            -- Actualizar el stock del producto
            actualizarStock(producto.id_producto, producto.unidades_producto);
        END LOOP;
    END LOOP;

    -- Generar factura
    v_Total := calcularTotalServicios(v_ID_Cita);
    v_ID_Factura := crearFactura(v_Total);

    -- Asociar factura con los servicios de la cita
    FOR i IN 1..v_Servicios.COUNT LOOP
        UPDATE citas_tablas.citas_servicios
        SET facturas_id_factura = v_ID_Factura
        WHERE id_cita = v_ID_Cita AND id_servicio = v_Servicios(i);
    END LOOP;

    COMMIT;
    DBMS_OUTPUT.PUT_LINE('Cita agendada exitosamente. ID_Cita: ' || v_ID_Cita || ', ID_Factura: ' || v_ID_Factura);
EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END agendarCita;
/

-- Funciones necesarias para el registro de clientes

CREATE OR REPLACE FUNCTION existeCedula(p_didentidad_cliente IN VARCHAR2) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM usuarios_tablas.clientes
    WHERE didentidad_cliente = p_didentidad_cliente;

    RETURN v_count > 0;
END;
/

CREATE OR REPLACE FUNCTION existeCorreo(p_email IN VARCHAR2) RETURN BOOLEAN IS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM usuarios_tablas.USUARIOS
    WHERE correo = p_email;

    RETURN v_count > 0;
END;
/

CREATE OR REPLACE FUNCTION insertarCliente(
    p_didentidad_cliente IN VARCHAR2,
    p_nombre IN VARCHAR2,
    p_apellido IN VARCHAR2,
    p_email IN VARCHAR2,
    p_telefono IN VARCHAR2,
    p_direccion IN VARCHAR2,
    p_fecha_registro IN DATE
) RETURN NUMBER IS
    v_ID_Cliente NUMBER;
BEGIN
    INSERT INTO usuarios_tablas.clientes (
        didentidad_cliente, nombre, apellido, email, telefono, direccion, fecha_registro
    ) VALUES (
        p_didentidad_cliente, p_nombre, p_apellido, p_email, p_telefono, p_direccion, p_fecha_registro
    )
    RETURNING id_cliente INTO v_ID_Cliente;

    RETURN v_ID_Cliente;
END;
/

CREATE OR REPLACE FUNCTION insertarUsuario(
    p_id_cliente IN NUMBER, -- Nuevo parámetro para el ID del cliente
    p_email IN VARCHAR2,
    p_contrasena IN VARCHAR2
) RETURN NUMBER IS
    v_ID_Usuario NUMBER;
BEGIN
    INSERT INTO usuarios_tablas.usuarios (
        id_cliente, correo, contrasena
    ) VALUES (
        p_id_cliente, p_email, p_contrasena
    )
    RETURNING id_usuario INTO v_ID_Usuario;

    RETURN v_ID_Usuario;
END;
/

-- Procedimiento para registrar un cliente
CREATE OR REPLACE PROCEDURE registrarCliente(
    p_didentidad_cliente IN VARCHAR2,
    p_nombre IN VARCHAR2,
    p_apellido IN VARCHAR2,
    p_email IN VARCHAR2,
    p_telefono IN VARCHAR2,
    p_direccion IN VARCHAR2,
    p_contrasena IN VARCHAR2,
    p_ID_Cliente OUT NUMBER -- Nuevo parámetro de salida para devolver el ID del cliente
) AS
    v_ID_Usuario NUMBER;
BEGIN
    -- Validar cédula o documento de identidad
    IF existeCedula(p_didentidad_cliente) THEN
        RAISE_APPLICATION_ERROR(-20010, 'La cédula o documento de identidad ya está registrada');
    END IF;

    -- Validar correo electrónico en la tabla de clientes
    IF existeCorreo(p_email) THEN
        RAISE_APPLICATION_ERROR(-20011, 'El correo electrónico ya está registrado');
    END IF;

    -- Insertar cliente en la tabla de clientes
    p_ID_Cliente := insertarCliente(
        p_didentidad_cliente,
        p_nombre,
        p_apellido,
        p_email,
        p_telefono,
        p_direccion,
        SYSDATE
    );

    -- Insertar usuario en la tabla de usuarios
    v_ID_Usuario := insertarUsuario(
        p_ID_Cliente, -- Pasar el ID del cliente
        p_email,
        p_contrasena
    );

    COMMIT;
    DBMS_OUTPUT.PUT_LINE('Cliente registrado exitosamente. ID_Cliente: ' || p_ID_Cliente || ', ID_Usuario: ' || v_ID_Usuario);
EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END registrarCliente;
/

-- Datos de prueba
DECLARE
    v_ID_Cliente NUMBER; -- Variable para capturar el ID del cliente registrado
BEGIN
    -- Registrar un cliente
    registrarCliente(
        p_didentidad_cliente => '123456789',
        p_nombre => 'Fran',
        p_apellido => 'Rojas',
        p_email => 'andrjs@outlook.com',
        p_telefono => '85075310',
        p_direccion => 'San José, Costa Rica',
        p_contrasena => 'contrasena123',
        p_ID_Cliente => v_ID_Cliente -- Capturar el ID del cliente registrado
    );

    -- Mostrar el ID del cliente registrado
    DBMS_OUTPUT.PUT_LINE('Cliente registrado con éxito. ID del cliente: ' || v_ID_Cliente);
END;
/
--24/04/2025
-- Datos de prueba
DECLARE
    v_ID_Cliente NUMBER := 41; -- Reemplaza con el ID del cliente
    v_ID_Mascota NUMBER := 1; -- Reemplaza con el ID de la mascota
    v_ID_Veterinario NUMBER := 1; -- Reemplaza con el ID del veterinario
    v_Fecha_Cita DATE := TO_DATE('2025-04-30', 'YYYY-MM-DD'); -- Reemplaza con la fecha deseada
    v_Servicios VARCHAR2(100) := '3,4'; -- Reemplaza con los IDs de los servicios separados por comas
BEGIN
    -- Llamar al procedimiento para agendar la cita
    agendarCita(
        p_ID_Cliente => v_ID_Cliente,
        p_ID_Mascota => v_ID_Mascota,
        p_ID_Veterinario => v_ID_Veterinario,
        p_Fecha_Cita => v_Fecha_Cita,
        p_Servicios => v_Servicios
    );

    DBMS_OUTPUT.PUT_LINE('Cita agendada exitosamente.');
END;
/
commit;

SELECT object_name, status
FROM user_objects
WHERE object_name = 'AGENDARCITA';

SELECT c.ID_CITA, c.fecha_cita, c.estado, m.nombre AS mascota, v.nombre AS veterinario
FROM citas_tablas.citas c
JOIN usuarios_tablas.mascotas m ON c.ID_MASCOTA = m.ID_MASCOTA
JOIN usuarios_tablas.veterinarios v ON c.id_veterinario = v.id_veterinario
WHERE m.id_cliente = 41;

SELECT f.ID_FACTURA, f.TOTAL, f.FECHA_FACTURA, 
       CASE 
           WHEN c.FECHA_CITA < SYSDATE THEN 'Pagada'
           ELSE 'Pendiente'
       END AS ESTADO
FROM CITAS_TABLAS.FACTURAS f
WHERE m.ID_CLIENTE = 41;