package shared

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "3306")
	dbname := envOr("DB_NAME", "patitas_rescate")
	user := envOr("DB_USER", "Progra_PAR")
	pass := envOr("DB_PASS", "PrograPAR_2026")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, pass, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al verificar conexión: %v", err)
	}
	log.Println("Conexión exitosa a la base de datos MariaDB.")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
