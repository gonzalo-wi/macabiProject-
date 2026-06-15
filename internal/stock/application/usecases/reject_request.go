package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type RejectRequestInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

type RejectRequest struct {
	repo          stockports.StockRepository
	projectReader projectports.ProjectMembership
	notifier      stockports.UserNotifier
}

func NewRejectRequest(
	repo stockports.StockRepository,
	projectReader projectports.ProjectMembership,
	notifier stockports.UserNotifier,
) *RejectRequest {
	return &RejectRequest{repo: repo, projectReader: projectReader, notifier: notifier}
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

	if err := uc.repo.UpdateRequestStatus(ctx, input.RequestID, stockdomain.RequestStatusRejected); err != nil {
		return err
	}

	if resource, err := uc.repo.FindResourceByID(ctx, req.ResourceID); err == nil {
		uc.notifier.NotifyRequesterRejected(ctx, input.RequestID, req.RequestedByID, resource.Name, req.Quantity)
	}
	return nil
}
