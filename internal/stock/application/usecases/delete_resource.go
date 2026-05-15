package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type DeleteResource struct {
	repo stockports.StockRepository
}

func NewDeleteResource(repo stockports.StockRepository) *DeleteResource {
	return &DeleteResource{repo: repo}
}

func (uc *DeleteResource) Execute(ctx context.Context, id string) error {
	if _, err := uc.repo.FindResourceByID(ctx, id); err != nil {
		return err
	}
	active, err := uc.repo.HasActiveRequests(ctx, id)
	if err != nil {
		return err
	}
	if active {
		return stockdomain.ErrResourceInUse
	}
	return uc.repo.DeleteResource(ctx, id)
}
