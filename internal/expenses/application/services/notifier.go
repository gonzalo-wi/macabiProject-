package expensesservices

import (
	"context"
	"fmt"
	"strings"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	projectports "macabi-back/internal/project/application/ports"
	"macabi-back/internal/shared/notifications"
	userports "macabi-back/internal/user/application/ports"
)

// ExpenseNotificationService unifica campana in-app, Web Push y email para gastos.
type ExpenseNotificationService struct {
	notifs       expensesports.ExpenseNotificationRepository
	coordinators projectports.ProjectCoordinatorReader
	emails       userports.UserEmailReader
	mailer       expensesports.ExpenseMailer
	repo         expensesports.ExpenseRepository
	push         notifications.PushNotifier
}

func NewExpenseNotificationService(
	notifs expensesports.ExpenseNotificationRepository,
	coordinators projectports.ProjectCoordinatorReader,
	emails userports.UserEmailReader,
	mailer expensesports.ExpenseMailer,
	repo expensesports.ExpenseRepository,
	push notifications.PushNotifier,
) *ExpenseNotificationService {
	return &ExpenseNotificationService{
		notifs:       notifs,
		coordinators: coordinators,
		emails:       emails,
		mailer:       mailer,
		repo:         repo,
		push:         push,
	}
}

func (s *ExpenseNotificationService) NotifyCoordinatorsNewPendingExpense(ctx context.Context, exp *expensesdomain.Expense) {
	coordinators, err := s.coordinators.FindProjectCoordinators(ctx, exp.ProjectID)
	if err != nil || len(coordinators) == 0 {
		return
	}

	filtered := make([]string, 0, len(coordinators))
	for _, id := range coordinators {
		if id != exp.SubmittedByUserID {
			filtered = append(filtered, id)
		}
	}
	if len(filtered) == 0 {
		return
	}
	coordinators = filtered

	amountLabel := formatAmount(exp)
	msg := fmt.Sprintf("Nuevo gasto pendiente: %s — %s", amountLabel, exp.Description)
	title := "Nuevo gasto pendiente"
	body := fmt.Sprintf("%s — %s", amountLabel, exp.Description)
	url := notifications.ExpenseDetail(exp.ID)

	for _, coordinatorID := range coordinators {
		_ = s.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
			UserID:    coordinatorID,
			ExpenseID: exp.ID,
			ProjectID: exp.ProjectID,
			Message:   msg,
		})
		notifications.PushToUser(ctx, s.push, coordinatorID, title, body, url)
	}

	projectName, _ := s.repo.FindProjectName(ctx, exp.ProjectID)
	if strings.TrimSpace(projectName) == "" {
		projectName = "Proyecto"
	}

	if emails, err := s.emails.FindEmailsByIDs(ctx, coordinators); err == nil {
		addrs := make([]string, 0, len(emails))
		for _, e := range emails {
			addrs = append(addrs, e)
		}
		_ = s.mailer.NotifyCoordinatorsNewExpense(ctx, addrs, amountLabel, exp.Description, projectName, exp.ID)
	}
}

func (s *ExpenseNotificationService) NotifySubmitterApproved(ctx context.Context, exp *expensesdomain.Expense) {
	msg := fmt.Sprintf("Tu gasto fue aprobado: %s — %s", formatAmount(exp), exp.Description)
	title := "Gasto aprobado"
	body := exp.Description + " — tu gasto fue aprobado"
	url := notifications.ExpenseDetail(exp.ID)

	_ = s.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
		UserID:    exp.SubmittedByUserID,
		ExpenseID: exp.ID,
		ProjectID: exp.ProjectID,
		Message:   msg,
	})
	notifications.PushToUser(ctx, s.push, exp.SubmittedByUserID, title, body, url)

	if email, err := s.emails.FindEmailByID(ctx, exp.SubmittedByUserID); err == nil {
		_ = s.mailer.NotifySubmitterApproved(ctx, email, formatAmount(exp), exp.Description, exp.ID)
	}
}

func (s *ExpenseNotificationService) NotifySubmitterRejected(ctx context.Context, exp *expensesdomain.Expense, reason string) {
	msg := fmt.Sprintf("Tu gasto fue rechazado: %s — %s", formatAmount(exp), exp.Description)
	if reason != "" {
		msg += fmt.Sprintf(" (motivo: %s)", reason)
	}
	title := "Gasto rechazado"
	body := exp.Description + " — tu gasto fue rechazado"
	url := notifications.ExpenseDetail(exp.ID)

	_ = s.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
		UserID:    exp.SubmittedByUserID,
		ExpenseID: exp.ID,
		ProjectID: exp.ProjectID,
		Message:   msg,
	})
	notifications.PushToUser(ctx, s.push, exp.SubmittedByUserID, title, body, url)

	if email, err := s.emails.FindEmailByID(ctx, exp.SubmittedByUserID); err == nil {
		_ = s.mailer.NotifySubmitterRejected(ctx, email, formatAmount(exp), exp.Description, reason, exp.ID)
	}
}

func formatAmount(exp *expensesdomain.Expense) string {
	currency := exp.Currency
	if currency == "" {
		currency = expensesdomain.DefaultCurrency
	}
	return fmt.Sprintf("%s %s", exp.Amount.StringFixed(2), currency)
}
