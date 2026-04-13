-- Migración: agregar columnas de estado y Stripe a FACTURAS
-- Ejecutar después de los scripts 1-5
USE patitas_rescate;

ALTER TABLE FACTURAS
    ADD COLUMN ESTADO VARCHAR(20) NOT NULL DEFAULT 'Pendiente',
    ADD COLUMN STRIPE_SESSION_ID VARCHAR(255) DEFAULT NULL;

-- Marcar como pagadas las facturas cuya cita ya pasó (migración de datos existentes)
UPDATE FACTURAS f
    JOIN CITAS_SERVICIOS cs ON cs.FACTURAS_ID_FACTURA = f.ID_FACTURA
    JOIN CITAS c ON cs.ID_CITA = c.ID_CITA
SET f.ESTADO = 'Pagada'
WHERE c.FECHA_CITA < NOW();
