package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type ListMyExpenses struct {
	repo expensesports.ExpenseRepository
}

func NewListMyExpenses(repo expensesports.ExpenseRepository) *ListMyExpenses {
	return &ListMyExpenses{repo: repo}
}

func (uc *ListMyExpenses) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error) {
	return uc.repo.ListMine(ctx, userID, params)
}
