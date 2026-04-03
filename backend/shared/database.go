package shared

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/sijms/go-ora/v2"
)

var DB *sql.DB

func InitDB() {
	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "1521")
	service := envOr("DB_SERVICE", "xe")
	user := envOr("DB_USER", "Progra_PAR")
	pass := envOr("DB_PASS", "PrograPAR_2026")

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		url.PathEscape(user), url.PathEscape(pass), host, port, service)

	var err error
	DB, err = sql.Open("oracle", dsn)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al verificar conexión: %v", err)
	}
	log.Println("Conexión exitosa a la base de datos Oracle.")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
