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

const senderDisplayName = "Macabi Madrijim"

// emailLayout wraps content in a responsive HTML email shell.
// headerColor is a CSS hex color, e.g. "#2563eb".
func emailLayout(headerColor, title, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#f4f6f9;font-family:Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f9;padding:32px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%%;background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <!-- Header -->
        <tr>
          <td style="background:%s;padding:28px 32px;">
            <p style="margin:0;font-size:13px;color:rgba(255,255,255,0.8);letter-spacing:1px;text-transform:uppercase;">Macabi Madrijim</p>
            <h1 style="margin:6px 0 0;font-size:22px;color:#ffffff;font-weight:700;">%s</h1>
          </td>
        </tr>
        <!-- Body -->
        <tr>
          <td style="padding:32px;">
            %s
          </td>
        </tr>
        <!-- Footer -->
        <tr>
          <td style="background:#f4f6f9;padding:20px 32px;border-top:1px solid #e5e7eb;">
            <p style="margin:0;font-size:12px;color:#9ca3af;text-align:center;">Este mensaje fue generado automáticamente por Macabi Madrijim. Por favor no respondas este correo.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, headerColor, title, body)
}

// detailsCard renders a two-column info card for resource/quantity.
func detailsCard(resourceName string, quantity int) string {
	return fmt.Sprintf(`
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f8fafc;border:1px solid #e5e7eb;border-radius:6px;margin:20px 0;">
  <tr>
    <td style="padding:14px 20px;border-bottom:1px solid #e5e7eb;">
      <span style="font-size:12px;color:#6b7280;text-transform:uppercase;letter-spacing:0.5px;">Recurso</span><br>
      <span style="font-size:16px;color:#111827;font-weight:600;">%s</span>
    </td>
  </tr>
  <tr>
    <td style="padding:14px 20px;">
      <span style="font-size:12px;color:#6b7280;text-transform:uppercase;letter-spacing:0.5px;">Cantidad</span><br>
      <span style="font-size:16px;color:#111827;font-weight:600;">%d unidad(es)</span>
    </td>
  </tr>
</table>`, resourceName, quantity)
}

// ctaButton renders a full-width CTA button.
func ctaButton(url, label, color string) string {
	return fmt.Sprintf(`
<table width="100%%" cellpadding="0" cellspacing="0" style="margin-top:24px;">
  <tr>
    <td align="center">
      <a href="%s" style="display:inline-block;background:%s;color:#ffffff;font-size:15px;font-weight:700;text-decoration:none;padding:14px 40px;border-radius:6px;">%s</a>
    </td>
  </tr>
  <tr>
    <td align="center" style="padding-top:12px;">
      <span style="font-size:12px;color:#9ca3af;">O copiá este enlace: <a href="%s" style="color:#6b7280;">%s</a></span>
    </td>
  </tr>
</table>`, url, color, label, url, url)
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
		ctaButton(requestURL, "Ver solicitud", "#2563eb"),
	)
	html := emailLayout("#2563eb", "Nueva solicitud de reserva", body)
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
	html := emailLayout("#16a34a", "Solicitud aprobada", body)
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
	html := emailLayout("#dc2626", "Solicitud rechazada", body)
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
