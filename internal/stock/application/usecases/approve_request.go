package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userdomain "macabi-back/internal/user/domain"
)

type ApproveRequestInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

type ApproveRequest struct {
	repo          stockports.StockRepository
	projectReader stockports.ProjectMemberReader
}

func NewApproveRequest(repo stockports.StockRepository, projectReader stockports.ProjectMemberReader) *ApproveRequest {
	return &ApproveRequest{repo: repo, projectReader: projectReader}
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

	// Check stock is still available (may have changed since request was created).
	resource, err := uc.repo.FindResourceByID(ctx, req.ResourceID)
	if err != nil {
		return err
	}
	if resource.AvailableStock < req.Quantity {
		return stockdomain.ErrInsufficientStock
	}

	// Atomic: status → RESERVADO + available_stock -= quantity.
	return uc.repo.ApproveRequest(ctx, input.RequestID)
}
