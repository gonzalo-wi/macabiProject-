package stockusecases

import (
	"context"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userdomain "macabi-back/internal/user/domain"
)

type ReturnRequestInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

type ReturnRequest struct {
	repo          stockports.StockRepository
	projectReader stockports.ProjectMemberReader
}

func NewReturnRequest(repo stockports.StockRepository, projectReader stockports.ProjectMemberReader) *ReturnRequest {
	return &ReturnRequest{repo: repo, projectReader: projectReader}
}

func (uc *ReturnRequest) Execute(ctx context.Context, input ReturnRequestInput) error {
	req, err := uc.repo.FindRequestByID(ctx, input.RequestID)
	if err != nil {
		return err
	}
	if req.Status != stockdomain.RequestStatusDelivered {
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
	if resource.Type != stockdomain.ResourceTypeReturnable {
		return stockdomain.ErrNotReturnable
	}

	return uc.repo.ReturnRequest(ctx, input.RequestID)
}
