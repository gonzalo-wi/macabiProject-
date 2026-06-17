package eventmail

import (
	"context"
	"fmt"
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
	return email.BrevoSend(ctx, m.client, m.apiKey, m.from, []email.BrevoRecipient{{Email: toEmail}}, "Recordatorio: respondé la jornada antes de mañana — Macabi Madrijim", html)
}
