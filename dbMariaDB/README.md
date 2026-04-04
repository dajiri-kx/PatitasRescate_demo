# Scripts de Base de Datos — MariaDB

## Requisitos

- MariaDB 10.5 o superior
- Acceso con usuario `root` o con privilegios de `CREATE DATABASE` y `CREATE USER`

## Orden de ejecución

Los scripts deben ejecutarse en orden numérico desde un cliente MariaDB con privilegios de administrador:

```bash
mysql -u root -p < "1. crear_base_datos.sql"
mysql -u root -p < "2. usuarios_tablas.sql"
mysql -u root -p < "3. servicios_tablas.sql"
mysql -u root -p < "4. citas_tablas.sql"
mysql -u root -p < "5. procedimientos.sql"
```

## Descripción de cada script

### 1. crear_base_datos.sql
Crea la base de datos `patitas_rescate` y el usuario de aplicación `Progra_PAR` con permisos SELECT, INSERT, UPDATE y DELETE sobre todas las tablas.

Equivale al `system.sql` de Oracle donde se creaban los schemas separados (`Usuarios_Tablas`, `Citas_Tablas`, `Servicios_Tablas`) y los grants cruzados. En MariaDB todo vive en una sola base de datos, así que este script es mucho más corto.

### 2. usuarios_tablas.sql
Crea las tablas de dominio de usuarios:
- **CLIENTES** — datos personales de los clientes
- **MASCOTAS** — mascotas asociadas a cada cliente
- **VETERINARIOS** — profesionales que atienden citas
- **USUARIOS** — credenciales de autenticación (correo + hash bcrypt)

Incluye datos semilla de 3 veterinarios.

### 3. servicios_tablas.sql
Crea las tablas de catálogo de servicios:
- **SERVICIOS** — servicios ofrecidos (cortes, uñas, etc.)
- **PRODUCTOS** — productos consumibles (shampoo, acondicionador, etc.)
- **SERVICIOS_PRODUCTOS** — relación de qué productos consume cada servicio y en qué cantidad

Incluye datos semilla de 4 servicios, 3 productos, y sus asociaciones.

### 4. citas_tablas.sql
Crea las tablas operativas:
- **CITAS** — citas agendadas con fecha, mascota y veterinario
- **FACTURAS** — facturas generadas al agendar una cita
- **CITAS_SERVICIOS** — tabla pivote que liga cada cita con sus servicios y su factura

### 5. procedimientos.sql
Crea funciones y procedimientos almacenados equivalentes a los de Oracle:

**Funciones de validación:**
- `existeCliente()`, `existeMascota()`, `existeVeterinario()`, `existeServicio()`, `existeProducto()`
- `mascotaPerteneceACliente()` — verifica propiedad mascota-cliente
- `mascotaTieneCitaActivaMismaFecha()` — evita citas duplicadas el mismo día
- `veterinarioDisponible()` — verifica que no haya conflicto de horario
- `existeCedula()`, `existeCorreo()` — validación de datos únicos en registro

**Funciones de cálculo:**
- `calcularTotalServicios()` — suma precios de servicios asociados a una cita

**Procedimientos principales:**
- `registrarCliente()` — valida datos únicos, inserta en CLIENTES y USUARIOS dentro de una transacción
- `agendarCita()` — valida todas las reglas de negocio, crea la cita, registra servicios, actualiza stock de productos, genera factura y la asocia. Todo en una transacción
- `actualizarStock()` — descuenta stock de un producto

## Diferencias con Oracle

| Concepto | Oracle | MariaDB |
|----------|--------|---------|
| Schemas | `USUARIOS_TABLAS.CLIENTES` | `CLIENTES` (todo en una BD) |
| Auto-incremento | Sequences + Triggers | `AUTO_INCREMENT` |
| Tipos texto | `VARCHAR2`, `CLOB` | `VARCHAR`, `TEXT` |
| Booleano en funciones | `RETURN BOOLEAN` | `RETURNS TINYINT(1)` |
| Obtener ID insertado | `RETURNING id INTO var` | `LAST_INSERT_ID()` |
| Errores en procedimientos | `RAISE_APPLICATION_ERROR(-200xx, msg)` | `SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = msg` |
| Parsear CSV | `SYS.ODCINUMBERLIST` + `CONNECT BY` | `WHILE` + `SUBSTRING_INDEX` |
| Filas afectadas | `SQL%ROWCOUNT` | `ROW_COUNT()` |
| Fecha actual | `SYSDATE` | `NOW()` |
| Truncar fecha | `TRUNC(fecha)` | `DATE(fecha)` |
| Formatear fecha | `TO_CHAR(fecha, 'YYYY-MM-DD')` | `DATE_FORMAT(fecha, '%Y-%m-%d')` |
| Nulo por defecto | `NVL(col, val)` | `IFNULL(col, val)` |

## Variables de entorno del backend Go

El backend espera estas variables (con valores por defecto):

| Variable | Default | Descripción |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Host de MariaDB |
| `DB_PORT` | `3306` | Puerto de MariaDB |
| `DB_NAME` | `patitas_rescate` | Nombre de la base de datos |
| `DB_USER` | `Progra_PAR` | Usuario de aplicación |
| `DB_PASS` | `PrograPAR_2026` | Contraseña del usuario |
