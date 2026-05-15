package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
)

type MarkNotificationRead struct {
	repo stockports.StockRepository
}

func NewMarkNotificationRead(repo stockports.StockRepository) *MarkNotificationRead {
	return &MarkNotificationRead{repo: repo}
}

func (uc *MarkNotificationRead) Execute(ctx context.Context, notifID, userID string) error {
	return uc.repo.MarkNotificationRead(ctx, notifID, userID)
}
