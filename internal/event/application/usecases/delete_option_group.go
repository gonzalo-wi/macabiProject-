package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type DeleteOptionGroup struct {
	repo eventports.Repository
}

func NewDeleteOptionGroup(repo eventports.Repository) *DeleteOptionGroup {
	return &DeleteOptionGroup{repo: repo}
}

func (uc *DeleteOptionGroup) Execute(ctx context.Context, id string) error {
	return uc.repo.DeleteOptionGroup(ctx, id)
}
