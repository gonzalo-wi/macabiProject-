package eventmail

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

type BrevoEventMailer struct {
	apiKey      string
	from        string
	frontendURL string
	client      *http.Client
}

func NewBrevoEventMailer(apiKey, fromEmail, frontendURL string) *BrevoEventMailer {
	return &BrevoEventMailer{
		apiKey:      strings.TrimSpace(apiKey),
		from:        strings.TrimSpace(fromEmail),
		frontendURL: strings.TrimSuffix(strings.TrimSpace(frontendURL), "/"),
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

type NoOpEventMailer struct{}

func (NoOpEventMailer) SendEventReminder(context.Context, string, string, string, time.Time, string) error {
	return nil
}

type brevoPayload struct {
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

func (m *BrevoEventMailer) SendEventReminder(ctx context.Context, toEmail, userName, eventTitle string, deadline time.Time, eventID string) error {
	toEmail = strings.TrimSpace(strings.ToLower(toEmail))
	if toEmail == "" {
		return nil
	}

	deadlineStr := deadline.Format("02/01/2006 15:04")
	eventURL := fmt.Sprintf("%s/app/eventos/%s", m.frontendURL, eventID)

	body := fmt.Sprintf(`
<p style="margin:0 0 12px;font-size:15px;color:#374151;">Hola <strong>%s</strong>,</p>
<p style="margin:0 0 16px;font-size:15px;color:#374151;">
  Te recordamos que mañana vence el plazo para responder la jornada <strong>%s</strong>.
</p>
%s
%s`,
		userName,
		eventTitle,
		email.DetailsCard([]email.DetailRow{
			{Label: "Jornada", Value: eventTitle},
			{Label: "Cierre de respuestas", Value: deadlineStr},
		}),
		email.CTAButton(eventURL, "Responder ahora", "#7c3aed"),
	)

	html := email.Layout("#7c3aed", "Recordatorio de jornada", body)
	return m.send(ctx, toEmail, "Recordatorio: respondé la jornada antes de mañana — Macabi Madrijim", html)
}

func (m *BrevoEventMailer) send(ctx context.Context, toEmail, subject, htmlContent string) error {
	if m.apiKey == "" || m.from == "" {
		return nil
	}
	payload := brevoPayload{
		Sender:      brevoSender{Email: m.from, Name: email.SenderDisplayName},
		To:          []brevoRecipient{{Email: toEmail}},
		Subject:     subject,
		HTMLContent: htmlContent,
	}
	raw, err := json.Marshal(payload)
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
