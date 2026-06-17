package stockports

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type NotificationRepository interface {
	SaveNotification(ctx context.Context, n *stockdomain.StockNotification) error
	ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.StockNotification], error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error)
	UnreadCount(ctx context.Context, userID string) (int64, error)
}
