-- Tabla de Servicios
CREATE TABLE servicios (
    id_servicio INT CONSTRAINT pk_servicios PRIMARY KEY,
    nombre_servicio VARCHAR2(100) NOT NULL,
    descripcion CLOB,
    precio DECIMAL(10,2) NOT NULL CONSTRAINT chk_servicios_precio CHECK (precio >= 0),
    duracion_minutos INT CONSTRAINT chk_servicios_duracion CHECK (duracion_minutos > 0)
);

-- Tabla de Productos
CREATE TABLE productos (
    id_producto INT CONSTRAINT pk_productos PRIMARY KEY,
    nombre_producto VARCHAR2(100) NOT NULL,
    categoria VARCHAR2(100),
    precio DECIMAL(10,2) NOT NULL CONSTRAINT chk_productos_precio CHECK (precio >= 0),
    stock INT NOT NULL CONSTRAINT chk_productos_stock CHECK (stock >= 0)
);

-- Tabla de Relación entre Servicios y Productos
CREATE TABLE servicios_productos (
    id_producto_servicio INT CONSTRAINT pk_servicios_productos PRIMARY KEY,
    id_producto INT NOT NULL CONSTRAINT fk_servicios_productos REFERENCES productos(id_producto),
    id_servicio INT NOT NULL CONSTRAINT fk_servicios_productos_servicios REFERENCES servicios(id_servicio),
    unidades_producto INT NOT NULL CONSTRAINT chk_sp_unidades CHECK (unidades_producto > 0),
    cantidad_consumida INT NOT NULL CONSTRAINT chk_sp_cantidad CHECK (cantidad_consumida > 0)
);

CREATE INDEX idx_servprod_producto ON servicios_productos(id_producto);
CREATE INDEX idx_servprod_servicio ON servicios_productos(id_servicio);

-- Permisos cross-schema
GRANT REFERENCES ON servicios_tablas.servicios TO citas_tablas;

-- Secuencias (START WITH 100 para evitar colisión con datos semilla)
CREATE SEQUENCE seq_id_servicio START WITH 100 INCREMENT BY 1;
CREATE SEQUENCE seq_id_producto START WITH 100 INCREMENT BY 1;
CREATE SEQUENCE seq_id_producto_servicio START WITH 100 INCREMENT BY 1;

-- Triggers de auto-incremento
CREATE OR REPLACE TRIGGER trg_id_servicio
BEFORE INSERT ON servicios_tablas.servicios
FOR EACH ROW
BEGIN
    IF :NEW.id_servicio IS NULL THEN
        :NEW.id_servicio := seq_id_servicio.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_id_producto
BEFORE INSERT ON servicios_tablas.productos
FOR EACH ROW
BEGIN
    IF :NEW.id_producto IS NULL THEN
        :NEW.id_producto := seq_id_producto.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_id_producto_servicio
BEFORE INSERT ON servicios_tablas.servicios_productos
FOR EACH ROW
BEGIN
    IF :NEW.id_producto_servicio IS NULL THEN
        :NEW.id_producto_servicio := seq_id_producto_servicio.NEXTVAL;
    END IF;
END;
/

-- Datos semilla: servicios
INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (1, 'Corte Mini', 'Corte de pelo para perros pequeños', 15000.00, 30);

INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (2, 'Corte Mediano', 'Corte de pelo para perros medianos', 25000.00, 45);

INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (3, 'Corte Grande', 'Corte de pelo para perros grandes', 35000.00, 60);

INSERT INTO servicios (id_servicio, nombre_servicio, descripcion, precio, duracion_minutos)
VALUES (4, 'Corte de Uñas', 'Corte de uñas para perros de cualquier tamaño', 10000.00, 20);

-- Datos semilla: productos
INSERT INTO productos (id_producto, nombre_producto, categoria, precio, stock)
VALUES (1, 'Shampoo burbuja', 'Higiene', 1500.00, 100);

INSERT INTO productos (id_producto, nombre_producto, categoria, precio, stock)
VALUES (2, 'Acondicionador burbuja', 'Higiene', 1350.00, 100);

INSERT INTO productos (id_producto, nombre_producto, categoria, precio, stock)
VALUES (3, 'Desinfectante', 'Higiene', 400.00, 50);

-- Asociar productos con servicios
-- Corte Mini
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (1, 1, 1, 1, 1);
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (2, 2, 1, 1, 1);
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (3, 3, 1, 1, 1);

-- Corte Mediano
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (4, 1, 2, 2, 2);
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (5, 2, 2, 2, 2);
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (6, 3, 2, 1, 1);

-- Corte Grande
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (7, 1, 3, 3, 3);
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (8, 2, 3, 3, 3);
INSERT INTO servicios_productos (id_producto_servicio, id_producto, id_servicio, unidades_producto, cantidad_consumida)
VALUES (9, 3, 3, 1, 1);

COMMIT;

-- Verificación
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
