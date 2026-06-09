package expensesservices

import (
	"context"
	"fmt"
	"strings"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	projectports "macabi-back/internal/project/application/ports"
	userports "macabi-back/internal/user/application/ports"
)

// ExpenseNotificationService handles in-app and email notifications for expense events.
// It implements CoordinatorNotifier and SubmitterNotifier from the usecases package.
type ExpenseNotificationService struct {
	notifs       expensesports.ExpenseNotificationRepository
	coordinators projectports.ProjectCoordinatorReader
	emails       userports.UserEmailReader
	mailer       expensesports.ExpenseMailer
	repo         expensesports.ExpenseRepository
}

func NewExpenseNotificationService(
	notifs expensesports.ExpenseNotificationRepository,
	coordinators projectports.ProjectCoordinatorReader,
	emails userports.UserEmailReader,
	mailer expensesports.ExpenseMailer,
	repo expensesports.ExpenseRepository,
) *ExpenseNotificationService {
	return &ExpenseNotificationService{
		notifs:       notifs,
		coordinators: coordinators,
		emails:       emails,
		mailer:       mailer,
		repo:         repo,
	}
}

func (s *ExpenseNotificationService) NotifyCoordinatorsNewPendingExpense(ctx context.Context, exp *expensesdomain.Expense) {
	coordinators, err := s.coordinators.FindProjectCoordinators(ctx, exp.ProjectID)
	if err != nil || len(coordinators) == 0 {
		return
	}

	amountLabel := formatAmount(exp)
	msg := fmt.Sprintf("Nuevo gasto pendiente: %s — %s", amountLabel, exp.Description)

	for _, coordinatorID := range coordinators {
		_ = s.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
			UserID:    coordinatorID,
			ExpenseID: exp.ID,
			ProjectID: exp.ProjectID,
			Message:   msg,
		})
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
	_ = s.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
		UserID:    exp.SubmittedByUserID,
		ExpenseID: exp.ID,
		ProjectID: exp.ProjectID,
		Message:   fmt.Sprintf("Tu gasto fue aprobado: %s — %s", formatAmount(exp), exp.Description),
	})

	if email, err := s.emails.FindEmailByID(ctx, exp.SubmittedByUserID); err == nil {
		_ = s.mailer.NotifySubmitterApproved(ctx, email, formatAmount(exp), exp.Description, exp.ID)
	}
}

func (s *ExpenseNotificationService) NotifySubmitterRejected(ctx context.Context, exp *expensesdomain.Expense, reason string) {
	msg := fmt.Sprintf("Tu gasto fue rechazado: %s — %s", formatAmount(exp), exp.Description)
	if reason != "" {
		msg += fmt.Sprintf(" (motivo: %s)", reason)
	}
	_ = s.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
		UserID:    exp.SubmittedByUserID,
		ExpenseID: exp.ID,
		ProjectID: exp.ProjectID,
		Message:   msg,
	})

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
