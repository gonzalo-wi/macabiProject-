package eventports

import (
	"context"
	"time"
)

type EventMailer interface {
	SendEventReminder(ctx context.Context, toEmail, userName, eventTitle string, deadline time.Time, eventID string) error
}
