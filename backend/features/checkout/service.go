/*
checkout/service.go — Integración con Stripe Checkout para pagos en línea.

FLUJO COMPLETO DE PAGO:
 1. Cliente ve facturas en dashboard → factura con Estado="Pendiente" muestra botón "Pagar".
 2. Frontend hace POST /api/checkout/crear-sesion con {id_factura}.
 3. CrearSesion():
    a. Valida que la factura existe Y pertenece al cliente (JOIN hasta MASCOTAS.ID_CLIENTE).
    b. Verifica que Estado=="Pendiente" (no permite doble pago).
    c. Crea una Stripe Checkout Session con el total en CRC (colón costarricense).
    d. Guarda el STRIPE_SESSION_ID en la factura para poder verificar después.
    e. Retorna la URL de Stripe Checkout al frontend.
 4. Frontend redirige al usuario a Stripe (window.location.href = url).
 5. Usuario paga en Stripe → Stripe redirige a /pago-felicidades/?session_id=X.
 6. Página de felicidades hace POST /api/checkout/verificar con {session_id}.
 7. VerificarPago():
    a. Busca la factura por STRIPE_SESSION_ID.
    b. Consulta a Stripe API si el pago fue completado (session.Status==complete).
    c. Si sí, actualiza Estado='Pagada' en FACTURAS.

SEGURIDAD:
- La factura se valida contra ID_CLIENTE de la sesión (no del body).
- No se puede pagar una factura ajena (JOIN falla → sql.ErrNoRows).
- No se puede pagar dos veces (estado != 'Pendiente' → error).

MONEDA: CRC (colón costarricense). Stripe requiere centavos → math.Round(total*100).
*/
package checkout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"

	stripe "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type CheckoutService struct {
	db *sql.DB
}

func NewCheckoutService(db *sql.DB) *CheckoutService {
	return &CheckoutService{db: db}
}

// CrearSesion crea una Stripe Checkout Session para la factura indicada.
// Paso 1: Validar propiedad (JOIN hasta MASCOTAS.ID_CLIENTE = idCliente).
// Paso 2: Verificar estado Pendiente.
// Paso 3: Crear sesión Stripe con precio en céntimos CRC.
// Paso 4: Guardar STRIPE_SESSION_ID en la factura.
// Paso 5: Retornar URL de checkout.
func (s *CheckoutService) CrearSesion(ctx context.Context, idFactura, idCliente int64) (string, error) {
	// Paso 1+2: Obtener total y estado, verificando propiedad del cliente.
	var total float64
	var estado string
	err := s.db.QueryRowContext(ctx,
		`SELECT f.TOTAL, f.ESTADO
		 FROM FACTURAS f
		 JOIN CITAS_SERVICIOS cs ON cs.FACTURAS_ID_FACTURA = f.ID_FACTURA
		 JOIN CITAS c ON cs.ID_CITA = c.ID_CITA
		 JOIN MASCOTAS m ON c.ID_MASCOTA = m.ID_MASCOTA
		 WHERE f.ID_FACTURA = ? AND m.ID_CLIENTE = ?
		 GROUP BY f.ID_FACTURA`,
		idFactura, idCliente,
	).Scan(&total, &estado)
	if err == sql.ErrNoRows {
		return "", errors.New("factura no encontrada o no autorizada")
	}
	if err != nil {
		return "", err
	}
	if estado != "Pendiente" {
		return "", errors.New("esta factura ya fue pagada")
	}

	// Paso 3: Stripe espera centavos → total * 100 redondeado.
	unitAmount := int64(math.Round(total * 100))

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Paso 3b: SuccessURL incluye {CHECKOUT_SESSION_ID} que Stripe reemplaza
	// automáticamente con el ID real de la sesión al redirigir.
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("crc"),
					UnitAmount: stripe.Int64(unitAmount),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(fmt.Sprintf("Factura #%d - Patitas al Rescate", idFactura)),
						Description: stripe.String("Pago de servicios veterinarios"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(baseURL + "/features/pago-felicidades/?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(baseURL + "/features/dashboard/"),
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("error al crear sesion de Stripe: %w", err)
	}

	// Paso 4: Guardar ID de sesión para poder verificar después del pago.
	_, err = s.db.ExecContext(ctx,
		`UPDATE FACTURAS SET STRIPE_SESSION_ID = ? WHERE ID_FACTURA = ?`,
		sess.ID, idFactura,
	)
	if err != nil {
		return "", err
	}

	// Paso 5: Retornar URL de Stripe Checkout para redirigir al cliente.
	return sess.URL, nil
}

// VerificarPago consulta a Stripe si el pago fue exitoso y actualiza la factura.
// Se llama desde la página /pago-felicidades/ que recibe session_id como query param.
func (s *CheckoutService) VerificarPago(ctx context.Context, sessionID string) error {
	// Buscar factura por su STRIPE_SESSION_ID almacenado.
	var idFactura int64
	var estado string
	err := s.db.QueryRowContext(ctx,
		`SELECT ID_FACTURA, ESTADO FROM FACTURAS WHERE STRIPE_SESSION_ID = ?`,
		sessionID,
	).Scan(&idFactura, &estado)
	if err == sql.ErrNoRows {
		return errors.New("sesion de pago no encontrada")
	}
	if err != nil {
		return err
	}
	// Idempotencia: si ya está pagada, no hacer nada.
	if estado == "Pagada" {
		return nil
	}

	// Consultar a Stripe API si la sesión fue completada.
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return fmt.Errorf("error al verificar con Stripe: %w", err)
	}

	if sess.Status != stripe.CheckoutSessionStatusComplete {
		return errors.New("el pago aun no ha sido completado")
	}

	// Marcar como pagada en la base de datos.
	_, err = s.db.ExecContext(ctx,
		`UPDATE FACTURAS SET ESTADO = 'Pagada' WHERE ID_FACTURA = ?`,
		idFactura,
	)
	return err
}
