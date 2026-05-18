package stockusecases

import (
	"context"
	"strings"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type UpdateResourceInput struct {
	ID         string
	Name       string
	Type       string
	TotalStock int
}

type UpdateResource struct {
	repo stockports.StockRepository
}

func NewUpdateResource(repo stockports.StockRepository) *UpdateResource {
	return &UpdateResource{repo: repo}
}

func (uc *UpdateResource) Execute(ctx context.Context, input UpdateResourceInput) (*stockdomain.Resource, error) {
	r, err := uc.repo.FindResourceByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, stockdomain.ErrEmptyResourceName
	}
	rtype := stockdomain.ResourceType(input.Type)
	if rtype != stockdomain.ResourceTypeReturnable && rtype != stockdomain.ResourceTypeConsumable {
		return nil, stockdomain.ErrInvalidResourceType
	}
	if input.TotalStock < 0 {
		return nil, stockdomain.ErrInvalidStock
	}

	diff := input.TotalStock - r.TotalStock
	r.Name = name
	r.Type = rtype
	r.TotalStock = input.TotalStock
	r.AvailableStock = r.AvailableStock + diff
	if r.AvailableStock < 0 {
		r.AvailableStock = 0
	}

	if err := uc.repo.UpdateResource(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}
