package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type SetEventProjects struct {
	repo eventports.EventRepository
}

func NewSetEventProjects(repo eventports.EventRepository) *SetEventProjects {
	return &SetEventProjects{repo: repo}
}

func (uc *SetEventProjects) Execute(ctx context.Context, eventID string, projectIDs []string) error {
	return uc.repo.SetInstanceProjects(ctx, eventID, projectIDs)
}
