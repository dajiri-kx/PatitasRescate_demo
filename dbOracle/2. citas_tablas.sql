-- Tabla de Citas
CREATE TABLE citas (
    id_cita INT CONSTRAINT pk_citas PRIMARY KEY,
    id_mascota INT NOT NULL CONSTRAINT fk_citas_mascotas REFERENCES usuarios_tablas.mascotas(id_mascota),
    id_veterinario INT NOT NULL CONSTRAINT fk_citas_veterinarios REFERENCES usuarios_tablas.veterinarios(id_veterinario),
    fecha_cita DATE NOT NULL,
    estado VARCHAR2(20) NOT NULL CONSTRAINT chk_citas_estado
        CHECK (estado IN ('Activa', 'Completada', 'Cancelada'))
);

CREATE INDEX idx_citas_mascota ON citas(id_mascota);
CREATE INDEX idx_citas_veterinario ON citas(id_veterinario);
CREATE INDEX idx_citas_estado_fecha ON citas(estado, fecha_cita);

-- Tabla de Facturas
CREATE TABLE facturas (
    id_factura INT CONSTRAINT pk_facturas PRIMARY KEY,
    fecha_factura DATE DEFAULT SYSDATE NOT NULL,
    total DECIMAL(10,2) NOT NULL CONSTRAINT chk_facturas_total CHECK (total >= 0)
);

-- Tabla de Relación entre Citas y Servicios
CREATE TABLE citas_servicios (
    id_cita_servicio INT GENERATED ALWAYS AS IDENTITY CONSTRAINT pk_citas_servicios PRIMARY KEY,
    id_cita INT NOT NULL CONSTRAINT fk_citas_servicios_citas REFERENCES citas(id_cita),
    id_servicio INT NOT NULL CONSTRAINT fk_citas_servicios_servicios REFERENCES servicios_tablas.servicios(id_servicio),
    facturas_id_factura INT CONSTRAINT fk_citas_servicios_facturas REFERENCES facturas(id_factura)
);

CREATE INDEX idx_citasserv_cita ON citas_servicios(id_cita);
CREATE INDEX idx_citasserv_servicio ON citas_servicios(id_servicio);
CREATE INDEX idx_citasserv_factura ON citas_servicios(facturas_id_factura);

-- Permisos cross-schema
GRANT REFERENCES ON citas_tablas.citas_servicios TO servicios_tablas;

-- Secuencias (START WITH 100 para evitar colisión con datos manuales)
CREATE SEQUENCE seq_id_cita START WITH 100 INCREMENT BY 1;
CREATE SEQUENCE seq_id_factura START WITH 100 INCREMENT BY 1;

-- Triggers de auto-incremento
CREATE OR REPLACE TRIGGER trg_id_cita
BEFORE INSERT ON citas_tablas.citas
FOR EACH ROW
BEGIN
    IF :NEW.id_cita IS NULL THEN
        :NEW.id_cita := seq_id_cita.NEXTVAL;
    END IF;
END;
/

CREATE OR REPLACE TRIGGER trg_id_factura
BEFORE INSERT ON citas_tablas.facturas
FOR EACH ROW
BEGIN
    IF :NEW.id_factura IS NULL THEN
        :NEW.id_factura := seq_id_factura.NEXTVAL;
    END IF;
END;
/
