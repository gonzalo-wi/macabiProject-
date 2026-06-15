package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type MarkEventNotificationRead struct {
	repo eventports.EventNotificationRepository
}

func NewMarkEventNotificationRead(repo eventports.EventNotificationRepository) *MarkEventNotificationRead {
	return &MarkEventNotificationRead{repo: repo}
}

func (uc *MarkEventNotificationRead) Execute(ctx context.Context, notifID, userID string) error {
	return uc.repo.MarkNotificationRead(ctx, notifID, userID)
}

type MarkAllEventNotificationsRead struct {
	repo eventports.EventNotificationRepository
}

func NewMarkAllEventNotificationsRead(repo eventports.EventNotificationRepository) *MarkAllEventNotificationsRead {
	return &MarkAllEventNotificationsRead{repo: repo}
}

func (uc *MarkAllEventNotificationsRead) Execute(ctx context.Context, userID string) error {
	_, err := uc.repo.MarkAllNotificationsRead(ctx, userID)
	return err
}
