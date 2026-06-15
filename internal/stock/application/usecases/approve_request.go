package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type ApproveRequestInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

type ApproveRequest struct {
	repo          stockports.StockRepository
	projectReader projectports.ProjectMembership
	notifier      stockports.UserNotifier
}

func NewApproveRequest(
	repo stockports.StockRepository,
	projectReader projectports.ProjectMembership,
	notifier stockports.UserNotifier,
) *ApproveRequest {
	return &ApproveRequest{repo: repo, projectReader: projectReader, notifier: notifier}
}

func (uc *ApproveRequest) Execute(ctx context.Context, input ApproveRequestInput) error {
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

	resource, err := uc.repo.FindResourceByID(ctx, req.ResourceID)
	if err != nil {
		return err
	}
	if resource.AvailableStock < req.Quantity {
		return stockdomain.ErrInsufficientStock
	}

	if err := uc.repo.ApproveRequest(ctx, input.RequestID); err != nil {
		return err
	}

	uc.notifier.NotifyRequesterApproved(ctx, input.RequestID, req.RequestedByID, resource.Name, req.Quantity)
	return nil
}
