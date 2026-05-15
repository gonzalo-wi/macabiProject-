package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userdomain "macabi-back/internal/user/domain"
)

type RejectRequestInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

type RejectRequest struct {
	repo          stockports.StockRepository
	projectReader stockports.ProjectMemberReader
}

func NewRejectRequest(repo stockports.StockRepository, projectReader stockports.ProjectMemberReader) *RejectRequest {
	return &RejectRequest{repo: repo, projectReader: projectReader}
}

func (uc *RejectRequest) Execute(ctx context.Context, input RejectRequestInput) error {
	req, err := uc.repo.FindRequestByID(ctx, input.RequestID)
	if err != nil {
		return err
	}
	if req.Status != stockdomain.RequestStatusPending {
		return stockdomain.ErrInvalidStatusTransition
	}

	if userdomain.Role(input.UserRole) != userdomain.RoleAdmin {
		ok, err := uc.projectReader.IsProjectCoordinator(ctx, req.ProjectID, input.UserID)
		if err != nil {
			return err
		}
		if !ok {
			return stockdomain.ErrForbidden
		}
	}

	return uc.repo.UpdateRequestStatus(ctx, input.RequestID, stockdomain.RequestStatusRejected)
}
