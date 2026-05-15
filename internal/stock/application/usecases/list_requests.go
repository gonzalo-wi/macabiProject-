package stockusecases

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type ListRequests struct {
	repo stockports.StockRepository
}

func NewListRequests(repo stockports.StockRepository) *ListRequests {
	return &ListRequests{repo: repo}
}

// Execute lists all requests (admin-facing). projectID is optional — pass "" for all.
func (uc *ListRequests) Execute(ctx context.Context, params pagination.Params, projectID string) (pagination.Result[stockdomain.RequestDetail], error) {
	return uc.repo.ListRequests(ctx, params, projectID)
}
