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

type BrevoTransactionalMailer struct {
	apiKey string
	from   string
	client *http.Client
}

func NewBrevoTransactionalMailer(apiKey, fromEmail string) *BrevoTransactionalMailer {
	return &BrevoTransactionalMailer{
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

const brevoSenderDisplayName = "Macabi Madrijim"

func (m *BrevoTransactionalMailer) SendInvitationLink(ctx context.Context, toEmail, acceptURL string) error {
	body := fmt.Sprintf(`
<p style="margin:0 0 12px;font-size:15px;color:#374151;">Fuiste invitado/a a unirte a la plataforma de gestión de <strong>Macabi Madrijim</strong>.</p>
<p style="margin:0 0 4px;font-size:14px;color:#6b7280;">Hacé clic en el botón para completar tu registro y elegir tu contraseña. El enlace tiene un tiempo de validez limitado.</p>
%s
<p style="margin:24px 0 0;font-size:13px;color:#9ca3af;">Si no esperabas esta invitación, podés ignorar este mensaje sin problema.</p>`,
		ctaButton(acceptURL, "Aceptar invitación", "#4f46e5"),
	)
	html := emailLayout("#4f46e5", "Te invitaron a Macabi Madrijim", "Completá tu registro para comenzar", body)
	return m.send(ctx, toEmail, "Invitación a Macabi Madrijim", html)
}

func (m *BrevoTransactionalMailer) send(ctx context.Context, toEmail, subject, html string) error {
	toEmail = strings.TrimSpace(strings.ToLower(toEmail))
	if toEmail == "" {
		return fmt.Errorf("empty recipient")
	}

	body := brevoSendEmailRequest{
		Sender:      brevoSender{Email: m.from, Name: brevoSenderDisplayName},
		To:          []brevoRecipient{{Email: toEmail}},
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
