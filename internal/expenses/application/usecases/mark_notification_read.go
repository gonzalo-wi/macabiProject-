package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
)

type MarkExpenseNotificationRead struct {
	repo expensesports.ExpenseRepository
}

func NewMarkExpenseNotificationRead(repo expensesports.ExpenseRepository) *MarkExpenseNotificationRead {
	return &MarkExpenseNotificationRead{repo: repo}
}

func (uc *MarkExpenseNotificationRead) Execute(ctx context.Context, notifID, userID string) error {
	return uc.repo.MarkNotificationRead(ctx, notifID, userID)
}
