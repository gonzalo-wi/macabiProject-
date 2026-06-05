package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type DeleteModule struct {
	repo eventports.ModuleRepository
}

func NewDeleteModule(repo eventports.ModuleRepository) *DeleteModule {
	return &DeleteModule{repo: repo}
}

func (uc *DeleteModule) Execute(ctx context.Context, id string) error {
	return uc.repo.DeleteModule(ctx, id)
}
