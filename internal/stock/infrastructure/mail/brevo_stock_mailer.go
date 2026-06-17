package stockmail

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"macabi-back/internal/shared/email"
)

type BrevoStockMailer struct {
	apiKey      string
	from        string
	frontendURL string
	client      *http.Client
}

func NewBrevoStockMailer(apiKey, fromEmail, frontendURL string) *BrevoStockMailer {
	return &BrevoStockMailer{
		apiKey:      strings.TrimSpace(apiKey),
		from:        strings.TrimSpace(fromEmail),
		frontendURL: strings.TrimSuffix(strings.TrimSpace(frontendURL), "/"),
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

func detailsCard(resourceName string, quantity int) string {
	return email.DetailsCard([]email.DetailRow{
		{Label: "Recurso", Value: resourceName},
		{Label: "Cantidad", Value: fmt.Sprintf("%d unidad(es)", quantity)},
	})
}

func (m *BrevoStockMailer) NotifyCoordinatorsNewRequest(ctx context.Context, coordinatorEmails []string, resourceName string, quantity int, requestID string) error {
	if len(coordinatorEmails) == 0 {
		return nil
	}
	to := make([]email.BrevoRecipient, 0, len(coordinatorEmails))
	for _, e := range coordinatorEmails {
		if e = strings.TrimSpace(strings.ToLower(e)); e != "" {
			to = append(to, email.BrevoRecipient{Email: e})
		}
	}
	if len(to) == 0 {
		return nil
	}
	requestURL := fmt.Sprintf("%s/app/stock/requests/%s", m.frontendURL, requestID)
	body := fmt.Sprintf(`
<p style="margin:0 0 8px;font-size:15px;color:#374151;">Hay una nueva solicitud de reserva <strong>pendiente de aprobación</strong>.</p>
<p style="margin:0;font-size:14px;color:#6b7280;">Revisá los detalles y tomá una acción desde la plataforma.</p>
%s
%s`,
		detailsCard(resourceName, quantity),
		email.CTAButton(requestURL, "Ver solicitud", "#2563eb"),
	)
	html := email.Layout("#2563eb", "Nueva solicitud de reserva", body)
	return email.BrevoSend(ctx, m.client, m.apiKey, m.from, to, "Nueva solicitud de reserva — Macabi Madrijim", html)
}

func (m *BrevoStockMailer) NotifyRequesterApproved(ctx context.Context, requesterEmail, resourceName string, quantity int) error {
	to := strings.TrimSpace(strings.ToLower(requesterEmail))
	if to == "" {
		return nil
	}
	body := fmt.Sprintf(`
<p style="margin:0 0 8px;font-size:15px;color:#374151;">Tu solicitud de reserva fue <strong style="color:#16a34a;">aprobada</strong> ✓</p>
<p style="margin:0;font-size:14px;color:#6b7280;">Te avisaremos cuando el material esté listo para retirarse.</p>
%s`,
		detailsCard(resourceName, quantity),
	)
	html := email.Layout("#16a34a", "Solicitud aprobada", body)
	return email.BrevoSend(ctx, m.client, m.apiKey, m.from, []email.BrevoRecipient{{Email: to}}, "Solicitud aprobada — Macabi Madrijim", html)
}

func (m *BrevoStockMailer) NotifyRequesterRejected(ctx context.Context, requesterEmail, resourceName string, quantity int) error {
	to := strings.TrimSpace(strings.ToLower(requesterEmail))
	if to == "" {
		return nil
	}
	body := fmt.Sprintf(`
<p style="margin:0 0 8px;font-size:15px;color:#374151;">Tu solicitud de reserva fue <strong style="color:#dc2626;">rechazada</strong>.</p>
<p style="margin:0;font-size:14px;color:#6b7280;">Contactá a tu coordinador para más información.</p>
%s`,
		detailsCard(resourceName, quantity),
	)
	html := email.Layout("#dc2626", "Solicitud rechazada", body)
	return email.BrevoSend(ctx, m.client, m.apiKey, m.from, []email.BrevoRecipient{{Email: to}}, "Solicitud rechazada — Macabi Madrijim", html)
}
