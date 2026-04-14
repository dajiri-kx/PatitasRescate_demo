<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/MariaDB-10.5+-003545?style=for-the-badge&logo=mariadb&logoColor=white" />
  <img src="https://img.shields.io/badge/Bootstrap-5.3-7952B3?style=for-the-badge&logo=bootstrap&logoColor=white" />
  <img src="https://img.shields.io/badge/Stripe-Payments-635BFF?style=for-the-badge&logo=stripe&logoColor=white" />
</p>

# 🐾 Patitas al Rescate

> **Sistema integral de gestión para clínica veterinaria** — Plataforma web full-stack que permite a dueños de mascotas agendar citas, gestionar el historial de sus mascotas y pagar facturas en línea, con portales dedicados para veterinarios y administradores.

---

## 📑 Tabla de Contenidos

- [Descripción General](#-descripción-general)
- [Stack Tecnológico](#-stack-tecnológico)
- [Arquitectura del Sistema](#-arquitectura-del-sistema)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Base de Datos](#-base-de-datos)
- [API Endpoints](#-api-endpoints)
- [Roles y Permisos](#-roles-y-permisos)
- [Funcionalidades](#-funcionalidades)
- [Instalación y Configuración](#-instalación-y-configuración)
- [Variables de Entorno](#-variables-de-entorno)
- [Seguridad](#-seguridad)

---

## 📖 Descripción General

**Patitas al Rescate** es un sistema de gestión veterinaria diseñado para clínicas de cuidado animal que ofrecen servicios de salud y estética. La plataforma centraliza la operación del negocio en tres portales:

| Portal | Descripción |
|--------|-------------|
| **🧑 Cliente** | Registro de mascotas, agendamiento de citas, visualización de facturas y pago en línea vía Stripe |
| **🩺 Veterinario** | Dashboard personal con citas asignadas del día, gestión de estados de citas |
| **⚙️ Administrador** | CRUD completo de servicios y veterinarios, vista de todos los clientes y citas del sistema, métricas del dashboard |

---

## 🛠 Stack Tecnológico

| Capa | Tecnología | Detalle |
|------|-----------|---------|
| **Backend** | Go 1.25 | Servidor HTTP con `net/http` (ServeMux nativo) |
| **Frontend** | HTML5 / CSS3 / JavaScript | Vanilla JS sin frameworks, Bootstrap 5.3.3 |
| **Base de Datos** | MariaDB 10.5+ | Compatible con MySQL 8.0+ |
| **Autenticación** | Gorilla Sessions | Cookies firmadas con HMAC, bcrypt para contraseñas |
| **Pagos** | Stripe Checkout API v82 | Moneda: CRC (Colón Costarricense) |
| **Driver SQL** | go-sql-driver/mysql | Conexión pool con `database/sql` |

---

## 🏗 Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────┐
│                      NAVEGADOR                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Cliente     │  │ Veterinario  │  │    Admin      │  │
│  │   Portal      │  │   Portal     │  │    Panel      │  │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  │
└─────────┼──────────────────┼──────────────────┼─────────┘
          │   HTTP + Cookies │                  │
          ▼──────────────────▼──────────────────▼
┌─────────────────────────────────────────────────────────┐
│                 GO SERVER (:8080)                         │
│                                                          │
│  ┌─────────┐  ┌──────────────────────────────────────┐  │
│  │  CORS   │→ │          HTTP Router (ServeMux)       │  │
│  │Middleware│  │                                      │  │
│  └─────────┘  │  /api/auth/*     → Auth Handlers     │  │
│               │  /api/citas/*    → Citas Handlers     │  │
│               │  /api/mascotas/* → Mascotas Handlers  │  │
│               │  /api/facturas/* → Facturas Handlers  │  │
│               │  /api/checkout/* → Stripe Handlers    │  │
│               │  /api/contacto   → Contacto Handler   │  │
│               │  /api/admin/*    → Admin Handlers     │  │
│               │  /api/vet/*      → Vet Handlers       │  │
│               │  /frontend/*     → Static Files       │  │
│               └──────────────────────────────────────┘  │
│                           │                              │
│               ┌───────────▼───────────┐                  │
│               │   Session Middleware   │                  │
│               │  RequireAuth           │                  │
│               │  RequireAdmin (rol=0)  │                  │
│               │  RequireVeterinario    │                  │
│               │       (rol=2)          │                  │
│               └───────────┬───────────┘                  │
└───────────────────────────┼──────────────────────────────┘
                            │ SQL (Prepared Statements)
                            ▼
                ┌───────────────────────┐
                │     MariaDB / MySQL    │
                │    patitas_rescate     │
                │                       │
                │  10 Tablas             │
                │  12 Funciones          │
                │   2 Procedimientos     │
                └───────────────────────┘
```

### Flujo de una Petición

1. El navegador envía una petición HTTP con `credentials: 'include'` (cookie de sesión)
2. El middleware CORS valida el origen (si está configurado)
3. El router despacha al handler correspondiente
4. El handler ejecuta `RequireAuth()` para extraer la sesión del cookie
5. El service ejecuta queries parametrizadas contra MariaDB
6. La respuesta se devuelve en formato JSON estandarizado

**Formato de respuesta estándar:**
```json
{
  "ok": true,
  "data": { ... },
  "message": "Operación exitosa"
}
```

---

## 📁 Estructura del Proyecto

```
PatitasRescate/
│
├── index.html                    # Landing page raíz
├── Planificador.md               # Documento de planificación del MVP
│
├── backend/
│   ├── main.go                   # Punto de entrada del servidor
│   ├── go.mod                    # Dependencias de Go
│   ├── shared/
│   │   ├── database.go           # Pool de conexiones MariaDB
│   │   ├── middleware.go         # Middleware CORS
│   │   ├── response.go          # Helpers de respuesta JSON
│   │   ├── session.go           # Sesiones con Gorilla (cookie HMAC)
│   │   └── stripe.go            # Inicialización de Stripe API
│   └── features/
│       ├── auth/                 # Registro, login, logout, check-session
│       ├── citas/                # Agendar, cancelar, listar citas
│       ├── mascotas/             # CRUD de mascotas
│       ├── facturas/             # Consulta de facturas
│       ├── checkout/             # Integración Stripe Checkout
│       ├── contacto/             # Formulario de contacto público
│       ├── admin/                # Panel administrativo
│       ├── veterinario/          # Portal veterinario
│       └── clientes/             # Servicio de perfil de cliente
│
├── frontend/
│   ├── shared/
│   │   ├── css/
│   │   │   ├── admin.css         # Estilos del panel admin
│   │   │   └── vet.css           # Estilos del portal veterinario
│   │   └── js/
│   │       ├── api.js            # Wrapper centralizado de fetch
│   │       ├── auth.js           # Estado de autenticación (localStorage)
│   │       ├── components.js     # Header/Footer dinámicos por rol
│   │       ├── admin-layout.js   # Layout y guards del admin
│   │       └── vet-layout.js     # Layout y guards del veterinario
│   └── features/
│       ├── home/                 # Página de inicio
│       ├── auth/                 # Login y registro
│       ├── dashboard/            # Dashboard del cliente
│       ├── mis-citas/            # Mis citas (cliente)
│       ├── mis-mascotas/         # Mis mascotas (cliente)
│       ├── mis-facturas/         # Mis facturas (cliente)
│       ├── citas/                # Agendar / Cancelar cita
│       ├── mascotas/             # Agregar mascota
│       ├── perfil/               # Perfil del usuario
│       ├── servicios/            # Catálogo de servicios
│       ├── contactenos/          # Formulario de contacto
│       ├── ubicacion/            # Página de ubicación
│       ├── pago-felicidades/     # Confirmación de pago exitoso
│       ├── admin/                # Panel admin (4 secciones)
│       └── veterinario/          # Portal veterinario
│
└── dbMariaDB/
    ├── patitas_rescate.sql       # Script completo de BD
    └── README.md                 # Instrucciones de la BD
```

### Patrón de Código Backend

Cada feature sigue el patrón **Handler → Service → Database**:

```go
// handlers.go — Decodifica HTTP, valida, llama al service
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
    svc := &Service{DB: db}
    mux.HandleFunc("POST /api/feature/action", actionHandler(svc))
}

// service.go — Lógica de negocio y queries SQL
func (s *Service) Action(ctx context.Context, params ...) (Result, error) {
    // Queries parametrizadas contra s.DB
}
```

---

## 🗄 Base de Datos

### Diagrama Entidad-Relación

```
┌──────────────┐     1:N     ┌──────────────┐     1:N     ┌──────────────┐
│   CLIENTES   │────────────▶│   MASCOTAS   │────────────▶│    CITAS      │
│              │             │              │             │              │
│ ID_CLIENTE   │             │ ID_MASCOTA   │             │ ID_CITA      │
│ DIDENTIDAD   │             │ ID_CLIENTE   │             │ ID_MASCOTA   │
│ NOMBRE       │             │ NOMBRE       │             │ ID_VETERINARIO│
│ APELLIDO     │             │ ESPECIE      │             │ FECHA_CITA   │
│ EMAIL        │             │ RAZA         │             │ ESTADO       │
│ TELEFONO     │             │ MESES        │             └───────┬──────┘
│ DIRECCION    │             └──────────────┘                     │
│ FECHA_REGISTRO│                                                 │ N:M
└──────┬───────┘                                                  │
       │ 1:1                                              ┌───────▼──────┐
┌──────▼───────┐     ┌──────────────┐                    │CITAS_SERVICIOS│
│   USUARIOS   │     │ VETERINARIOS │                    │              │
│              │     │              │                    │ ID_CITA      │
│ ID_USUARIO   │     │ ID_VETERINARIO│◀──────────────────│ ID_SERVICIO  │
│ ID_CLIENTE   │     │ NOMBRE       │                    │ ID_FACTURA   │
│ CORREO       │     │ ESPECIALIDAD │                    └──┬────┬──────┘
│ CONTRASENA   │     │ TELEFONO     │                       │    │
│ ROL (0,1,2)  │     │ CORREO       │                       │    │
│ ID_VETERINARIO│     │ ROL          │                ┌─────▼┐  ┌▼────────┐
└──────────────┘     └──────────────┘                │SERVICIOS│ │FACTURAS │
                                                     │        │ │         │
┌──────────────┐     ┌──────────────────┐            │ID_SERV │ │ID_FACTURA│
│  PRODUCTOS   │     │SERVICIOS_PRODUCTOS│            │NOMBRE  │ │FECHA    │
│              │◀────│                  │────────────▶│PRECIO  │ │TOTAL    │
│ ID_PRODUCTO  │     │ ID_PRODUCTO      │            │DURACIÓN│ │ESTADO   │
│ NOMBRE       │     │ ID_SERVICIO      │            │CATEGORÍA│ │STRIPE_ID│
│ CATEGORIA    │     │ UNIDADES         │            └────────┘ └─────────┘
│ PRECIO       │     │ CANT_CONSUMIDA   │
│ STOCK        │     └──────────────────┘
└──────────────┘
```

### Tablas (10)

| Tabla | Descripción | Registros Seed |
|-------|-------------|----------------|
| `CLIENTES` | Dueños de mascotas registrados | 1 (admin) |
| `USUARIOS` | Credenciales de acceso (bcrypt) | 1 (admin) |
| `VETERINARIOS` | Veterinarios y estilistas | 3 |
| `MASCOTAS` | Mascotas vinculadas a clientes | — |
| `SERVICIOS` | Catálogo de servicios | 4 |
| `PRODUCTOS` | Insumos con control de stock | 3 |
| `SERVICIOS_PRODUCTOS` | Relación servicios ↔ productos | 3 |
| `CITAS` | Citas agendadas | — |
| `CITAS_SERVICIOS` | Relación citas ↔ servicios ↔ facturas | — |
| `FACTURAS` | Facturas generadas por citas | — |

### Funciones Almacenadas (12)

| Función | Propósito |
|---------|-----------|
| `existeCliente()` | Verifica existencia de cliente |
| `existeMascota()` | Verifica existencia de mascota |
| `mascotaPerteneceACliente()` | Valida que la mascota es del cliente |
| `mascotaTieneCitaActivaMismaFecha()` | Detecta conflictos de horario |
| `existeVeterinario()` | Verifica existencia de veterinario |
| `veterinarioDisponible()` | Valida disponibilidad del vet |
| `existeServicio()` | Verifica existencia de servicio |
| `existeProducto()` | Verifica existencia de producto |
| `calcularTotalServicios()` | Suma precios de servicios seleccionados |
| `existeCedula()` | Previene duplicados de cédula |
| `existeCorreo()` | Previene duplicados de correo |

### Procedimientos Almacenados (2)

| Procedimiento | Descripción |
|---------------|-------------|
| `registrarCliente()` | Transacción: INSERT en `CLIENTES` + `USUARIOS` con validaciones |
| `agendarCita()` | Transacción compleja de 11 pasos: validaciones → crear cita → vincular servicios → descontar stock → generar factura |

---

## 🔌 API Endpoints

### Autenticación (`/api/auth/`)

| Método | Ruta | Auth | Descripción |
|--------|------|:----:|-------------|
| `POST` | `/api/auth/register` | ❌ | Registrar nuevo cliente |
| `POST` | `/api/auth/login` | ❌ | Iniciar sesión (email + contraseña) |
| `POST` | `/api/auth/logout` | ✅ | Cerrar sesión |
| `GET` | `/api/auth/check-session` | ✅ | Validar sesión activa |

### Citas (`/api/citas/`)

| Método | Ruta | Auth | Descripción |
|--------|------|:----:|-------------|
| `GET` | `/api/citas` | ✅ | Todas las citas del cliente |
| `GET` | `/api/citas/activas` | ✅ | Citas activas (para cancelación) |
| `GET` | `/api/citas/veterinarios` | ✅ | Lista de veterinarios |
| `GET` | `/api/citas/servicios` | ✅ | Servicios (filtrable por `?categoria=`) |
| `POST` | `/api/citas/agendar` | ✅ | Agendar cita |
| `POST` | `/api/citas/cancelar` | ✅ | Cancelar cita activa |

### Mascotas (`/api/mascotas/`)

| Método | Ruta | Auth | Descripción |
|--------|------|:----:|-------------|
| `GET` | `/api/mascotas` | ✅ | Todas las mascotas del cliente |
| `GET` | `/api/mascotas/nombres` | ✅ | IDs + nombres (para dropdowns) |
| `POST` | `/api/mascotas/agregar` | ✅ | Registrar nueva mascota |

### Facturas (`/api/facturas/`)

| Método | Ruta | Auth | Descripción |
|--------|------|:----:|-------------|
| `GET` | `/api/facturas` | ✅ | Facturas del cliente |

### Pagos (`/api/checkout/`)

| Método | Ruta | Auth | Descripción |
|--------|------|:----:|-------------|
| `POST` | `/api/checkout/crear-sesion` | ✅ | Crear sesión de Stripe Checkout |
| `POST` | `/api/checkout/verificar` | ✅ | Verificar pago completado |

### Contacto (`/api/contacto/`)

| Método | Ruta | Auth | Descripción |
|--------|------|:----:|-------------|
| `POST` | `/api/contacto` | ❌ | Enviar formulario de contacto |

### Administración (`/api/admin/`) — Solo rol Admin

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/admin/stats` | Métricas del dashboard |
| `GET/POST` | `/api/admin/servicios` | Listar / Crear servicios |
| `POST` | `/api/admin/servicios/editar` | Editar servicio |
| `POST` | `/api/admin/servicios/eliminar` | Eliminar servicio |
| `GET/POST` | `/api/admin/veterinarios` | Listar / Crear veterinarios |
| `POST` | `/api/admin/veterinarios/editar` | Editar veterinario |
| `POST` | `/api/admin/veterinarios/eliminar` | Eliminar veterinario |
| `GET` | `/api/admin/clientes` | Listar clientes |
| `GET` | `/api/admin/citas` | Todas las citas del sistema |
| `POST` | `/api/admin/citas/estado` | Cambiar estado de cita |

### Portal Veterinario (`/api/vet/`) — Solo rol Veterinario

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/vet/stats` | Métricas personales del veterinario |
| `GET` | `/api/vet/citas` | Citas asignadas al veterinario |
| `POST` | `/api/vet/citas/estado` | Actualizar estado (solo Completada/Cancelada) |

---

## 👥 Roles y Permisos

```
┌─────────────────────────────────────────────────────────────────┐
│                        ROLES DEL SISTEMA                         │
├──────────┬──────┬───────────────────────────────────────────────┤
│   Rol    │  ID  │              Permisos                          │
├──────────┼──────┼───────────────────────────────────────────────┤
│ Admin    │  0   │ • Todo lo del cliente                          │
│          │      │ • CRUD de servicios y veterinarios             │
│          │      │ • Ver todos los clientes (lectura)             │
│          │      │ • Ver y gestionar todas las citas              │
│          │      │ • Dashboard con métricas del sistema           │
├──────────┼──────┼───────────────────────────────────────────────┤
│ Cliente  │  1   │ • Registrarse e iniciar sesión                 │
│          │      │ • Gestionar mascotas (agregar, ver)            │
│          │      │ • Agendar y cancelar citas                     │
│          │      │ • Ver facturas y pagar vía Stripe              │
│          │      │ • Editar perfil personal                       │
├──────────┼──────┼───────────────────────────────────────────────┤
│Veterinario│  2  │ • Ver citas asignadas                          │
│          │      │ • Marcar citas como Completada/Cancelada       │
│          │      │ • Dashboard con métricas personales            │
└──────────┴──────┴───────────────────────────────────────────────┘
```

**Enforcement:**
- **Server-side:** Middleware `RequireAuth()`, `RequireAdmin()`, `RequireVeterinario()` validan rol desde la cookie de sesión
- **Client-side:** Guards en JavaScript (`Auth.requireAuth()`, `adminRequireAuth()`, `vetRequireAuth()`) redirigen a login si no hay sesión

---

## ✨ Funcionalidades

### 🏠 Público
- Landing page con banner, testimonios y tarjeta de contacto
- Catálogo de servicios con filtro por categoría
- Formulario de contacto
- Página de ubicación
- Registro e inicio de sesión

### 🧑 Portal Cliente
- **Mascotas** — Agregar y visualizar mascotas (nombre, especie, raza, edad)
- **Citas** — Agendar citas seleccionando mascota, veterinario, servicios y fecha/hora
- **Cancelar Citas** — Cancelar citas activas
- **Facturas** — Ver historial de facturas con estado (Pendiente/Pagada)
- **Pagos** — Pagar facturas pendientes vía Stripe Checkout
- **Perfil** — Ver y editar información personal

### 🩺 Portal Veterinario
- Dashboard con métricas: citas del día, pendientes, completadas
- Agenda con todas las citas asignadas incluyendo detalles de mascota, cliente y servicios
- Actualizar estado de citas (Completada / Cancelada)

### ⚙️ Panel Administrador
- Dashboard con métricas del sistema (servicios, veterinarios, clientes, citas activas)
- **Servicios** — CRUD completo (nombre, descripción, precio, duración, categoría)
- **Veterinarios** — CRUD completo (nombre, cédula, especialidad, teléfono, correo, rol)
- **Clientes** — Vista de lectura de todos los clientes registrados
- **Citas** — Vista de todas las citas con control de estado

### 💳 Flujo de Pago

```
Cliente ve factura PENDIENTE
        │
        ▼
POST /api/checkout/crear-sesion
        │
        ▼
Stripe genera URL de checkout ───▶ Redirige al cliente
        │
        ▼
Cliente paga en Stripe
        │
        ▼
Redirige a /pago-felicidades/?session_id=X
        │
        ▼
POST /api/checkout/verificar
        │
        ▼
Factura actualizada a PAGADA ✅
```

---

## 🚀 Instalación y Configuración

### Prerrequisitos

- **Go** 1.25+ → [golang.org/dl](https://golang.org/dl/)
- **MariaDB** 10.5+ o **MySQL** 8.0+
- **Stripe Account** (opcional, para pagos) → [stripe.com](https://stripe.com)

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/PatitasRescate.git
cd PatitasRescate
```

### 2. Configurar la base de datos

```bash
# Conectarse a MariaDB/MySQL como root
mysql -u root -p

# Ejecutar el script completo
source dbMariaDB/patitas_rescate.sql;
```

Esto crea:
- Base de datos `patitas_rescate`
- 10 tablas con integridad referencial
- 12 funciones y 2 procedimientos almacenados
- Usuario de BD `Progra_PAR` con permisos adecuados
- Datos semilla (3 veterinarios, 4 servicios, 3 productos, 1 admin)

### 3. Configurar variables de entorno

```bash
# Base de datos
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=Patitas67D
export DB_USER=demopar
export DB_PASS=demopar99

# Sesiones (cambiar en producción)
export SESSION_KEY="tu-clave-secreta-de-32-bytes-min"

# Stripe (opcional)
export STRIPE_SECRET_KEY=sk_test_...

# CORS (opcional, para desarrollo con frontend separado)
export CORS_ORIGIN=http://localhost:5500
```

### 4. Ejecutar el servidor

```bash
cd backend
go run main.go
```

El servidor inicia en `http://localhost:8080` y sirve:
- **API:** `http://localhost:8080/api/*`
- **Frontend:** `http://localhost:8080/frontend/features/*`
- **Landing:** `http://localhost:8080/`

### 5. Credenciales por defecto

| Rol | Email | Contraseña |
|-----|-------|------------|
| Admin | `admin@patitas.com` | `Admin123!` |

---

## ⚙️ Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Host de la base de datos |
| `DB_PORT` | `3306` | Puerto de MariaDB/MySQL |
| `DB_NAME` | `Patitas67D` | Nombre de la base de datos |
| `DB_USER` | `demopar` | Usuario de la base de datos |
| `DB_PASS` | `demopar99` | Contraseña de la base de datos |
| `PORT` | `8080` | Puerto del servidor HTTP |
| `SESSION_KEY` | *(hardcoded)* | Clave HMAC para firmar cookies de sesión |
| `STRIPE_SECRET_KEY` | *(ninguno)* | API key de Stripe (modo test o live) |
| `CORS_ORIGIN` | *(deshabilitado)* | Origen permitido para CORS |
| `BASE_URL` | `http://localhost:8080` | URL base para redirecciones de Stripe |

---

## 🔒 Seguridad

### Autenticación y Sesiones
- Contraseñas hasheadas con **bcrypt** (costo 10)
- Cookies de sesión **HttpOnly** (inaccesibles desde JavaScript)
- Cookies firmadas con **HMAC** (detección de manipulación)
- Expiración de sesión a las **24 horas**
- **SameSite=Lax** como mitigación básica de CSRF

### Validación de Datos
- Consultas SQL **parametrizadas** en todo el backend (prevención de SQL injection)
- Validación server-side con **regex** (cédula, teléfono, email, complejidad de contraseña)
- **Escape de HTML** (`escapeHtml()`) en el frontend para contenido generado por usuarios
- **Whitelist** de valores válidos (estados de citas, roles)

### Control de Acceso
- Verificación de **propiedad** en todas las operaciones del cliente (mascota → cita → factura)
- `idCliente` se extrae de la **cookie de sesión**, nunca del body de la petición
- Middleware de rol server-side (`RequireAdmin`, `RequireVeterinario`)
- Guards client-side con redirección automática

### Infraestructura
- API key de Stripe almacenada en **variable de entorno**
- Credenciales de BD configurables vía **variables de entorno**
- CORS configurable y **deshabilitado por defecto**

---

<p align="center">
  Desarrollado con ❤️ para el cuidado de nuestras mascotas
</p>
