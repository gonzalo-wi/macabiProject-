package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type ListNotifications struct {
	repo stockports.StockRepository
}

func NewListNotifications(repo stockports.StockRepository) *ListNotifications {
	return &ListNotifications{repo: repo}
}

func (uc *ListNotifications) Execute(ctx context.Context, userID string) ([]stockdomain.StockNotification, error) {
	return uc.repo.ListNotificationsByUser(ctx, userID)
}
