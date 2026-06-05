package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type DeleteOption struct {
	repo eventports.OptionRepository
}

func NewDeleteOption(repo eventports.OptionRepository) *DeleteOption {
	return &DeleteOption{repo: repo}
}

func (uc *DeleteOption) Execute(ctx context.Context, id string) error {
	return uc.repo.DeleteOption(ctx, id)
}
