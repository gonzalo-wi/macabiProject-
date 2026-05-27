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
		To:      []brevoRecipient{{Email: toEmail}},
		Subject: "Restablecer contraseña — Macabi Madrijim",
		HTMLContent: emailLayout(
			"#d97706",
			"Restablecer contraseña",
			"Seguí los pasos para elegir una nueva contraseña",
			fmt.Sprintf(`
<p style="margin:0 0 12px;font-size:15px;color:#374151;">Recibimos una solicitud para restablecer la contraseña de tu cuenta en <strong>Macabi Madrijim</strong>.</p>
<p style="margin:0 0 4px;font-size:14px;color:#6b7280;">Hacé clic en el botón para elegir una nueva contraseña. El enlace es de un solo uso y tiene un tiempo de validez limitado.</p>
%s
<p style="margin:24px 0 0;font-size:13px;color:#9ca3af;">Si no solicitaste este cambio, podés ignorar este mensaje. Tu contraseña actual no se modificará.</p>`,
				ctaButton(resetURL, "Restablecer contraseña", "#d97706"),
			),
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
