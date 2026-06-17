package stockusecases

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type ListResources struct {
	repo stockports.ResourceRepository
}

func NewListResources(repo stockports.ResourceRepository) *ListResources {
	return &ListResources{repo: repo}
}

func (uc *ListResources) Execute(ctx context.Context, params pagination.Params) (pagination.Result[stockdomain.Resource], error) {
	return uc.repo.FindAllResources(ctx, params)
}
