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
	return m.send(ctx, toEmail, "Invitación a Macabi Madrijim",
		fmt.Sprintf(`<p>Te invitaron a crear tu cuenta en Macabi Madrijim.</p>
<p><a href="%s">Completá tu registro y elegí tu contraseña</a></p>
<p>Si no esperabas este correo, ignoralo.</p>`, acceptURL))
}

func (m *BrevoTransactionalMailer) send(ctx context.Context, toEmail, subject, html string) error {
	toEmail = strings.TrimSpace(strings.ToLower(toEmail))
	if toEmail == "" {
		return fmt.Errorf("empty recipient")
	}

	body := brevoSendEmailRequest{
		Sender:  brevoSender{Email: m.from, Name: brevoSenderDisplayName},
		To:      []brevoRecipient{{Email: toEmail}},
		Subject: subject,
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
