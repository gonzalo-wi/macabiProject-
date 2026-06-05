package eventports

import (
	"context"

	eventdomain "macabi-back/internal/event/domain"
)

type ReminderTarget struct {
	UserID string
	Name   string
	Email  string
}

type ReminderRepository interface {
	FindEventsWithDeadlineTomorrow(ctx context.Context) ([]eventdomain.EventInstance, error)
	FindUsersWithoutResponse(ctx context.Context, eventID string) ([]ReminderTarget, error)
}
