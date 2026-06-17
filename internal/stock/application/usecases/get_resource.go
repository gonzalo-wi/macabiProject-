package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type GetResource struct {
	repo stockports.ResourceRepository
}

func NewGetResource(repo stockports.ResourceRepository) *GetResource {
	return &GetResource{repo: repo}
}

func (uc *GetResource) Execute(ctx context.Context, id string) (*stockdomain.Resource, error) {
	return uc.repo.FindResourceByID(ctx, id)
}
