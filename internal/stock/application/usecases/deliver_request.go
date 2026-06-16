package stockusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userdomain "macabi-back/internal/user/domain"
)

type DeliverRequestInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

type DeliverRequest struct {
	repo          stockports.StockRepository
	projectReader projectports.ProjectMembership
}

func NewDeliverRequest(repo stockports.StockRepository, projectReader projectports.ProjectMembership) *DeliverRequest {
	return &DeliverRequest{repo: repo, projectReader: projectReader}
}

func (uc *DeliverRequest) Execute(ctx context.Context, input DeliverRequestInput) error {
	req, err := uc.repo.FindRequestByID(ctx, input.RequestID)
	if err != nil {
		return err
	}
	if req.Status != stockdomain.RequestStatusApproved {
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
	return uc.repo.UpdateRequestStatus(ctx, input.RequestID, stockdomain.RequestStatusDelivered)
}
