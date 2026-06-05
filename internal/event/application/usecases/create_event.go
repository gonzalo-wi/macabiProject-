package eventusecases

import (
	"context"
	"strings"
	"time"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type CreateEventInput struct {
	Title    string
	Type     eventdomain.EventType
	StartsAt time.Time
	Deadline *time.Time
	Status   eventdomain.EventStatus
}

type CreateEvent struct {
	repo eventports.EventRepository
}

func NewCreateEvent(repo eventports.EventRepository) *CreateEvent {
	return &CreateEvent{repo: repo}
}

func (uc *CreateEvent) Execute(ctx context.Context, input CreateEventInput) (*eventdomain.EventInstance, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	e := &eventdomain.EventInstance{
		Title:              title,
		Type:               input.Type,
		StartsAt:           input.StartsAt,
		ResponseDeadlineAt: input.Deadline,
		Status:             input.Status,
	}
	if err := uc.repo.CreateInstance(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}
