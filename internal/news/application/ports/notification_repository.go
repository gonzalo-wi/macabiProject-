package newsports

import (
	"context"

	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"
)

type NewsNotificationRepository interface {
	// SaveNotifications inserta en lote (una noticia => muchos destinatarios).
	SaveNotifications(ctx context.Context, ns []*newsdomain.NewsNotification) error
	ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[newsdomain.NewsNotification], error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error)
	UnreadNotificationCount(ctx context.Context, userID string) (int64, error)
	DeleteNotificationsByNewsID(ctx context.Context, newsID string) error
}
