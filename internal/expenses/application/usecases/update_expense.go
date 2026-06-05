package expensesusecases

import (
	"context"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	userdomain "macabi-back/internal/user/domain"
)

type UpdateExpense struct {
	repo     expensesports.ExpenseRepository
	projects expensesports.ProjectMembership
}

func NewUpdateExpense(repo expensesports.ExpenseRepository, projects expensesports.ProjectMembership) *UpdateExpense {
	return &UpdateExpense{repo: repo, projects: projects}
}

type UpdateExpenseInput struct {
	ID          string
	ActorID     string
	UserRole    string
	AmountStr   *string
	Description *string
	ExpenseDate *time.Time
	ReceiptPath *string
	Currency    *string
	CategoryID  *string // nil = no change, "" = remove, "uuid" = set
}

func (uc *UpdateExpense) Execute(ctx context.Context, input UpdateExpenseInput) (*expensesdomain.Expense, error) {
	exp, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	isAdmin := userdomain.Role(input.UserRole) == userdomain.RoleAdmin
	coord, err := uc.projects.IsProjectCoordinator(ctx, exp.ProjectID, input.ActorID)
	if err != nil {
		return nil, err
	}

	owner := exp.SubmittedByUserID == input.ActorID

	canEdit := false
	if isAdmin || coord {
		canEdit = true
	}
	if owner && exp.Status == expensesdomain.StatusPending {
		canEdit = true
	}
	if !canEdit {
		return nil, expensesdomain.ErrForbidden
	}

	restrictFinancial := !(isAdmin || coord)
	if owner && exp.Status == expensesdomain.StatusPending {
		restrictFinancial = false
	}
	if restrictFinancial && (input.Currency != nil || input.AmountStr != nil) {
		return nil, expensesdomain.ErrForbidden
	}

	if owner && !(isAdmin || coord) && input.ReceiptPath != nil && exp.Status != expensesdomain.StatusPending {
		return nil, expensesdomain.ErrForbidden
	}

	if input.AmountStr != nil {
		amt, err := decimal.NewFromString(strings.TrimSpace(*input.AmountStr))
		if err != nil {
			return nil, expensesdomain.ErrInvalidAmount
		}
		if amt.LessThanOrEqual(decimal.Zero) {
			return nil, expensesdomain.ErrInvalidAmount
		}
		exp.Amount = amt
	}
	if input.Description != nil {
		d := strings.TrimSpace(*input.Description)
		if d == "" {
			return nil, expensesdomain.ErrMissingRequiredField
		}
		exp.Description = d
	}
	if input.ExpenseDate != nil {
		exp.ExpenseDate = input.ExpenseDate.UTC().Truncate(24 * time.Hour)
	}
	if input.Currency != nil && strings.TrimSpace(*input.Currency) != "" {
		exp.Currency = strings.TrimSpace(*input.Currency)
	}
	if input.ReceiptPath != nil && *input.ReceiptPath != "" {
		p := strings.TrimSpace(*input.ReceiptPath)
		if !expensesdomain.ValidateReceiptObjectKey(exp.ProjectID, exp.ID, p) {
			return nil, expensesdomain.ErrInvalidReceiptPath
		}
		exp.ReceiptStoragePath = &p
	}
	if input.CategoryID != nil {
		if strings.TrimSpace(*input.CategoryID) == "" {
			exp.CategoryID = nil
		} else {
			id := strings.TrimSpace(*input.CategoryID)
			exp.CategoryID = &id
		}
	}

	if err := uc.repo.Update(ctx, exp); err != nil {
		return nil, err
	}
	return uc.repo.FindByID(ctx, input.ID)
}
