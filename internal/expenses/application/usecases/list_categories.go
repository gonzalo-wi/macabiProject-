package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
)

type ListExpenseCategories struct {
	repo expensesports.CategoryRepository
}

func NewListExpenseCategories(repo expensesports.CategoryRepository) *ListExpenseCategories {
	return &ListExpenseCategories{repo: repo}
}

func (uc *ListExpenseCategories) Execute(ctx context.Context) ([]expensesdomain.ExpenseCategory, error) {
	return uc.repo.ListCategories(ctx)
}
