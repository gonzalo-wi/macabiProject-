package eventusecases

import (
	"context"
	"strings"
	"time"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type UpdateEventInput struct {
	ID       string
	Title    string
	Type     eventdomain.EventType
	StartsAt time.Time
	Deadline *time.Time
	Status   eventdomain.EventStatus
}

type UpdateEvent struct {
	repo eventports.Repository
}

func NewUpdateEvent(repo eventports.Repository) *UpdateEvent {
	return &UpdateEvent{repo: repo}
}

func (uc *UpdateEvent) Execute(ctx context.Context, input UpdateEventInput) (*eventdomain.EventInstance, error) {
	e, err := uc.repo.FindInstanceByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	e.Title = title
	e.Type = input.Type
	e.StartsAt = input.StartsAt
	e.ResponseDeadlineAt = input.Deadline
	e.Status = input.Status
	if err := uc.repo.UpdateInstance(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}
