package stockusecases

import (
	"context"
	"fmt"
	"time"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	projectports "macabi-back/internal/project/application/ports"
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
	pushNotifier  stockports.UserPushNotifier
}

func NewCreateRequest(repo stockports.StockRepository, projectReader projectports.ProjectMemberReader, emailReader userports.UserEmailReader, mailer stockports.StockMailer, pushNotifier stockports.UserPushNotifier) *CreateRequest {
	return &CreateRequest{repo: repo, projectReader: projectReader, emailReader: emailReader, mailer: mailer, pushNotifier: pushNotifier}
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
		uc.notifyCoordinators(ctx, req.ID, input.ProjectID, resource.Name, input.Quantity)
	}

	return req, nil
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
		addrs := make([]string, 0, len(emails))
		for _, e := range emails {
			addrs = append(addrs, e)
		}
		_ = uc.mailer.NotifyCoordinatorsNewRequest(ctx, addrs, resourceName, quantity, requestID)
	}

	for _, coordinatorID := range coordinators {
		uc.pushNotifier.Notify(ctx, coordinatorID,
			"Nueva solicitud de reserva",
			fmt.Sprintf("%d unidad(es) de \"%s\" esperan aprobación", quantity, resourceName),
			fmt.Sprintf("/app/mis-proyectos/%s/recursos", projectID),
		)
	}
}
