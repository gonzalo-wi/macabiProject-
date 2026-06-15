package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type ListEventNotifications struct {
	repo eventports.EventNotificationRepository
}

func NewListEventNotifications(repo eventports.EventNotificationRepository) *ListEventNotifications {
	return &ListEventNotifications{repo: repo}
}

func (uc *ListEventNotifications) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[eventdomain.EventNotification], error) {
	return uc.repo.ListNotificationsByUser(ctx, userID, params)
}
