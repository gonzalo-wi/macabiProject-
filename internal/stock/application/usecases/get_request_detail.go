package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type GetRequestDetail struct {
	repo stockports.StockRepository
}

func NewGetRequestDetail(repo stockports.StockRepository) *GetRequestDetail {
	return &GetRequestDetail{repo: repo}
}

func (uc *GetRequestDetail) Execute(ctx context.Context, id string) (*stockdomain.RequestDetail, error) {
	return uc.repo.FindRequestDetailByID(ctx, id)
}
