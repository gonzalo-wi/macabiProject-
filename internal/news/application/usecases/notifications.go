package newsusecases

import (
	"context"

	newsports "macabi-back/internal/news/application/ports"
	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"
)

type ListNewsNotifications struct {
	repo newsports.NewsNotificationRepository
}

func NewListNewsNotifications(repo newsports.NewsNotificationRepository) *ListNewsNotifications {
	return &ListNewsNotifications{repo: repo}
}

func (uc *ListNewsNotifications) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[newsdomain.NewsNotification], error) {
	return uc.repo.ListNotificationsByUser(ctx, userID, params)
}

type MarkNewsNotificationRead struct {
	repo newsports.NewsNotificationRepository
}

func NewMarkNewsNotificationRead(repo newsports.NewsNotificationRepository) *MarkNewsNotificationRead {
	return &MarkNewsNotificationRead{repo: repo}
}

func (uc *MarkNewsNotificationRead) Execute(ctx context.Context, notifID, userID string) error {
	return uc.repo.MarkNotificationRead(ctx, notifID, userID)
}

type MarkAllNewsNotificationsRead struct {
	repo newsports.NewsNotificationRepository
}

func NewMarkAllNewsNotificationsRead(repo newsports.NewsNotificationRepository) *MarkAllNewsNotificationsRead {
	return &MarkAllNewsNotificationsRead{repo: repo}
}

func (uc *MarkAllNewsNotificationsRead) Execute(ctx context.Context, userID string) error {
	_, err := uc.repo.MarkAllNotificationsRead(ctx, userID)
	return err
}

type UnreadNewsNotificationCount struct {
	repo newsports.NewsNotificationRepository
}

func NewUnreadNewsNotificationCount(repo newsports.NewsNotificationRepository) *UnreadNewsNotificationCount {
	return &UnreadNewsNotificationCount{repo: repo}
}

func (uc *UnreadNewsNotificationCount) Execute(ctx context.Context, userID string) (int64, error) {
	return uc.repo.UnreadNotificationCount(ctx, userID)
}
