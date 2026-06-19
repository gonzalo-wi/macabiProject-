package newsmail

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

type BrevoNewsMailer struct {
	apiKey      string
	from        string
	frontendURL string
	client      *http.Client
}

func NewBrevoNewsMailer(apiKey, fromEmail, frontendURL string) *BrevoNewsMailer {
	return &BrevoNewsMailer{
		apiKey:      strings.TrimSpace(apiKey),
		from:        strings.TrimSpace(fromEmail),
		frontendURL: strings.TrimSuffix(strings.TrimSpace(frontendURL), "/"),
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

type noopNewsMailer struct{}

func NewNoOpNewsMailer() *noopNewsMailer { return &noopNewsMailer{} }

func (noopNewsMailer) NotifyMembersNewNews(context.Context, []string, string, string, string) error {
	return nil
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

func (m *BrevoNewsMailer) newsURL(newsID string) string {
	return fmt.Sprintf("%s/app/noticias/%s", m.frontendURL, newsID)
}

func (m *BrevoNewsMailer) NotifyMembersNewNews(ctx context.Context, recipientEmails []string, title, summary, newsID string) error {
	if m.apiKey == "" || m.from == "" || len(recipientEmails) == 0 {
		return nil
	}

	rows := []email.DetailRow{{Label: "Noticia", Value: title}}
	if s := strings.TrimSpace(summary); s != "" {
		rows = append(rows, email.DetailRow{Label: "Resumen", Value: s})
	}
	body := fmt.Sprintf(`
<p style="margin:0 0 8px;font-size:15px;color:#374151;">Se publicó una <strong>nueva noticia</strong> en Macabi Madrijim.</p>
%s
%s`,
		email.DetailsCard(rows),
		email.CTAButton(m.newsURL(newsID), "Ver noticia", "#7c3aed"),
	)
	html := email.Layout("#7c3aed", "Nueva noticia", body)
	subject := "Nueva noticia — Macabi Madrijim"

	// Envío individual: cada miembro recibe su propio correo sin ver al resto.
	var firstErr error
	for _, e := range recipientEmails {
		addr := strings.TrimSpace(strings.ToLower(e))
		if addr == "" {
			continue
		}
		if err := m.send(ctx, addr, subject, html); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (m *BrevoNewsMailer) send(ctx context.Context, to, subject, htmlContent string) error {
	reqBody := brevoSendEmailRequest{
		Sender:      brevoSender{Email: m.from, Name: email.SenderDisplayName},
		To:          []brevoRecipient{{Email: to}},
		Subject:     subject,
		HTMLContent: htmlContent,
	}
	raw, err := json.Marshal(reqBody)
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
