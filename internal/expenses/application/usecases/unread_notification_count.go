package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
)

type UnreadExpenseNotificationCount struct {
	repo expensesports.ExpenseNotificationRepository
}

func NewUnreadExpenseNotificationCount(repo expensesports.ExpenseNotificationRepository) *UnreadExpenseNotificationCount {
	return &UnreadExpenseNotificationCount{repo: repo}
}

func (uc *UnreadExpenseNotificationCount) Execute(ctx context.Context, userID string) (int64, error) {
	return uc.repo.UnreadNotificationCount(ctx, userID)
}
