package expensesports

import (
	"context"

	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type ExpenseNotificationRepository interface {
	SaveNotification(ctx context.Context, n *expensesdomain.ExpenseNotification) error
	ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseNotification], error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error)
	UnreadNotificationCount(ctx context.Context, userID string) (int64, error)
	DeleteNotificationsByExpenseID(ctx context.Context, expenseID string) error
}
