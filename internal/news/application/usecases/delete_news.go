package newsusecases

import (
	"context"

	newsports "macabi-back/internal/news/application/ports"
)

type DeleteNews struct {
	repo   newsports.NewsRepository
	notifs newsports.NewsNotificationRepository
	store  newsports.ImageStorage
}

func NewDeleteNews(repo newsports.NewsRepository, notifs newsports.NewsNotificationRepository, store newsports.ImageStorage) *DeleteNews {
	return &DeleteNews{repo: repo, notifs: notifs, store: store}
}

func (uc *DeleteNews) Execute(ctx context.Context, id string) error {
	n, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if uc.store != nil && n.ImageStoragePath != nil && *n.ImageStoragePath != "" {
		_ = uc.store.DeleteObject(ctx, *n.ImageStoragePath)
	}
	_ = uc.notifs.DeleteNotificationsByNewsID(ctx, id)
	return uc.repo.DeleteByID(ctx, id)
}
