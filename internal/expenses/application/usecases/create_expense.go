package expensesusecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	userdomain "macabi-back/internal/user/domain"
)

type CreateExpense struct {
	repo         expensesports.ExpenseRepository
	notifs       expensesports.ExpenseNotificationRepository
	projects     expensesports.ProjectMembership
	coordinators expensesports.ProjectCoordinatorReader
	emails       expensesports.UserEmailReader
	mailer       expensesports.ExpenseMailer
}

func NewCreateExpense(
	repo expensesports.ExpenseRepository,
	notifs expensesports.ExpenseNotificationRepository,
	projects expensesports.ProjectMembership,
	coordinators expensesports.ProjectCoordinatorReader,
	emails expensesports.UserEmailReader,
	mailer expensesports.ExpenseMailer,
) *CreateExpense {
	return &CreateExpense{
		repo:         repo,
		notifs:       notifs,
		projects:     projects,
		coordinators: coordinators,
		emails:       emails,
		mailer:       mailer,
	}
}

type CreateExpenseInput struct {
	ProjectID   string
	SubmittedBy string
	UserRole    string
	AmountStr   string
	Description string
	ExpenseDate time.Time
	Currency    string
	CategoryID  *string
}

func (uc *CreateExpense) Execute(ctx context.Context, input CreateExpenseInput) (*expensesdomain.Expense, error) {
	pid := strings.TrimSpace(input.ProjectID)
	if pid == "" {
		return nil, expensesdomain.ErrMissingRequiredField
	}

	isAdmin := userdomain.Role(input.UserRole) == userdomain.RoleAdmin
	if !isAdmin {
		ok, err := uc.projects.IsProjectMember(ctx, pid, input.SubmittedBy)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, expensesdomain.ErrForbidden
		}
	}

	amount, err := decimal.NewFromString(strings.TrimSpace(input.AmountStr))
	if err != nil {
		return nil, expensesdomain.ErrInvalidAmount
	}

	exp, err := expensesdomain.NewExpense(pid, input.SubmittedBy, amount, input.Description, input.ExpenseDate)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Currency) != "" {
		exp.Currency = strings.TrimSpace(input.Currency)
	}
	if input.CategoryID != nil && strings.TrimSpace(*input.CategoryID) != "" {
		id := strings.TrimSpace(*input.CategoryID)
		exp.CategoryID = &id
	}

	coord, err := uc.projects.IsProjectCoordinator(ctx, pid, input.SubmittedBy)
	if err != nil {
		return nil, err
	}
	if coord {
		now := time.Now().UTC()
		u := input.SubmittedBy
		exp.Status = expensesdomain.StatusApproved
		exp.ApprovedByUserID = &u
		exp.ApprovedAt = &now
	}

	if err := uc.repo.Save(ctx, exp); err != nil {
		return nil, err
	}

	if exp.Status == expensesdomain.StatusPending {
		uc.notifyCoordinatorsNewExpense(ctx, exp)
	}

	return exp, nil
}

func (uc *CreateExpense) notifyCoordinatorsNewExpense(ctx context.Context, exp *expensesdomain.Expense) {
	coordinators, err := uc.coordinators.FindProjectCoordinators(ctx, exp.ProjectID)
	if err != nil || len(coordinators) == 0 {
		return
	}

	amountLabel := formatExpenseAmount(exp)
	msg := fmt.Sprintf("Nuevo gasto pendiente: %s — %s", amountLabel, exp.Description)

	for _, coordinatorID := range coordinators {
		notif := &expensesdomain.ExpenseNotification{
			UserID:    coordinatorID,
			ExpenseID: exp.ID,
			ProjectID: exp.ProjectID,
			Message:   msg,
		}
		_ = uc.notifs.SaveNotification(ctx, notif)
	}

	projectName, _ := uc.repo.FindProjectName(ctx, exp.ProjectID)
	if strings.TrimSpace(projectName) == "" {
		projectName = "Proyecto"
	}

	if emails, err := uc.emails.FindEmailsByIDs(ctx, coordinators); err == nil {
		addrs := make([]string, 0, len(emails))
		for _, e := range emails {
			addrs = append(addrs, e)
		}
		_ = uc.mailer.NotifyCoordinatorsNewExpense(ctx, addrs, amountLabel, exp.Description, projectName, exp.ID)
	}
}
