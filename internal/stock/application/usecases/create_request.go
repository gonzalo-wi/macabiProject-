package stockusecases

import (
	"context"
	"fmt"
	"time"

	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
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
	projectReader stockports.ProjectMemberReader
	emailReader   stockports.UserEmailReader
	mailer        stockports.StockMailer
	pushNotifier  stockports.UserPushNotifier
}

func NewCreateRequest(repo stockports.StockRepository, projectReader stockports.ProjectMemberReader, emailReader stockports.UserEmailReader, mailer stockports.StockMailer, pushNotifier stockports.UserPushNotifier) *CreateRequest {
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
		// Coordinators get their request auto-approved (atomic stock deduction).
		if err := uc.repo.ApproveRequest(ctx, req.ID); err != nil {
			return nil, err
		}
		req.Status = stockdomain.RequestStatusApproved
	} else {
		// Madrij and others: notify coordinators so they can approve/reject.
		coordinators, err := uc.projectReader.FindProjectCoordinators(ctx, input.ProjectID)
		if err == nil {
			msg := fmt.Sprintf("Nueva solicitud de reserva: %d unidad(es) de \"%s\"", input.Quantity, resource.Name)
			for _, coordinatorID := range coordinators {
				notif := &stockdomain.StockNotification{
					UserID:    coordinatorID,
					RequestID: req.ID,
					Message:   msg,
				}
				_ = uc.repo.SaveNotification(ctx, notif)
			}
			// Best-effort email + push to all coordinators.
			if emails, err := uc.emailReader.FindEmailsByIDs(ctx, coordinators); err == nil {
				addrs := make([]string, 0, len(emails))
				for _, e := range emails {
					addrs = append(addrs, e)
				}
				_ = uc.mailer.NotifyCoordinatorsNewRequest(ctx, addrs, resource.Name, input.Quantity)
			}
			for _, coordinatorID := range coordinators {
				uc.pushNotifier.Notify(ctx, coordinatorID,
					"Nueva solicitud de reserva",
					fmt.Sprintf("%d unidad(es) de \"%s\" esperan aprobación", input.Quantity, resource.Name),
					fmt.Sprintf("/app/mis-proyectos/%s/recursos", input.ProjectID),
				)
			}
		}
	}

	return req, nil
}
