-- Tabla de Servicios
CREATE TABLE servicios (
    id_servicio INT CONSTRAINT pk_servicios PRIMARY KEY, -- identificacion interna de la base de datos
    nombre_servicio VARCHAR2(100), -- nombre del servicio ofrecido
    descripcion CLOB, -- descripcion detallada del servicio
    precio DECIMAL(10,2), -- costo del servicio
    duracion_minutos INT -- duracion estimada del servicio en minutos
);

-- Tabla de Productos
CREATE TABLE productos (
    id_producto INT CONSTRAINT pk_productos PRIMARY KEY, -- identificacion interna de la base de datos
    nombre_producto VARCHAR2(100), -- nombre del producto
    categoria VARCHAR2(100), -- categoria del producto
    precio DECIMAL(10,2), -- precio del producto
    stock INT -- cantidad disponible en inventario
);

-- Tabla de Relación entre Servicios y Productos

CREATE TABLE servicios_productos (
    id_producto_servicio INT CONSTRAINT pk_servicios_productos PRIMARY KEY, -- Identificador único para la tabla
    id_producto INT CONSTRAINT fk_servicios_productos REFERENCES productos(id_producto), -- Referencia a la tabla de productos
    id_servicio INT CONSTRAINT fk_servicios_productos_servicios REFERENCES servicios(id_servicio), -- Referencia a la tabla de servicios
    unidades_producto INT, -- Cantidad de unidades del producto consumido
    cantidad_consumida INT -- Cantidad total del producto consumido
);


-- Otorgar permiso REFERENCES para la tabla servicios
GRANT REFERENCES ON servicios_tablas.servicios TO citas_tablas;

-- Otorgar permisos al usuario progra para las tablas de servicios_tablas
GRANT UPDATE ON servicios_tablas.productos TO progra;
GRANT REFERENCES ON servicios_tablas.servicios TO progra;

--
GRANT SELECT ON SERVICIOS_TABLAS.servicios TO Progra_PAR;
GRANT SELECT ON SERVICIOS_TABLAS.productos TO Progra_PAR;
GRANT SELECT ON SERVICIOS_TABLAS.servicios_productos TO Progra_PAR;
GRANT INSERT, UPDATE ON SERVICIOS_TABLAS.servicios_productos TO Progra_PAR;

-- Otorgar permisos REFERENCES para las tablas relacionadas
GRANT REFERENCES ON SERVICIOS_TABLAS.servicios TO Progra_PAR; -- Permite referencias desde otras tablas a servicios
GRANT REFERENCES ON SERVICIOS_TABLAS.productos TO Progra_PAR; -- Permite referencias desde otras tablas a productos
GRANT REFERENCES ON SERVICIOS_TABLAS.servicios_productos TO Progra_PAR; -- Permite referencias desde otras tablas a servicios_productos
GRANT UPDATE ON SERVICIOS_TABLAS.PRODUCTOS TO Progra_PAR;
GRANT SELECT ON servicios_tablas.servicios TO Progra_PAR;


-- Insertar el servicio "Corte de Pelo" con tres variantes
INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (1, 'Corte Mini', 'Corte de pelo para perros pequeños', 15000.00, 30);

INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (2, 'Corte Mediano', 'Corte de pelo para perros medianos', 25000.00, 45);

INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (3, 'Corte Grande', 'Corte de pelo para perros grandes', 35000.00, 60);

-- Insertar el servicio "Corte de Uñas"
INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (4, 'Corte de Uñas', 'Corte de uñas para perros de cualquier tamaño', 10000.00, 20);

-- Insertar productos necesarios para el corte de pelo
INSERT INTO productos (id_producto, nombre_producto, categoria, precio, stock)
VALUES (1, 'Shampoo burbuja', 'Higiene', 1500.00, 100);

INSERT INTO productos (id_producto, nombre_producto, categoria, precio, stock)
VALUES (2, 'Acondicionador burbuja', 'Higiene', 1350.00, 100);

INSERT INTO productos (id_producto, nombre_producto, categoria, precio, stock)
VALUES (3, 'Desinfectante', 'Higiene', 400.00, 50);

-- Asociar productos con el servicio "Corte Mini"
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (1, 1, 1, 1, 1); -- Shampoo
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (2, 2, 1, 1, 1); -- Acondicionador
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (3, 3, 1, 1, 1); -- Desinfectante

-- Asociar productos con el servicio "Corte Mediano"
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (4, 1, 2, 2, 2); -- Shampoo
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (5, 2, 2, 2, 2); -- Acondicionador
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (6, 3, 2, 1, 1); -- Desinfectante

-- Asociar productos con el servicio "Corte Grande"
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (7, 1, 3, 3, 3); -- Shampoo
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (8, 2, 3, 3, 3); -- Acondicionador
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (9, 3, 3, 1, 1); -- Desinfectante

SELECT 
    s.nombre_servicio AS Servicio,
    p.nombre_producto AS Producto,
    sp.unidades_producto AS Unidades_Requeridas,
    sp.cantidad_consumida AS Cantidad_Consumida
FROM 
    servicios_tablas.servicios s
LEFT JOIN 
    servicios_tablas.servicios_productos sp ON s.id_servicio = sp.id_servicio
LEFT JOIN 
    servicios_tablas.productos p ON sp.id_producto = p.id_producto
ORDER BY 
    s.nombre_servicio, p.nombre_producto;

GRANT SELECT ON servicios_tablas.servicios TO Progra_PAR;
GRANT SELECT ON servicios_tablas.servicios_productos TO Progra_PAR;
GRANT SELECT ON servicios_tablas.productos TO Progra_PAR;

-- /24/04/25

CREATE SEQUENCE seq_id_servicio
START WITH 1
INCREMENT BY 1;

CREATE OR REPLACE TRIGGER trg_id_servicio
BEFORE INSERT ON servicios_tablas.servicios
FOR EACH ROW
BEGIN
    IF :NEW.id_servicio IS NULL THEN
        :NEW.id_servicio := seq_id_servicio.NEXTVAL;
    END IF;
END;
/

CREATE SEQUENCE seq_id_producto
START WITH 1
INCREMENT BY 1;

CREATE OR REPLACE TRIGGER trg_id_producto
BEFORE INSERT ON servicios_tablas.productos
FOR EACH ROW
BEGIN
    IF :NEW.id_producto IS NULL THEN
        :NEW.id_producto := seq_id_producto.NEXTVAL;
    END IF;
END;
/

CREATE SEQUENCE seq_id_producto_servicio
START WITH 1
INCREMENT BY 1;

CREATE OR REPLACE TRIGGER trg_id_producto_servicio
BEFORE INSERT ON servicios_tablas.servicios_productos
FOR EACH ROW
BEGIN
    IF :NEW.id_producto_servicio IS NULL THEN
        :NEW.id_producto_servicio := seq_id_producto_servicio.NEXTVAL;
    END IF;
END;
/