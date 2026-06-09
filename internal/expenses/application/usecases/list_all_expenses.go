package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type ListAllExpenses struct {
	repo expensesports.ExpenseQueryRepository
}

func NewListAllExpenses(repo expensesports.ExpenseQueryRepository) *ListAllExpenses {
	return &ListAllExpenses{repo: repo}
}

func (uc *ListAllExpenses) Execute(ctx context.Context, filter expensesports.ExpenseListFilter, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error) {
	return uc.repo.ListAll(ctx, filter, params)
}
