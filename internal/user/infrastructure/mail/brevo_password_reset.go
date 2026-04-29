package usermail

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

type BrevoPasswordResetMailer struct {
	apiKey string
	from   string
	client *http.Client
}

func NewBrevoPasswordResetMailer(apiKey, fromEmail string) *BrevoPasswordResetMailer {
	key := strings.TrimSpace(apiKey)
	from := strings.TrimSpace(fromEmail)
	return &BrevoPasswordResetMailer{
		apiKey: key,
		from:   from,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (m *BrevoPasswordResetMailer) SendResetLink(ctx context.Context, toEmail, resetURL string) error {
	toEmail = strings.TrimSpace(strings.ToLower(toEmail))
	if toEmail == "" {
		return fmt.Errorf("empty recipient")
	}
	if resetURL == "" {
		return fmt.Errorf("empty reset url")
	}

	body := brevoSendEmailRequest{
		Sender: brevoSender{
			Email: m.from,
			Name:  brevoSenderDisplayName,
		},
		To:          []brevoRecipient{{Email: toEmail}},
		Subject:     "Restablecer contraseña — Macabi Madrijim",
		HTMLContent: fmt.Sprintf(
			`<p>Recibimos una solicitud para restablecer tu contraseña en Macabi Madrijim.</p>
<p><a href="%s">Hacé clic aquí para elegir una nueva contraseña</a></p>
<p>Si no pediste este cambio, podés ignorar este mensaje. El enlace deja de ser válido tras usarse o al pasar el tiempo indicado.</p>`,
			resetURL,
		),
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
