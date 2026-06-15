package stockusecases

import (
	"context"
	"time"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type CreateRequestInput struct {
	ProjectID      string
	ResourceID     string
	RequestedByID  string
	UserRole       string
	Quantity       int
	WithdrawalDate time.Time
	ReturnDate     *time.Time
	Notes          string
}

type CreateRequest struct {
	repo          stockports.StockRepository
	projectReader projectports.ProjectMemberReader
	notifier      stockports.UserNotifier
}

func NewCreateRequest(
	repo stockports.StockRepository,
	projectReader projectports.ProjectMemberReader,
	notifier stockports.UserNotifier,
) *CreateRequest {
	return &CreateRequest{repo: repo, projectReader: projectReader, notifier: notifier}
}

func (uc *CreateRequest) Execute(ctx context.Context, input CreateRequestInput) (*stockdomain.ResourceRequest, error) {
	resource, err := uc.repo.FindResourceByID(ctx, input.ResourceID)
	if err != nil {
		return nil, err
	}

	if resource.AvailableStock < input.Quantity {
		return nil, stockdomain.ErrInsufficientStock
	}

	if userdomain.Role(input.UserRole) != userdomain.RoleAdmin {
		ok, err := uc.projectReader.IsProjectMember(ctx, input.ProjectID, input.RequestedByID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, stockdomain.ErrForbidden
		}
	}

	req, err := stockdomain.NewResourceRequest(
		input.ProjectID,
		input.ResourceID,
		input.RequestedByID,
		input.Quantity,
		input.WithdrawalDate,
		input.ReturnDate,
		input.Notes,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.SaveRequest(ctx, req); err != nil {
		return nil, err
	}

	isCoordinator, _ := uc.projectReader.IsProjectCoordinator(ctx, input.ProjectID, input.RequestedByID)

	if isCoordinator {
		if err := uc.repo.ApproveRequest(ctx, req.ID); err != nil {
			return nil, err
		}
		req.Status = stockdomain.RequestStatusApproved
	} else {
		uc.notifier.NotifyCoordinatorsNewRequest(ctx, req.ID, input.ProjectID, resource.Name, input.Quantity)
	}

	return req, nil
}
