package eventusecases

import (
	"context"
	"log"

	eventports "macabi-back/internal/event/application/ports"
)

type SendEventReminders struct {
	repo   eventports.ReminderRepository
	mailer eventports.EventMailer
}

func NewSendEventReminders(repo eventports.ReminderRepository, mailer eventports.EventMailer) *SendEventReminders {
	return &SendEventReminders{repo: repo, mailer: mailer}
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
			if err := uc.mailer.SendEventReminder(ctx, t.Email, t.Name, event.Title, *event.ResponseDeadlineAt, event.ID); err != nil {
				log.Printf("send_event_reminders: SendEventReminder to %s: %v", t.Email, err)
			}
		}
		log.Printf("send_event_reminders: sent %d reminders for event %q (%s)", len(targets), event.Title, event.ID)
	}
	return nil
}
