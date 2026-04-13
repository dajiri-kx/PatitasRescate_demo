package shared

import (
	"log"
	"os"

	stripe "github.com/stripe/stripe-go/v82"
)

func InitStripe() {
	key := os.Getenv("STRIPE_SECRET_KEY")
	if key == "" {
		log.Println("ADVERTENCIA: STRIPE_SECRET_KEY no configurada. El checkout no funcionará.")
		return
	}
	stripe.Key = key
	log.Println("Stripe inicializado correctamente.")
}
