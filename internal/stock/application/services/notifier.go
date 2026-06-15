package stockservices

import (
	"context"
	"fmt"

	projectports "macabi-back/internal/project/application/ports"
	"macabi-back/internal/shared/notifications"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userports "macabi-back/internal/user/application/ports"
)

// StockNotificationService unifica campana in-app, Web Push y email para stock.
type StockNotificationService struct {
	notifs       stockports.StockRepository
	coordinators projectports.ProjectCoordinatorReader
	emails       userports.UserEmailReader
	mailer       stockports.StockMailer
	push         notifications.PushNotifier
}

func NewStockNotificationService(
	notifs stockports.StockRepository,
	coordinators projectports.ProjectCoordinatorReader,
	emails userports.UserEmailReader,
	mailer stockports.StockMailer,
	push notifications.PushNotifier,
) *StockNotificationService {
	return &StockNotificationService{
		notifs:       notifs,
		coordinators: coordinators,
		emails:       emails,
		mailer:       mailer,
		push:         push,
	}
}

func (s *StockNotificationService) NotifyCoordinatorsNewRequest(
	ctx context.Context,
	requestID, projectID, resourceName string,
	quantity int,
) {
	coordinators, err := s.coordinators.FindProjectCoordinators(ctx, projectID)
	if err != nil || len(coordinators) == 0 {
		return
	}

	msg := fmt.Sprintf("Nueva solicitud de reserva: %d unidad(es) de \"%s\"", quantity, resourceName)
	title := "Nueva solicitud de reserva"
	body := fmt.Sprintf("%d unidad(es) de \"%s\" esperan aprobación", quantity, resourceName)
	url := notifications.StockRequestDetail(requestID)

	for _, coordinatorID := range coordinators {
		_ = s.notifs.SaveNotification(ctx, &stockdomain.StockNotification{
			UserID:    coordinatorID,
			RequestID: requestID,
			Message:   msg,
		})
		notifications.PushToUser(ctx, s.push, coordinatorID, title, body, url)
	}

	if emails, err := s.emails.FindEmailsByIDs(ctx, coordinators); err == nil {
		addrs := make([]string, 0, len(emails))
		for _, e := range emails {
			addrs = append(addrs, e)
		}
		_ = s.mailer.NotifyCoordinatorsNewRequest(ctx, addrs, resourceName, quantity, requestID)
	}
}

func (s *StockNotificationService) NotifyRequesterApproved(
	ctx context.Context,
	requestID, requesterID, resourceName string,
	quantity int,
) {
	msg := fmt.Sprintf("Tu pedido fue aprobado: %d unidad(es) de \"%s\"", quantity, resourceName)
	title := "Solicitud aprobada"
	body := resourceName + " — tu reserva fue aprobada"
	url := notifications.StockRequestDetail(requestID)

	_ = s.notifs.SaveNotification(ctx, &stockdomain.StockNotification{
		UserID:    requesterID,
		RequestID: requestID,
		Message:   msg,
	})
	notifications.PushToUser(ctx, s.push, requesterID, title, body, url)

	if email, err := s.emails.FindEmailByID(ctx, requesterID); err == nil {
		_ = s.mailer.NotifyRequesterApproved(ctx, email, resourceName, quantity)
	}
}

func (s *StockNotificationService) NotifyRequesterRejected(
	ctx context.Context,
	requestID, requesterID, resourceName string,
	quantity int,
) {
	msg := fmt.Sprintf("Tu pedido fue rechazado: %d unidad(es) de \"%s\"", quantity, resourceName)
	title := "Solicitud rechazada"
	body := resourceName + " — tu reserva fue rechazada"
	url := notifications.StockRequestDetail(requestID)

	_ = s.notifs.SaveNotification(ctx, &stockdomain.StockNotification{
		UserID:    requesterID,
		RequestID: requestID,
		Message:   msg,
	})
	notifications.PushToUser(ctx, s.push, requesterID, title, body, url)

	if email, err := s.emails.FindEmailByID(ctx, requesterID); err == nil {
		_ = s.mailer.NotifyRequesterRejected(ctx, email, resourceName, quantity)
	}
}
