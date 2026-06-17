package stockusecases

import (
	"context"
	"fmt"
	"time"

	projectports "macabi-back/internal/project/application/ports"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userports "macabi-back/internal/user/application/ports"
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
	emailReader   userports.UserEmailReader
	mailer        stockports.StockMailer
}

func NewCreateRequest(repo stockports.StockRepository, projectReader projectports.ProjectMemberReader, emailReader userports.UserEmailReader, mailer stockports.StockMailer) *CreateRequest {
	return &CreateRequest{repo: repo, projectReader: projectReader, emailReader: emailReader, mailer: mailer}
}

func (uc *CreateRequest) Execute(ctx context.Context, input CreateRequestInput) (*stockdomain.ResourceRequest, error) {
	resource, err := uc.repo.FindResourceByID(ctx, input.ResourceID)
	if err != nil {
		return nil, err
	}
	if resource.AvailableStock < input.Quantity {
		return nil, stockdomain.ErrInsufficientStock
	}
	if err := uc.authorizeRequest(ctx, input); err != nil {
		return nil, err
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
	if err := uc.approveOrNotify(ctx, req, input.ProjectID, input.RequestedByID, resource.Name, input.Quantity); err != nil {
		return nil, err
	}
	return req, nil
}

func (uc *CreateRequest) authorizeRequest(ctx context.Context, input CreateRequestInput) error {
	if userdomain.Role(input.UserRole) == userdomain.RoleAdmin {
		return nil
	}
	ok, err := uc.projectReader.IsProjectMember(ctx, input.ProjectID, input.RequestedByID)
	if err != nil {
		return err
	}
	if !ok {
		return stockdomain.ErrForbidden
	}
	return nil
}

func (uc *CreateRequest) approveOrNotify(ctx context.Context, req *stockdomain.ResourceRequest, projectID, requestedByID, resourceName string, quantity int) error {
	isCoordinator, err := uc.projectReader.IsProjectCoordinator(ctx, projectID, requestedByID)
	if err != nil {
		return err
	}
	if isCoordinator {
		if err := uc.repo.ApproveRequest(ctx, req.ID); err != nil {
			return err
		}
		req.Status = stockdomain.RequestStatusApproved
		return nil
	}
	uc.notifyCoordinators(ctx, req.ID, projectID, resourceName, quantity)
	return nil
}

func (uc *CreateRequest) notifyCoordinators(ctx context.Context, requestID, projectID, resourceName string, quantity int) {
	coordinators, err := uc.projectReader.FindProjectCoordinators(ctx, projectID)
	if err != nil {
		return
	}
	msg := fmt.Sprintf("Nueva solicitud de reserva: %d unidad(es) de \"%s\"", quantity, resourceName)
	for _, coordinatorID := range coordinators {
		_ = uc.repo.SaveNotification(ctx, &stockdomain.StockNotification{
			UserID:    coordinatorID,
			RequestID: requestID,
			Message:   msg,
		})
	}
	if emails, err := uc.emailReader.FindEmailsByIDs(ctx, coordinators); err == nil {
		_ = uc.mailer.NotifyCoordinatorsNewRequest(ctx, emails, resourceName, quantity, requestID)
	}
}
