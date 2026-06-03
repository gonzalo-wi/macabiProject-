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

type MarkAllExpenseNotificationsRead struct {
	repo expensesports.ExpenseRepository
}

func NewMarkAllExpenseNotificationsRead(repo expensesports.ExpenseRepository) *MarkAllExpenseNotificationsRead {
	return &MarkAllExpenseNotificationsRead{repo: repo}
}

func (uc *MarkAllExpenseNotificationsRead) Execute(ctx context.Context, userID string) error {
	_, err := uc.repo.MarkAllNotificationsRead(ctx, userID)
	return err
}
