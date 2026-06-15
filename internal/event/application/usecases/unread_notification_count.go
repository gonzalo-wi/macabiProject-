package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type UnreadEventNotificationCount struct {
	repo eventports.EventNotificationRepository
}

func NewUnreadEventNotificationCount(repo eventports.EventNotificationRepository) *UnreadEventNotificationCount {
	return &UnreadEventNotificationCount{repo: repo}
}

func (uc *UnreadEventNotificationCount) Execute(ctx context.Context, userID string) (int64, error) {
	return uc.repo.UnreadNotificationCount(ctx, userID)
}
