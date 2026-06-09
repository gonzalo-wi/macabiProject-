package expensesusecases

import (
	"context"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

// CoordinatorNotifier notifies project coordinators of new pending expenses.
type CoordinatorNotifier interface {
	NotifyCoordinatorsNewPendingExpense(ctx context.Context, exp *expensesdomain.Expense)
}

type CreateExpense struct {
	repo     expensesports.ExpenseRepository
	projects projectports.ProjectMembership
	notifier CoordinatorNotifier
}

func NewCreateExpense(
	repo expensesports.ExpenseRepository,
	projects projectports.ProjectMembership,
	notifier CoordinatorNotifier,
) *CreateExpense {
	return &CreateExpense{repo: repo, projects: projects, notifier: notifier}
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
		uc.notifier.NotifyCoordinatorsNewPendingExpense(ctx, exp)
	}

	return exp, nil
}
