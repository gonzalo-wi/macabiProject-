package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type DeleteModule struct {
	repo eventports.Repository
}

func NewDeleteModule(repo eventports.Repository) *DeleteModule {
	return &DeleteModule{repo: repo}
}

func (uc *DeleteModule) Execute(ctx context.Context, id string) error {
	return uc.repo.DeleteModule(ctx, id)
}
