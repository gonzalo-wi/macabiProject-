package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
)

type UnreadExpenseNotificationCount struct {
	repo expensesports.ExpenseRepository
}

func NewUnreadExpenseNotificationCount(repo expensesports.ExpenseRepository) *UnreadExpenseNotificationCount {
	return &UnreadExpenseNotificationCount{repo: repo}
}

func (uc *UnreadExpenseNotificationCount) Execute(ctx context.Context, userID string) (int64, error) {
	return uc.repo.UnreadNotificationCount(ctx, userID)
}
