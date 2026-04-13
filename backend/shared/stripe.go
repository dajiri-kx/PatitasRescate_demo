/*
stripe.go — Inicialización de Stripe para pagos.

FLUJO DE DATOS:
 1. main.go llama a InitStripe() al arrancar.
 2. Se lee STRIPE_SECRET_KEY del entorno y se asigna a stripe.Key (global del SDK).
 3. Después, features/checkout/service.go usa el SDK de Stripe para crear
    sesiones de checkout y verificar pagos.
 4. Si la key no está configurada, el servidor arranca pero el módulo de pagos
    no funcionará (retornará errores de Stripe al intentar crear sesiones).

La API key de Stripe es SECRETA (sk_...) — nunca se expone al frontend.
El frontend solo recibe la URL de checkout a la que redirigir al usuario.
*/
package shared

import (
	"log"
	"os"

	stripe "github.com/stripe/stripe-go/v82"
)

// InitStripe configura la clave secreta de Stripe globalmente.
// Es opcional — si no está configurada, solo se muestra una advertencia.
func InitStripe() {
	key := os.Getenv("STRIPE_SECRET_KEY")
	if key == "" {
		log.Println("ADVERTENCIA: STRIPE_SECRET_KEY no configurada. El checkout no funcionará.")
		return
	}
	stripe.Key = key
	log.Println("Stripe inicializado correctamente.")
}
