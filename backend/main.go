package main

import (
	"log"
	"net/http"
	"os"

	"patitas-backend/features/auth"
	"patitas-backend/features/citas"
	"patitas-backend/features/contacto"
	"patitas-backend/features/facturas"
	"patitas-backend/features/mascotas"
	"patitas-backend/shared"
)

func main() {
	shared.InitDB()
	defer shared.DB.Close()
	shared.InitSessionStore()

	mux := http.NewServeMux()

	// API routes
	auth.RegisterRoutes(mux, shared.DB)
	citas.RegisterRoutes(mux, shared.DB)
	mascotas.RegisterRoutes(mux, shared.DB)
	facturas.RegisterRoutes(mux, shared.DB)
	contacto.RegisterRoutes(mux)

	// Serve frontend static files
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend"))))

	// Root -> index.html (redirects to /frontend/features/home/)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../index.html")
	})

	handler := shared.CORS(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
