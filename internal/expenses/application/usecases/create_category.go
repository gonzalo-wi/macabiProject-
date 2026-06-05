package expensesusecases

import (
	"context"
	"strings"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
)

type CreateExpenseCategory struct {
	repo expensesports.CategoryRepository
}

func NewCreateExpenseCategory(repo expensesports.CategoryRepository) *CreateExpenseCategory {
	return &CreateExpenseCategory{repo: repo}
}

func (uc *CreateExpenseCategory) Execute(ctx context.Context, name string) (*expensesdomain.ExpenseCategory, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, expensesdomain.ErrMissingRequiredField
	}
	cat := &expensesdomain.ExpenseCategory{Name: name}
	if err := uc.repo.SaveCategory(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}
