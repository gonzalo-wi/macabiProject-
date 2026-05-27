package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type ListExpenseNotifications struct {
	repo expensesports.ExpenseRepository
}

func NewListExpenseNotifications(repo expensesports.ExpenseRepository) *ListExpenseNotifications {
	return &ListExpenseNotifications{repo: repo}
}

func (uc *ListExpenseNotifications) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseNotification], error) {
	return uc.repo.ListNotificationsByUser(ctx, userID, params)
}
