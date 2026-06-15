package eventservices

import (
	"context"
	"fmt"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/notifications"
)

// EventReminderNotifier unifica campana in-app, Web Push y email para recordatorios de jornadas.
type EventReminderNotifier struct {
	notifs eventports.EventNotificationRepository
	mailer eventports.EventMailer
	push   notifications.PushNotifier
}

func NewEventReminderNotifier(
	notifs eventports.EventNotificationRepository,
	mailer eventports.EventMailer,
	push notifications.PushNotifier,
) *EventReminderNotifier {
	return &EventReminderNotifier{notifs: notifs, mailer: mailer, push: push}
}

func (s *EventReminderNotifier) NotifyResponseReminder(
	ctx context.Context,
	target eventports.ReminderTarget,
	event eventdomain.EventInstance,
) {
	if event.ResponseDeadlineAt == nil {
		return
	}

	exists, err := s.notifs.HasReminderForEvent(ctx, target.UserID, event.ID)
	if err != nil || exists {
		return
	}

	msg := fmt.Sprintf("Recordatorio: mañana vence el plazo para responder \"%s\"", event.Title)
	title := "Recordatorio de jornada"
	body := fmt.Sprintf("Respondé antes de mañana — %s", event.Title)
	url := notifications.EventRespond(event.ID)

	_ = s.notifs.SaveNotification(ctx, &eventdomain.EventNotification{
		UserID:          target.UserID,
		EventInstanceID: event.ID,
		Message:         msg,
	})
	notifications.PushToUser(ctx, s.push, target.UserID, title, body, url)

	_ = s.mailer.SendEventReminder(
		ctx,
		target.Email,
		target.Name,
		event.Title,
		*event.ResponseDeadlineAt,
		event.ID,
	)
}
