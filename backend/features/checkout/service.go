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

func (s *CheckoutService) CrearSesion(ctx context.Context, idFactura, idCliente int64) (string, error) {
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

	unitAmount := int64(math.Round(total * 100))

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

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
		SuccessURL: stripe.String(baseURL + "/frontend/features/pago-felicidades/?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(baseURL + "/frontend/features/dashboard/"),
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("error al crear sesion de Stripe: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE FACTURAS SET STRIPE_SESSION_ID = ? WHERE ID_FACTURA = ?`,
		sess.ID, idFactura,
	)
	if err != nil {
		return "", err
	}

	return sess.URL, nil
}

func (s *CheckoutService) VerificarPago(ctx context.Context, sessionID string) error {
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
	if estado == "Pagada" {
		return nil
	}

	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return fmt.Errorf("error al verificar con Stripe: %w", err)
	}

	if sess.Status != stripe.CheckoutSessionStatusComplete {
		return errors.New("el pago aun no ha sido completado")
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE FACTURAS SET ESTADO = 'Pagada' WHERE ID_FACTURA = ?`,
		idFactura,
	)
	return err
}
