-- Migración: Categorías de servicios
-- Ejecutar después de 6. checkout_migracion.sql
USE patitas_rescate;

-- 1. Agregar columna CATEGORIA a SERVICIOS
ALTER TABLE SERVICIOS
    ADD COLUMN CATEGORIA VARCHAR(50) NOT NULL DEFAULT 'General';

-- 2. Clasificar servicios existentes (estéticos)
UPDATE SERVICIOS SET CATEGORIA = 'Estética' WHERE ID_SERVICIO IN (1, 2, 3, 4);

-- 3. Nuevos servicios por categoría
-- Consulta Veterinaria
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Consulta General', 'Evaluación médica completa del estado de salud', 20000.00, 30, 'Consulta');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Consulta Especializada', 'Evaluación con veterinario especialista', 35000.00, 45, 'Consulta');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Control Post-Operatorio', 'Revisión de seguimiento post cirugía', 15000.00, 20, 'Consulta');

-- Vacunación
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Vacuna Antirrábica', 'Vacunación obligatoria contra la rabia', 12000.00, 15, 'Vacunación');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Vacuna Múltiple Canina', 'Parvovirus, moquillo, hepatitis, leptospirosis', 18000.00, 15, 'Vacunación');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Vacuna Triple Felina', 'Panleucopenia, rinotraqueitis, calicivirus', 16000.00, 15, 'Vacunación');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Desparasitación', 'Tratamiento antiparasitario interno y externo', 8000.00, 15, 'Vacunación');

-- Cirugía
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Esterilización Canina', 'Cirugía de esterilización para perros', 65000.00, 90, 'Cirugía');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Esterilización Felina', 'Cirugía de esterilización para gatos', 45000.00, 60, 'Cirugía');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Limpieza Dental', 'Profilaxis dental bajo anestesia', 40000.00, 60, 'Cirugía');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Cirugía de Tejidos Blandos', 'Procedimientos quirúrgicos generales', 80000.00, 120, 'Cirugía');

-- Diagnóstico por Imágenes
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Radiografía', 'Estudio radiográfico digital', 25000.00, 20, 'Diagnóstico');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Ecografía', 'Ultrasonido abdominal o cardíaco', 35000.00, 30, 'Diagnóstico');
INSERT INTO SERVICIOS (NOMBRE_SERVICIO, DESCRIPCION, PRECIO, DURACION_MINUTOS, CATEGORIA)
VALUES ('Exámenes de Laboratorio', 'Hemograma, química sanguínea, urianálisis', 22000.00, 15, 'Diagnóstico');
