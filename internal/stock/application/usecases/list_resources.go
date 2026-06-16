package stockusecases

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type ListResources struct {
	repo stockports.StockRepository
}

func NewListResources(repo stockports.StockRepository) *ListResources {
	return &ListResources{repo: repo}
}

func (uc *ListResources) Execute(ctx context.Context, filter stockports.ResourceListFilter, params pagination.Params) (pagination.Result[stockdomain.Resource], error) {
	return uc.repo.FindAllResources(ctx, filter, params)
}
