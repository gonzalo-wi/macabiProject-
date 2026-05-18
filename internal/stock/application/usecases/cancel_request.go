package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userdomain "macabi-back/internal/user/domain"
)

type CancelRequestInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

type CancelRequest struct {
	repo stockports.StockRepository
}

func NewCancelRequest(repo stockports.StockRepository) *CancelRequest {
	return &CancelRequest{repo: repo}
}

func (uc *CancelRequest) Execute(ctx context.Context, input CancelRequestInput) error {
	req, err := uc.repo.FindRequestByID(ctx, input.RequestID)
	if err != nil {
		return err
	}

	if req.Status != stockdomain.RequestStatusPending {
		return stockdomain.ErrCannotCancel
	}

	if userdomain.Role(input.UserRole) != userdomain.RoleAdmin {
		if req.RequestedByID != input.UserID {
			return stockdomain.ErrNotOwner
		}
	}

	return uc.repo.UpdateRequestStatus(ctx, input.RequestID, stockdomain.RequestStatusCancelled)
}
