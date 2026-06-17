package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	"macabi-back/internal/shared/pagination"
)

type ListNotifications struct {
	repo stockports.NotificationRepository
}

func NewListNotifications(repo stockports.NotificationRepository) *ListNotifications {
	return &ListNotifications{repo: repo}
}

func (uc *ListNotifications) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.StockNotification], error) {
	return uc.repo.ListNotificationsByUser(ctx, userID, params)
}
