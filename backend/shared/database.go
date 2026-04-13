/*
database.go — Conexión global a MariaDB.

FLUJO DE DATOS:
 1. main.go llama a InitDB() al arrancar el servidor.
 2. Se construye un DSN (Data Source Name) usando variables de entorno o valores
    por defecto para desarrollo local.
 3. sql.Open("mysql", dsn) crea un pool de conexiones (no abre conexiones todavía).
 4. DB.Ping() verifica que la base de datos responde.
 5. La variable global DB queda disponible para todos los paquetes (features/).
    Cada feature recibe DB como parámetro en su RegisterRoutes(mux, DB),
    y lo inyecta en su Service struct.

El driver "github.com/go-sql-driver/mysql" es compatible con MariaDB.
Se usa parseTime=true para que TIME/DATETIME se escaneen como time.Time,
y charset=utf8mb4 para soportar caracteres especiales y emojis.
*/
package shared

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL/MariaDB — se registra automáticamente con database/sql
)

// DB es el pool global de conexiones. Se inicializa una sola vez en InitDB()
// y se comparte entre todos los features. database/sql maneja el pooling internamente.
var DB *sql.DB

// InitDB construye el DSN desde variables de entorno y abre el pool de conexiones.
// Se llama una sola vez desde main.go antes de registrar cualquier ruta.
// Si la conexión falla, el servidor se detiene con log.Fatalf (no tiene sentido
// correr sin base de datos).
func InitDB() {
	// Variables de entorno permiten configurar la conexión en producción
	// sin modificar el código. Los valores por defecto son para desarrollo local.
	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "3306")
	dbname := envOr("DB_NAME", "Patitas67D")
	user := envOr("DB_USER", "demopar")
	pass := envOr("DB_PASS", "demopar99")

	// Formato DSN del driver mysql: user:pass@tcp(host:port)/dbname?params
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, pass, host, port, dbname)

	var err error
	// sql.Open NO abre una conexión real; solo valida el DSN y prepara el pool.
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	// Ping() fuerza una conexión real para verificar que el servidor responde.
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al verificar conexión: %v", err)
	}
	log.Println("Conexión exitosa a la base de datos MariaDB.")
}

// envOr devuelve el valor de la variable de entorno key, o fallback si está vacía.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
