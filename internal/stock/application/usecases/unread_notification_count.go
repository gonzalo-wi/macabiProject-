package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
)

type UnreadNotificationCount struct {
	repo stockports.NotificationRepository
}

func NewUnreadNotificationCount(repo stockports.NotificationRepository) *UnreadNotificationCount {
	return &UnreadNotificationCount{repo: repo}
}

func (uc *UnreadNotificationCount) Execute(ctx context.Context, userID string) (int64, error) {
	return uc.repo.UnreadCount(ctx, userID)
}
