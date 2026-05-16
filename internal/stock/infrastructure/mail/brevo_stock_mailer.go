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
)

type BrevoStockMailer struct {
	apiKey string
	from   string
	client *http.Client
}

func NewBrevoStockMailer(apiKey, fromEmail string) *BrevoStockMailer {
	return &BrevoStockMailer{
		apiKey: strings.TrimSpace(apiKey),
		from:   strings.TrimSpace(fromEmail),
		client: &http.Client{Timeout: 30 * time.Second},
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

const senderDisplayName = "Macabi Madrijim"

func (m *BrevoStockMailer) NotifyCoordinatorsNewRequest(ctx context.Context, coordinatorEmails []string, resourceName string, quantity int) error {
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
	html := fmt.Sprintf(
		`<p>Hay una nueva solicitud de reserva pendiente de aprobación:</p>
<ul>
  <li><strong>Recurso:</strong> %s</li>
  <li><strong>Cantidad:</strong> %d unidad(es)</li>
</ul>
<p>Ingresá a la plataforma para aprobarla o rechazarla.</p>`,
		resourceName, quantity,
	)
	return m.send(ctx, to, "Nueva solicitud de reserva — Macabi Madrijim", html)
}

func (m *BrevoStockMailer) NotifyRequesterApproved(ctx context.Context, requesterEmail, resourceName string, quantity int) error {
	to := strings.TrimSpace(strings.ToLower(requesterEmail))
	if to == "" {
		return nil
	}
	html := fmt.Sprintf(
		`<p>Tu solicitud de reserva fue <strong>aprobada</strong>.</p>
<ul>
  <li><strong>Recurso:</strong> %s</li>
  <li><strong>Cantidad:</strong> %d unidad(es)</li>
</ul>
<p>Te avisaremos cuando el material esté listo para retirarse.</p>`,
		resourceName, quantity,
	)
	return m.send(ctx, []brevoRecipient{{Email: to}}, "Solicitud aprobada — Macabi Madrijim", html)
}

func (m *BrevoStockMailer) NotifyRequesterRejected(ctx context.Context, requesterEmail, resourceName string, quantity int) error {
	to := strings.TrimSpace(strings.ToLower(requesterEmail))
	if to == "" {
		return nil
	}
	html := fmt.Sprintf(
		`<p>Tu solicitud de reserva fue <strong>rechazada</strong>.</p>
<ul>
  <li><strong>Recurso:</strong> %s</li>
  <li><strong>Cantidad:</strong> %d unidad(es)</li>
</ul>
<p>Contactá a tu coordinador para más información.</p>`,
		resourceName, quantity,
	)
	return m.send(ctx, []brevoRecipient{{Email: to}}, "Solicitud rechazada — Macabi Madrijim", html)
}

func (m *BrevoStockMailer) send(ctx context.Context, to []brevoRecipient, subject, html string) error {
	body := brevoSendEmailRequest{
		Sender:      brevoSender{Email: m.from, Name: senderDisplayName},
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
