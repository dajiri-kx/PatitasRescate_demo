-- Desde System
-- NOTA: Cambiar las contraseñas antes de desplegar en producción.
-- Orden de ejecución recomendado: 1 → 5 → 3 → 2 → 4

-- Creación del usuario Usuarios_Tablas
CREATE USER Usuarios_Tablas IDENTIFIED BY UsuTablas_2026
DEFAULT TABLESPACE users
TEMPORARY TABLESPACE temp
QUOTA 50M ON users;

GRANT CREATE SESSION, CREATE TABLE TO Usuarios_Tablas;
GRANT CREATE SEQUENCE TO Usuarios_Tablas;
GRANT CREATE TRIGGER TO Usuarios_Tablas;

-- Creación del usuario Citas_Tablas
CREATE USER Citas_Tablas IDENTIFIED BY CitTablas_2026
DEFAULT TABLESPACE users
TEMPORARY TABLESPACE temp
QUOTA 50M ON users;

GRANT CREATE SESSION, CREATE TABLE TO Citas_Tablas;
GRANT CREATE SEQUENCE TO Citas_Tablas;
GRANT CREATE TRIGGER TO Citas_Tablas;

-- Creación del usuario Servicios_Tablas
CREATE USER Servicios_Tablas IDENTIFIED BY SrvTablas_2026
DEFAULT TABLESPACE users
TEMPORARY TABLESPACE temp
QUOTA 50M ON users;

GRANT CREATE SESSION, CREATE TABLE TO Servicios_Tablas;
GRANT CREATE SEQUENCE TO Servicios_Tablas;
GRANT CREATE TRIGGER TO Servicios_Tablas;

-- Creación del usuario de aplicación Progra_PAR
CREATE USER Progra_PAR IDENTIFIED BY PrograPAR_2026
DEFAULT TABLESPACE users
TEMPORARY TABLESPACE temp
QUOTA 10M ON users;

-- Privilegios específicos (sin roles genéricos RESOURCE/CONNECT)
GRANT CREATE SESSION TO Progra_PAR;
GRANT CREATE PROCEDURE, CREATE TRIGGER TO Progra_PAR;

-- =============================================
-- Los siguientes GRANTs deben ejecutarse
-- DESPUÉS de crear las tablas (scripts 5, 3, 2)
-- =============================================

-- Permisos sobre tablas de usuarios_tablas
GRANT SELECT ON usuarios_tablas.clientes TO Progra_PAR;
GRANT SELECT ON usuarios_tablas.mascotas TO Progra_PAR;
GRANT SELECT ON usuarios_tablas.veterinarios TO Progra_PAR;
GRANT SELECT ON usuarios_tablas.usuarios TO Progra_PAR;
GRANT INSERT, UPDATE ON usuarios_tablas.clientes TO Progra_PAR;
GRANT INSERT, UPDATE ON usuarios_tablas.usuarios TO Progra_PAR;
GRANT INSERT, UPDATE, DELETE ON usuarios_tablas.mascotas TO Progra_PAR;
GRANT REFERENCES ON usuarios_tablas.clientes TO Progra_PAR;
GRANT REFERENCES ON usuarios_tablas.mascotas TO Progra_PAR;
GRANT REFERENCES ON usuarios_tablas.veterinarios TO Progra_PAR;
GRANT REFERENCES ON usuarios_tablas.usuarios TO Progra_PAR;

-- Permisos sobre tablas de servicios_tablas
GRANT SELECT ON servicios_tablas.servicios TO Progra_PAR;
GRANT SELECT ON servicios_tablas.productos TO Progra_PAR;
GRANT SELECT ON servicios_tablas.servicios_productos TO Progra_PAR;
GRANT INSERT, UPDATE ON servicios_tablas.servicios_productos TO Progra_PAR;
GRANT UPDATE ON servicios_tablas.productos TO Progra_PAR;
GRANT REFERENCES ON servicios_tablas.servicios TO Progra_PAR;
GRANT REFERENCES ON servicios_tablas.productos TO Progra_PAR;

-- Permisos sobre tablas de citas_tablas
GRANT SELECT ON citas_tablas.citas TO Progra_PAR;
GRANT SELECT ON citas_tablas.citas_servicios TO Progra_PAR;
GRANT SELECT ON citas_tablas.facturas TO Progra_PAR;
GRANT INSERT, UPDATE, DELETE ON citas_tablas.citas TO Progra_PAR;
GRANT INSERT, UPDATE, DELETE ON citas_tablas.citas_servicios TO Progra_PAR;
GRANT INSERT, UPDATE ON citas_tablas.facturas TO Progra_PAR;
GRANT REFERENCES ON citas_tablas.citas TO Progra_PAR;
GRANT REFERENCES ON citas_tablas.citas_servicios TO Progra_PAR;
GRANT REFERENCES ON citas_tablas.facturas TO Progra_PAR;

-- Verificación del diccionario de datos
SELECT username, default_tablespace, temporary_tablespace
FROM dba_users
WHERE username IN ('USUARIOS_TABLAS', 'CITAS_TABLAS', 'SERVICIOS_TABLAS', 'PROGRA_PAR');

SELECT *
FROM dba_ts_quotas
WHERE username IN ('USUARIOS_TABLAS', 'CITAS_TABLAS', 'SERVICIOS_TABLAS', 'PROGRA_PAR');

SELECT owner, object_name, object_type
FROM dba_objects
WHERE owner IN ('USUARIOS_TABLAS', 'CITAS_TABLAS', 'SERVICIOS_TABLAS')
ORDER BY 1, 3;

SELECT *
FROM dba_sys_privs
WHERE grantee IN ('USUARIOS_TABLAS', 'CITAS_TABLAS', 'SERVICIOS_TABLAS', 'PROGRA_PAR');

SELECT *
FROM dba_tab_privs
WHERE grantee = 'PROGRA_PAR';
