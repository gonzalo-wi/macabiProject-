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
	emailReader   stockports.UserEmailReader
	mailer        stockports.StockMailer
	pushNotifier  stockports.UserPushNotifier
}

func NewApproveRequest(repo stockports.StockRepository, projectReader stockports.ProjectMemberReader, emailReader stockports.UserEmailReader, mailer stockports.StockMailer, pushNotifier stockports.UserPushNotifier) *ApproveRequest {
	return &ApproveRequest{repo: repo, projectReader: projectReader, emailReader: emailReader, mailer: mailer, pushNotifier: pushNotifier}
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
	if err := uc.repo.ApproveRequest(ctx, input.RequestID); err != nil {
		return err
	}
	// Best-effort email + push to the requester.
	if email, err := uc.emailReader.FindEmailByID(ctx, req.RequestedByID); err == nil {
		_ = uc.mailer.NotifyRequesterApproved(ctx, email, resource.Name, req.Quantity)
	}
	uc.pushNotifier.Notify(ctx, req.RequestedByID,
		"Solicitud aprobada",
		resource.Name+" — tu reserva fue aprobada",
		"/stock/requests/my",
	)
	return nil
}
