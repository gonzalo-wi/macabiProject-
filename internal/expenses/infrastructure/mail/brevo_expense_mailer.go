package expensesmail

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"macabi-back/internal/shared/email"
)

type BrevoExpenseMailer struct {
	apiKey      string
	from        string
	frontendURL string
	client      *http.Client
}

func NewBrevoExpenseMailer(apiKey, fromEmail, frontendURL string) *BrevoExpenseMailer {
	return &BrevoExpenseMailer{
		apiKey:      strings.TrimSpace(apiKey),
		from:        strings.TrimSpace(fromEmail),
		frontendURL: strings.TrimSuffix(strings.TrimSpace(frontendURL), "/"),
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

type noopExpenseMailer struct{}

func NewNoOpExpenseMailer() *noopExpenseMailer {
	return &noopExpenseMailer{}
}

func (noopExpenseMailer) NotifyCoordinatorsNewExpense(context.Context, []string, string, string, string, string) error {
	return nil
}

func (noopExpenseMailer) NotifySubmitterApproved(context.Context, string, string, string, string) error {
	return nil
}

func (noopExpenseMailer) NotifySubmitterRejected(context.Context, string, string, string, string, string) error {
	return nil
}

func (m *BrevoExpenseMailer) expenseURL(expenseID string) string {
	return fmt.Sprintf("%s/app/gastos/%s", m.frontendURL, expenseID)
}

func (m *BrevoExpenseMailer) NotifyCoordinatorsNewExpense(ctx context.Context, coordinatorEmails []string, amount, description, projectName, expenseID string) error {
	if len(coordinatorEmails) == 0 {
		return nil
	}
	to := make([]email.BrevoRecipient, 0, len(coordinatorEmails))
	for _, e := range coordinatorEmails {
		if e = strings.TrimSpace(strings.ToLower(e)); e != "" {
			to = append(to, email.BrevoRecipient{Email: e})
		}
	}
	if len(to) == 0 {
		return nil
	}
	body := fmt.Sprintf(`
<p style="margin:0 0 8px;font-size:15px;color:#374151;">Hay un nuevo gasto <strong>pendiente de aprobación</strong>.</p>
<p style="margin:0;font-size:14px;color:#6b7280;">Revisá los detalles y aprobalo o rechazalo desde la plataforma.</p>
%s
%s`,
		email.DetailsCard([]email.DetailRow{
			{Label: "Proyecto", Value: projectName},
			{Label: "Monto", Value: amount},
			{Label: "Descripción", Value: description},
		}),
		email.CTAButton(m.expenseURL(expenseID), "Ver gasto", "#2563eb"),
	)
	html := email.Layout("#2563eb", "Nuevo gasto pendiente", body)
	return email.BrevoSend(ctx, m.client, m.apiKey, m.from, to, "Nuevo gasto pendiente — Macabi Madrijim", html)
}

func (m *BrevoExpenseMailer) NotifySubmitterApproved(ctx context.Context, submitterEmail, amount, description, expenseID string) error {
	to := strings.TrimSpace(strings.ToLower(submitterEmail))
	if to == "" {
		return nil
	}
	body := fmt.Sprintf(`
<p style="margin:0 0 8px;font-size:15px;color:#374151;">Tu gasto fue <strong style="color:#16a34a;">aprobado</strong> ✓</p>
%s
%s`,
		email.DetailsCard([]email.DetailRow{
			{Label: "Monto", Value: amount},
			{Label: "Descripción", Value: description},
		}),
		email.CTAButton(m.expenseURL(expenseID), "Ver gasto", "#16a34a"),
	)
	html := email.Layout("#16a34a", "Gasto aprobado", body)
	return email.BrevoSend(ctx, m.client, m.apiKey, m.from, []email.BrevoRecipient{{Email: to}}, "Gasto aprobado — Macabi Madrijim", html)
}

func (m *BrevoExpenseMailer) NotifySubmitterRejected(ctx context.Context, submitterEmail, amount, description, reason, expenseID string) error {
	to := strings.TrimSpace(strings.ToLower(submitterEmail))
	if to == "" {
		return nil
	}
	rows := []email.DetailRow{
		{Label: "Monto", Value: amount},
		{Label: "Descripción", Value: description},
	}
	if r := strings.TrimSpace(reason); r != "" {
		rows = append(rows, email.DetailRow{Label: "Motivo", Value: r})
	}
	body := fmt.Sprintf(`
<p style="margin:0 0 8px;font-size:15px;color:#374151;">Tu gasto fue <strong style="color:#dc2626;">rechazado</strong>.</p>
<p style="margin:0;font-size:14px;color:#6b7280;">Contactá a tu coordinador para más información.</p>
%s
%s`,
		email.DetailsCard(rows),
		email.CTAButton(m.expenseURL(expenseID), "Ver gasto", "#dc2626"),
	)
	html := email.Layout("#dc2626", "Gasto rechazado", body)
	return email.BrevoSend(ctx, m.client, m.apiKey, m.from, []email.BrevoRecipient{{Email: to}}, "Gasto rechazado — Macabi Madrijim", html)
}
