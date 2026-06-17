package stockusecases

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type ListMyRequests struct {
	repo stockports.RequestRepository
}

func NewListMyRequests(repo stockports.RequestRepository) *ListMyRequests {
	return &ListMyRequests{repo: repo}
}

func (uc *ListMyRequests) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.RequestDetail], error) {
	return uc.repo.ListMyRequests(ctx, userID, params)
}
