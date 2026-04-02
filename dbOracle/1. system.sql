-- desde System

-- Creación del usuario Usuarios_Tablas
create user Usuarios_Tablas identified by Usuarios_Tablas
default tablespace users 
temporary tablespace temp
quota unlimited on users;

-- Otorgando privilegios de sistema
grant create session, create table to Usuarios_Tablas;
grant create sequence to Usuarios_Tablas;
grant create trigger to Usuarios_Tablas;

-- Creación del usuario Citas_Tablas
create user Citas_Tablas identified by Citas_Tablas
default tablespace users 
temporary tablespace temp
quota unlimited on users;

-- Otorgando privilegios de sistema
grant create session, create table to Citas_Tablas;

-- Creación del usuario Servicios_Tablas
create user Servicios_Tablas identified by Servicios_Tablas
default tablespace users 
temporary tablespace temp
quota unlimited on users;

-- Otorgando privilegios de sistema
grant create session, create table to Servicios_Tablas;

-- creación de usuario Progra.
CREATE USER Progra_PAR IDENTIFIED BY Progra_PAR
DEFAULT TABLESPACE users
TEMPORARY TABLESPACE temp
QUOTA UNLIMITED ON users;

-- Otorgar privilegios básicos
GRANT CREATE SESSION TO Progra_PAR;
GRANT CONNECT, RESOURCE TO Progra_PAR;

-- Permitir la creación de procedimientos, funciones y triggers
GRANT CREATE PROCEDURE, CREATE TRIGGER TO Progra_PAR;

-- Otorgar privilegios para consultar tablas
GRANT SELECT ON usuarios_tablas.clientes TO progra;
GRANT SELECT ON usuarios_tablas.mascotas TO progra;
GRANT SELECT ON usuarios_tablas.veterinarios TO progra;
GRANT SELECT ON servicios_tablas.servicios TO progra;
GRANT SELECT ON servicios_tablas.productos TO progra;
GRANT SELECT ON citas_tablas.citas TO progra;

-- Otorgar privilegios para insertar y actualizar registros
GRANT INSERT, UPDATE ON citas_tablas.citas TO progra;
GRANT INSERT, UPDATE ON citas_tablas.citas_servicios TO progra;
GRANT UPDATE ON servicios_tablas.productos TO progra;
GRANT INSERT ON usuarios_tablas.clientes TO progra;
GRANT INSERT ON usuarios_tablas.usuarios TO progra;
GRANT INSERT ON citas_tablas.facturas TO progra;

-- Revisar diccionario de datos
select username, default_tablespace, temporary_tablespace
from dba_users
where username like '%TABLAS';

select  *
from dba_ts_quotas
where username like '%TABLAS';

select owner, object_name, object_type
from dba_objects
where owner like '%TABLAS'
order by 1,3;

select *
from dba_sys_privs
where grantee like '%TABLAS';

select *
from dba_tab_privs
where grantee like '%TABLAS';

SELECT constraint_name, table_name, r_constraint_name
FROM user_constraints
WHERE table_name = 'CITAS_SERVICIOS' AND constraint_type = 'R';


SELECT * 
FROM DBA_TAB_PRIVS 
WHERE GRANTEE = 'PROGRA_PAR';

SELECT * 
FROM DBA_SYS_PRIVS 
WHERE GRANTEE = 'PROGRA_PAR';

SELECT *
FROM ALL_TAB_PRIVS 
WHERE GRANTEE = 'PROGRA_PAR';


GRANT SELECT ON USUARIOS_TABLAS.CLIENTES TO Progra_PAR;
GRANT SELECT ON USUARIOS_TABLAS.USUARIOS TO Progra_PAR;
GRANT SELECT ON USUARIOS_TABLAS.MASCOTAS TO Progra_PAR;
GRANT SELECT ON USUARIOS_TABLAS.VETERINARIOS TO Progra_PAR;

GRANT SELECT ON CITAS_TABLAS.CITAS TO Progra_PAR;
GRANT SELECT ON CITAS_TABLAS.CITAS_SERVICIOS TO Progra_PAR;
GRANT SELECT ON CITAS_TABLAS.FACTURAS TO Progra_PAR;

GRANT SELECT ON SERVICIOS_TABLAS.SERVICIOS TO Progra_PAR;
GRANT SELECT ON SERVICIOS_TABLAS.PRODUCTOS TO Progra_PAR;
GRANT SELECT ON SERVICIOS_TABLAS.SERVICIOS_PRODUCTOS TO Progra_PAR;

--
GRANT INSERT ON USUARIOS_TABLAS.CLIENTES TO Progra_PAR;
GRANT INSERT ON USUARIOS_TABLAS.USUARIOS TO Progra_PAR;

GRANT INSERT ON CITAS_TABLAS.CITAS TO Progra_PAR;
GRANT INSERT ON CITAS_TABLAS.FACTURAS TO Progra_PAR;
GRANT INSERT ON CITAS_TABLAS.CITAS_SERVICIOS TO Progra_PAR;

GRANT INSERT ON SERVICIOS_TABLAS.SERVICIOS_PRODUCTOS TO Progra_PAR;

--
GRANT UPDATE ON SERVICIOS_TABLAS.PRODUCTOS TO Progra_PAR;
GRANT UPDATE ON CITAS_TABLAS.FACTURAS TO Progra_PAR;
GRANT UPDATE ON CITAS_TABLAS.CITAS_SERVICIOS TO Progra_PAR;

--
GRANT REFERENCES ON USUARIOS_TABLAS.CLIENTES TO Progra_PAR;
GRANT REFERENCES ON USUARIOS_TABLAS.MASCOTAS TO Progra_PAR;
GRANT REFERENCES ON USUARIOS_TABLAS.VETERINARIOS TO Progra_PAR;
GRANT REFERENCES ON USUARIOS_TABLAS.USUARIOS TO Progra_PAR;

GRANT REFERENCES ON CITAS_TABLAS.CITAS TO Progra_PAR;
GRANT REFERENCES ON CITAS_TABLAS.FACTURAS TO Progra_PAR;
GRANT REFERENCES ON CITAS_TABLAS.CITAS_SERVICIOS TO Progra_PAR;

GRANT REFERENCES ON SERVICIOS_TABLAS.SERVICIOS TO Progra_PAR;
GRANT REFERENCES ON SERVICIOS_TABLAS.PRODUCTOS TO Progra_PAR;
GRANT REFERENCES ON SERVICIOS_TABLAS.SERVICIOS_PRODUCTOS TO Progra_PAR;

GRANT CREATE PROCEDURE TO Progra_PAR;
GRANT CREATE TRIGGER TO Progra_PAR;
GRANT SELECT ON servicios_tablas.servicios_productos TO Progra_PAR;


GRANT SELECT ON servicios_tablas.servicios_productos TO Progra_PAR;
GRANT SELECT ON servicios_tablas.servicios TO Progra_PAR;
GRANT SELECT ON servicios_tablas.productos TO Progra_PAR;
GRANT CREATE TRIGGER TO Progra_PAR;

-- /24/04/25
GRANT CREATE TRIGGER TO Usuarios_Tablas;
GRANT CREATE SEQUENCE TO Usuarios_Tablas;

-- /24/04/25
GRANT CREATE TRIGGER TO Citas_Tablas;
GRANT CREATE SEQUENCE TO Citas_Tablas;

-- /24/04/25
GRANT CREATE TRIGGER TO Servicios_Tablas;
GRANT CREATE SEQUENCE TO Servicios_Tablas;

/*
-- Eliminar el usuario Usuarios_Tablas y todos sus objetos
DROP USER Usuarios_Tablas CASCADE;
-- Eliminar el usuario Citas_Tablas y todos sus objetos
DROP USER Citas_Tablas CASCADE;

-- Eliminar el usuario Servicios_Tablas y todos sus objetos
DROP USER Progra CASCADE;

-- Eliminar el usuario Progra_PAR y todos sus objetos
DROP USER Progra_PAR CASCADE;

-- Eliminar el usuario Progra y todos sus objetos
DROP USER Progra CASCADE;*/