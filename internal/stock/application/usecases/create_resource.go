package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type CreateResourceInput struct {
	Name       string
	Type       string
	TotalStock int
}

type CreateResource struct {
	repo stockports.StockRepository
}

func NewCreateResource(repo stockports.StockRepository) *CreateResource {
	return &CreateResource{repo: repo}
}

func (uc *CreateResource) Execute(ctx context.Context, input CreateResourceInput) (*stockdomain.Resource, error) {
	r, err := stockdomain.NewResource(input.Name, stockdomain.ResourceType(input.Type), input.TotalStock)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.SaveResource(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}
