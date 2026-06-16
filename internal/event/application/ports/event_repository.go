package eventports

import (
	"context"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type EventListFilter struct {
	Query  string
	Status string
}

type EventRepository interface {
	CreateInstance(ctx context.Context, e *eventdomain.EventInstance) error
	UpdateInstance(ctx context.Context, e *eventdomain.EventInstance) error
	DeleteInstance(ctx context.Context, id string) error
	FindInstanceByID(ctx context.Context, id string) (*eventdomain.EventInstance, error)
	ListInstances(ctx context.Context, filter EventListFilter, params pagination.Params) (pagination.Result[eventdomain.EventInstance], error)
	LoadEventDetail(ctx context.Context, id string) (*eventdomain.EventDetail, error)
	SetInstanceProjects(ctx context.Context, eventID string, projectIDs []string) error
}
