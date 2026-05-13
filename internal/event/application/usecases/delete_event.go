package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type DeleteEvent struct {
	repo eventports.Repository
}

func NewDeleteEvent(repo eventports.Repository) *DeleteEvent {
	return &DeleteEvent{repo: repo}
}

func (uc *DeleteEvent) Execute(ctx context.Context, id string) error {
	return uc.repo.DeleteInstance(ctx, id)
}
