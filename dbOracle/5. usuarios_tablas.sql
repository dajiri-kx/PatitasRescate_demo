-- Tabla de Clientes
CREATE TABLE clientes (
    id_cliente INT CONSTRAINT pk_clientes PRIMARY KEY, -- identificacion interna de la base de datos
    didentidad_cliente VARCHAR2(20), -- se agrega este campo ya que se requiere una identificacion valida para el usaurio y no se puede usar el id_cliente como tal
    nombre VARCHAR2(100),
    apellido VARCHAR2(100),
    email VARCHAR2(100) CONSTRAINT uq_clientes_email UNIQUE,
    telefono VARCHAR2(15),
    direccion VARCHAR2(255),
    fecha_registro DATE
);

-- Tabla de Mascotas
CREATE TABLE mascotas (
    id_mascota INT CONSTRAINT pk_mascotas PRIMARY KEY, -- identificacion interna de la base de datos
    id_cliente INT CONSTRAINT fk_mascotas_clientes REFERENCES clientes(ID_Cliente), -- se agrega este campo ya que se requiere una identificacion valida para el usaurio y no se puede usar el id_cliente como tal
    nombre VARCHAR2(100),
    especie VARCHAR2(50),
    raza VARCHAR2(50),
    meses INT,
    historial_medica CLOB
);

-- Tabla de Veterinarios
CREATE TABLE veterinarios (
    id_veterinario NUMBER(10) CONSTRAINT pk_veterinarios PRIMARY KEY, -- identificacion interna de la base de datos
    didentidad_veterinario VARCHAR2(20), -- se agrega este campo ya que se requiere una identificacion valida para el usaurio y no se puede usar el id_veterinario como tal
    nombre VARCHAR2(100),
    especialidad VARCHAR2(100),
    telefono VARCHAR2(15),
    correo VARCHAR2(100) CONSTRAINT uq_veterinarios_correo UNIQUE,
    rol VARCHAR2(20)
);

-- Otorgar permiso REFERENCES para la tabla mascotas
GRANT REFERENCES ON usuarios_tablas.mascotas TO citas_tablas;

-- Otorgar permiso REFERENCES para la tabla veterinarios
GRANT REFERENCES ON usuarios_tablas.veterinarios TO citas_tablas;

-- Tabla de Usuarios
CREATE TABLE usuarios (
    id_usuario INT CONSTRAINT pk_usuarios PRIMARY KEY, -- Identificación interna de la base de datos
    id_cliente INT CONSTRAINT fk_usuarios_clientes REFERENCES clientes(id_cliente), -- Relación con la tabla de clientes
    correo VARCHAR2(100) CONSTRAINT uq_usuarios_correo UNIQUE, -- Correo único para el usuario
    contrasena VARCHAR2(255) -- Contraseña encriptada
);

-- Otorgar permisos al usuario progra para las tablas de usuarios_tablas
GRANT INSERT ON usuarios_tablas.clientes TO progra;
GRANT INSERT ON usuarios_tablas.usuarios TO progra;

GRANT SELECT ON clientes TO Progra_PAR;
GRANT SELECT ON usuarios TO Progra_PAR;
GRANT SELECT ON mascotas TO Progra_PAR;
GRANT SELECT ON veterinarios TO Progra_PAR;

GRANT INSERT, UPDATE ON clientes TO Progra_PAR;
GRANT INSERT, UPDATE ON usuarios TO Progra_PAR;

-- Otorgar permisos REFERENCES para las tablas relacionadas
GRANT REFERENCES ON clientes TO Progra_PAR; -- Permite referencias desde otras tablas a clientes
GRANT REFERENCES ON mascotas TO Progra_PAR; -- Permite referencias desde otras tablas a mascotas
GRANT REFERENCES ON veterinarios TO Progra_PAR; -- Permite referencias desde otras tablas a veterinarios
GRANT REFERENCES ON usuarios TO Progra_PAR; -- Permite referencias desde otras tablas a usuarios

GRANT ALTER ON usuarios_tablas.clientes TO Progra_PAR;
GRANT INSERT ON usuarios_tablas.clientes TO Progra_PAR;

CREATE SEQUENCE seq_id_cliente
START WITH 1
INCREMENT BY 1;

CREATE OR REPLACE TRIGGER trg_id_cliente
BEFORE INSERT ON usuarios_tablas.clientes
FOR EACH ROW
BEGIN
    IF :NEW.id_cliente IS NULL THEN
        :NEW.id_cliente := seq_id_cliente.NEXTVAL;
    END IF;
END;
/

CREATE SEQUENCE seq_id_usuario
START WITH 1
INCREMENT BY 1;

CREATE OR REPLACE TRIGGER trg_id_usuario
BEFORE INSERT ON usuarios_tablas.usuarios
FOR EACH ROW
BEGIN
    IF :NEW.id_usuario IS NULL THEN
        :NEW.id_usuario := seq_id_cliente.NEXTVAL;
    END IF;
END;
/

CREATE SEQUENCE seq_id_mascota
START WITH 1
INCREMENT BY 1;

CREATE OR REPLACE TRIGGER trg_id_mascota
BEFORE INSERT ON usuarios_tablas.mascotas
FOR EACH ROW
BEGIN
    IF :NEW.id_mascota IS NULL THEN
        :NEW.id_mascota := seq_id_mascota.NEXTVAL;
    END IF;
END;
/

-- 21/04/2025

-- Inserción de veterinarios estéticos
INSERT INTO veterinarios (id_veterinario, didentidad_veterinario, nombre, especialidad, telefono, correo, rol)
VALUES (1, 'VET123456789', 'Dr. Juan Pérez', 'Estética Canina', '88881234', 'juan.perez@veterinaria.com', 'Estético');

INSERT INTO veterinarios (id_veterinario, didentidad_veterinario, nombre, especialidad, telefono, correo, rol)
VALUES (2, 'VET987654321', 'Dra. María López', 'Estética Felina', '88885678', 'maria.lopez@veterinaria.com', 'Estético');

INSERT INTO veterinarios (id_veterinario, didentidad_veterinario, nombre, especialidad, telefono, correo, rol)
VALUES (3, 'VET456789123', 'Dr. Carlos Ramírez', 'Estética General', '88889012', 'carlos.ramirez@veterinaria.com', 'Estético');

INSERT INTO MASCOTAS (
    ID_MASCOTA, ID_CLIENTE, NOMBRE, ESPECIE, RAZA, MESES, HISTORIAL_MEDICA
) VALUES (
    1, 41, 'Max', 'Perro', 'Golden Retriever', 24, EMPTY_CLOB()
);

commit;

CREATE SEQUENCE seq_id_veterinario
START WITH 1
INCREMENT BY 1;

CREATE OR REPLACE TRIGGER trg_id_veterinario
BEFORE INSERT ON usuarios_tablas.veterinarios
FOR EACH ROW
BEGIN
    IF :NEW.id_veterinario IS NULL THEN
        :NEW.id_veterinario := seq_id_veterinario.NEXTVAL;
    END IF;
END;
/

