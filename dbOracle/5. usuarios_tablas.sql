-- Tabla de Clientes
CREATE TABLE clientes (
    id_cliente INT CONSTRAINT pk_clientes PRIMARY KEY,
    didentidad_cliente VARCHAR2(20) NOT NULL,
    nombre VARCHAR2(100) NOT NULL,
    apellido VARCHAR2(100) NOT NULL,
    email VARCHAR2(100) NOT NULL CONSTRAINT uq_clientes_email UNIQUE,
    telefono VARCHAR2(15) NOT NULL,
    direccion VARCHAR2(255),
    fecha_registro DATE DEFAULT SYSDATE NOT NULL
);

CREATE INDEX idx_clientes_didentidad ON clientes(didentidad_cliente);

-- Tabla de Mascotas
CREATE TABLE mascotas (
    id_mascota INT CONSTRAINT pk_mascotas PRIMARY KEY,
    id_cliente INT NOT NULL CONSTRAINT fk_mascotas_clientes REFERENCES clientes(id_cliente),
    nombre VARCHAR2(100) NOT NULL,
    especie VARCHAR2(50) NOT NULL,
    raza VARCHAR2(50),
    meses INT CONSTRAINT chk_mascotas_meses CHECK (meses >= 0)
);

CREATE INDEX idx_mascotas_cliente ON mascotas(id_cliente);

-- Tabla de Veterinarios
CREATE TABLE veterinarios (
    id_veterinario NUMBER(10) CONSTRAINT pk_veterinarios PRIMARY KEY,
    didentidad_veterinario VARCHAR2(20) NOT NULL,
    nombre VARCHAR2(100) NOT NULL,
    especialidad VARCHAR2(100),
    telefono VARCHAR2(15),
    correo VARCHAR2(100) CONSTRAINT uq_veterinarios_correo UNIQUE,
    rol VARCHAR2(20)
);

-- Permisos REFERENCES para citas_tablas
GRANT REFERENCES ON usuarios_tablas.mascotas TO citas_tablas;
GRANT REFERENCES ON usuarios_tablas.veterinarios TO citas_tablas;

-- Tabla de Usuarios (autenticación)
CREATE TABLE usuarios (
    id_usuario INT CONSTRAINT pk_usuarios PRIMARY KEY,
    id_cliente INT NOT NULL CONSTRAINT fk_usuarios_clientes REFERENCES clientes(id_cliente),
    correo VARCHAR2(100) NOT NULL CONSTRAINT uq_usuarios_correo UNIQUE,
    contrasena VARCHAR2(255) NOT NULL -- Almacena hash bcrypt, nunca texto plano
);

CREATE INDEX idx_usuarios_cliente ON usuarios(id_cliente);

-- Secuencias (START WITH 100 para evitar colisión con datos semilla)
CREATE SEQUENCE seq_id_cliente START WITH 100 INCREMENT BY 1;
CREATE SEQUENCE seq_id_usuario START WITH 100 INCREMENT BY 1;
CREATE SEQUENCE seq_id_mascota START WITH 100 INCREMENT BY 1;
CREATE SEQUENCE seq_id_veterinario START WITH 100 INCREMENT BY 1;

-- Triggers de auto-incremento
CREATE OR REPLACE TRIGGER trg_id_cliente
BEFORE INSERT ON usuarios_tablas.clientes
FOR EACH ROW
BEGIN
    IF :NEW.id_cliente IS NULL THEN
        :NEW.id_cliente := seq_id_cliente.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_id_usuario
BEFORE INSERT ON usuarios_tablas.usuarios
FOR EACH ROW
BEGIN
    IF :NEW.id_usuario IS NULL THEN
        :NEW.id_usuario := seq_id_usuario.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_id_mascota
BEFORE INSERT ON usuarios_tablas.mascotas
FOR EACH ROW
BEGIN
    IF :NEW.id_mascota IS NULL THEN
        :NEW.id_mascota := seq_id_mascota.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_id_veterinario
BEFORE INSERT ON usuarios_tablas.veterinarios
FOR EACH ROW
BEGIN
    IF :NEW.id_veterinario IS NULL THEN
        :NEW.id_veterinario := seq_id_veterinario.NEXTVAL;
    END IF;
END;
/

-- Datos semilla: veterinarios
INSERT INTO veterinarios (id_veterinario, didentidad_veterinario, nombre, especialidad, telefono, correo, rol)
VALUES (1, 'VET123456789', 'Dr. Juan Pérez', 'Estética Canina', '88881234', 'juan.perez@veterinaria.com', 'Estético');

INSERT INTO veterinarios (id_veterinario, didentidad_veterinario, nombre, especialidad, telefono, correo, rol)
VALUES (2, 'VET987654321', 'Dra. María López', 'Estética Felina', '88885678', 'maria.lopez@veterinaria.com', 'Estético');

INSERT INTO veterinarios (id_veterinario, didentidad_veterinario, nombre, especialidad, telefono, correo, rol)
VALUES (3, 'VET456789123', 'Dr. Carlos Ramírez', 'Estética General', '88889012', 'carlos.ramirez@veterinaria.com', 'Estético');

COMMIT;
