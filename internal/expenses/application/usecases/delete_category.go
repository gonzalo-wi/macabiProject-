package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
)

type DeleteExpenseCategory struct {
	repo expensesports.CategoryRepository
}

func NewDeleteExpenseCategory(repo expensesports.CategoryRepository) *DeleteExpenseCategory {
	return &DeleteExpenseCategory{repo: repo}
}

func (uc *DeleteExpenseCategory) Execute(ctx context.Context, id string) error {
	return uc.repo.DeleteCategory(ctx, id)
}
