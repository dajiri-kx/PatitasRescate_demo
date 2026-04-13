/*
main.go — Punto de entrada del backend de Patitas al Rescate.

ORDEN DE INICIALIZACIÓN:
1. InitDB()           → Conecta a MariaDB (Patitas67D), verifica con Ping.
2. InitSessionStore() → Crea el CookieStore para sesiónes con gorilla/sessions.
3. InitStripe()       → Configura la API key de Stripe (opcional, falla silenciosamente).
4. Registrar rutas    → Cada feature registra sus endpoints en el mux.
5. Servir frontend    → Archivos estáticos desde ../frontend/.
6. Envolver con CORS  → shared.CORS añade headers para peticiones cross-origin.
7. ListenAndServe     → Puerto 8080 (o env PORT).

ARQUITECTURA:
El servidor Go sirve TANTO la API (/api/*) como los archivos estáticos
del frontend (/frontend/*). No hay servidor de frontend separado.
Esto simplifica el despliegue: un solo proceso sirve todo.

REGISTRO DE FEATURES:
Cada paquete de feature tiene su propio RegisterRoutes(mux, db).
El mux de Go 1.22+ soporta patterns como "GET /api/auth/login",
por lo que no se necesita un router externo (e.g., gorilla/mux).
*/
package main

import (
	"log"
	"net/http"
	"os"

	"patitas-backend/features/admin"
	"patitas-backend/features/auth"
	"patitas-backend/features/checkout"
	"patitas-backend/features/citas"
	"patitas-backend/features/contacto"
	"patitas-backend/features/facturas"
	"patitas-backend/features/mascotas"
	"patitas-backend/features/veterinario"
	"patitas-backend/shared"
)

func main() {
	// 1. Inicializar infraestructura compartida.
	shared.InitDB()
	defer shared.DB.Close() // Cerrar pool de conexiones al salir.
	shared.InitSessionStore()
	shared.InitStripe()

	mux := http.NewServeMux()

	// 2. Registrar rutas API de cada feature.
	// Cada RegisterRoutes crea su propio service con la referencia a DB.
	auth.RegisterRoutes(mux, shared.DB)
	citas.RegisterRoutes(mux, shared.DB)
	mascotas.RegisterRoutes(mux, shared.DB)
	facturas.RegisterRoutes(mux, shared.DB)
	contacto.RegisterRoutes(mux) // No usa DB (stub).
	checkout.RegisterRoutes(mux, shared.DB)
	admin.RegisterRoutes(mux, shared.DB)
	veterinario.RegisterRoutes(mux, shared.DB)

	// 3. Servir archivos estáticos del frontend.
	// StripPrefix quita /frontend/ del path para que coincida con ../frontend/ en disco.
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend"))))

	// 4. Raíz → index.html (landing page que redirige a /frontend/features/home/).
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../index.html")
	})

	// 5. Envolver todo con middleware CORS.
	handler := shared.CORS(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
