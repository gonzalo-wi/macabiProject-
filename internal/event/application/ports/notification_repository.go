package eventports

import (
	"context"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type EventNotificationRepository interface {
	SaveNotification(ctx context.Context, n *eventdomain.EventNotification) error
	ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[eventdomain.EventNotification], error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error)
	UnreadNotificationCount(ctx context.Context, userID string) (int64, error)
	HasReminderForEvent(ctx context.Context, userID, eventInstanceID string) (bool, error)
}
