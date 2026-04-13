# Scripts de Base de Datos — MariaDB

## Requisitos

- MariaDB 10.5 o superior
- Acceso con usuario `root` o con privilegios de `CREATE DATABASE` y `CREATE USER`

## Ejecución

Un solo script unificado crea toda la base de datos:

```bash
mysql -u root -p < patitas_rescate.sql
```

> Los scripts numerados (`1. crear_base_datos.sql`, etc.) son la versión anterior
> separada por pasos. Se conservan como referencia pero **`patitas_rescate.sql` es
> el script definitivo** y el único que se necesita ejecutar.

## Qué hace el script

### Base de datos y usuario
Crea la base de datos `patitas_rescate` (utf8mb4) y el usuario de aplicación
`Progra_PAR` con permisos SELECT, INSERT, UPDATE y DELETE.

### Tablas de usuarios
- **CLIENTES** — datos personales de los dueños de mascotas
- **VETERINARIOS** — profesionales que atienden citas
- **MASCOTAS** — mascotas asociadas a cada cliente (FK → CLIENTES)
- **USUARIOS** — credenciales de autenticación (FK → CLIENTES, FK → VETERINARIOS)
  - `ROL`: 0 = Admin, 1 = Cliente (default), 2 = Veterinario
  - `ID_VETERINARIO`: nullable, solo para rol 2

### Tablas de catálogo
- **SERVICIOS** — catálogo de servicios con `CATEGORIA` para filtrar en el frontend
- **PRODUCTOS** — consumibles usados durante los servicios
- **SERVICIOS_PRODUCTOS** — relación N:N de qué productos consume cada servicio

### Tablas operativas
- **CITAS** — citas agendadas (estado: Activa → Completada | Cancelada)
- **FACTURAS** — generadas al agendar, con `ESTADO` (Pendiente → Pagada) y `STRIPE_SESSION_ID`
- **CITAS_SERVICIOS** — tabla puente que liga citas con servicios y con su factura

### Funciones y procedimientos
- Funciones de validación: `existeCliente()`, `existeMascota()`, `mascotaPerteneceACliente()`, etc.
- `registrarCliente()` — transacción: valida duplicados → inserta CLIENTES + USUARIOS
- `agendarCita()` — transacción: valida reglas → crea cita → registra servicios → descuenta stock → genera factura
- `actualizarStock()` / `calcularTotalServicios()`

### Datos semilla
- 3 veterinarios, 4 servicios (categoría "Estética"), 3 productos con stock, asociaciones servicio-producto
- 1 usuario administrador (`admin@patitas.com` / `Admin123!`, ROL=0)

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