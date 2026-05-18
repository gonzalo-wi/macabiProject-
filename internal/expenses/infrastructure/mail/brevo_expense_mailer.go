package expensesmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

type BrevoExpenseMailer struct {
	apiKey string
	from   string
	client *http.Client
}

func NewBrevoExpenseMailer(apiKey, fromEmail string) *BrevoExpenseMailer {
	return &BrevoExpenseMailer{
		apiKey: strings.TrimSpace(apiKey),
		from:   strings.TrimSpace(fromEmail),
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type noopExpenseMailer struct{}

func NewNoOpExpenseMailer() *noopExpenseMailer {
	return &noopExpenseMailer{}
}

func (noopExpenseMailer) NotifyCoordinatorsNewExpense(context.Context, []string, string, string, string) error {
	return nil
}

func (noopExpenseMailer) NotifySubmitterApproved(context.Context, string, string, string) error {
	return nil
}

func (noopExpenseMailer) NotifySubmitterRejected(context.Context, string, string, string, string) error {
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

const senderDisplayName = "Macabi Madrijim"

func (m *BrevoExpenseMailer) NotifyCoordinatorsNewExpense(ctx context.Context, coordinatorEmails []string, amount, description, projectName string) error {
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
	htmlBody := fmt.Sprintf(
		`<p>Hay un nuevo gasto pendiente de aprobación en el proyecto <strong>%s</strong>:</p>
<ul>
  <li><strong>Monto:</strong> %s</li>
  <li><strong>Descripción:</strong> %s</li>
</ul>
<p>Ingresá a la plataforma para aprobarlo o rechazarlo.</p>`,
		html.EscapeString(projectName),
		html.EscapeString(amount),
		html.EscapeString(description),
	)
	return m.send(ctx, to, "Nuevo gasto pendiente — Macabi Madrijim", htmlBody)
}

func (m *BrevoExpenseMailer) NotifySubmitterApproved(ctx context.Context, submitterEmail, amount, description string) error {
	to := strings.TrimSpace(strings.ToLower(submitterEmail))
	if to == "" {
		return nil
	}
	htmlBody := fmt.Sprintf(
		`<p>Tu gasto fue <strong>aprobado</strong>.</p>
<ul>
  <li><strong>Monto:</strong> %s</li>
  <li><strong>Descripción:</strong> %s</li>
</ul>`,
		html.EscapeString(amount),
		html.EscapeString(description),
	)
	return m.send(ctx, []brevoRecipient{{Email: to}}, "Gasto aprobado — Macabi Madrijim", htmlBody)
}

func (m *BrevoExpenseMailer) NotifySubmitterRejected(ctx context.Context, submitterEmail, amount, description, reason string) error {
	to := strings.TrimSpace(strings.ToLower(submitterEmail))
	if to == "" {
		return nil
	}
	reasonBlock := ""
	if strings.TrimSpace(reason) != "" {
		reasonBlock = fmt.Sprintf(`<li><strong>Motivo:</strong> %s</li>`, html.EscapeString(strings.TrimSpace(reason)))
	}
	htmlBody := fmt.Sprintf(
		`<p>Tu gasto fue <strong>rechazado</strong>.</p>
<ul>
  <li><strong>Monto:</strong> %s</li>
  <li><strong>Descripción:</strong> %s</li>
  %s
</ul>
<p>Contactá a tu coordinador para más información.</p>`,
		html.EscapeString(amount),
		html.EscapeString(description),
		reasonBlock,
	)
	return m.send(ctx, []brevoRecipient{{Email: to}}, "Gasto rechazado — Macabi Madrijim", htmlBody)
}

func (m *BrevoExpenseMailer) send(ctx context.Context, to []brevoRecipient, subject, htmlContent string) error {
	if m.apiKey == "" || m.from == "" {
		return nil
	}
	body := brevoSendEmailRequest{
		Sender:      brevoSender{Email: m.from, Name: senderDisplayName},
		To:          to,
		Subject:     subject,
		HTMLContent: htmlContent,
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
