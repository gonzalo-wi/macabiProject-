package stockmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type brevoSendEmailRequest struct {
	Sender      brevoSender      `json:"sender"`
	To          []brevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

type brevoSender struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type brevoRecipient struct {
	Email string `json:"email"`
}

// detailsCard renders the resource/quantity info card via the shared helper.
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
	to := make([]brevoRecipient, 0, len(coordinatorEmails))
	for _, e := range coordinatorEmails {
		if e = strings.TrimSpace(strings.ToLower(e)); e != "" {
			to = append(to, brevoRecipient{Email: e})
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
	return m.send(ctx, to, "Nueva solicitud de reserva — Macabi Madrijim", html)
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
	return m.send(ctx, []brevoRecipient{{Email: to}}, "Solicitud aprobada — Macabi Madrijim", html)
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
	return m.send(ctx, []brevoRecipient{{Email: to}}, "Solicitud rechazada — Macabi Madrijim", html)
}

func (m *BrevoStockMailer) send(ctx context.Context, to []brevoRecipient, subject, html string) error {
	body := brevoSendEmailRequest{
		Sender:      brevoSender{Email: m.from, Name: email.SenderDisplayName},
		To:          to,
		Subject:     subject,
		HTMLContent: html,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal brevo payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("brevo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", m.apiKey)

	res, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("brevo http: %w", err)
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("brevo: status %d: %s", res.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
