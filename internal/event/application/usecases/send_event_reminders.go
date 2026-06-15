package eventusecases

import (
	"context"
	"log"

	eventports "macabi-back/internal/event/application/ports"
	eventservices "macabi-back/internal/event/application/services"
)

type SendEventReminders struct {
	repo     eventports.ReminderRepository
	notifier *eventservices.EventReminderNotifier
}

func NewSendEventReminders(repo eventports.ReminderRepository, notifier *eventservices.EventReminderNotifier) *SendEventReminders {
	return &SendEventReminders{repo: repo, notifier: notifier}
}

func (uc *SendEventReminders) Execute(ctx context.Context) error {
	events, err := uc.repo.FindEventsWithDeadlineTomorrow(ctx)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		if event.ResponseDeadlineAt == nil {
			continue
		}
		targets, err := uc.repo.FindUsersWithoutResponse(ctx, event.ID)
		if err != nil {
			log.Printf("send_event_reminders: FindUsersWithoutResponse(%s): %v", event.ID, err)
			continue
		}
		for _, t := range targets {
			uc.notifier.NotifyResponseReminder(ctx, t, event)
		}
		log.Printf("send_event_reminders: notified %d users for event %q (%s)", len(targets), event.Title, event.ID)
	}
	return nil
}
