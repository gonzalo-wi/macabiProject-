package eventusecases

import (
	"context"
	"strings"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type UpdateOptionInput struct {
	ID           string
	Label        string
	MaxCapacity  *int
	SortOrder    int
	CurrentCount *int
}

type UpdateOption struct {
	repo eventports.OptionRepository
}

func NewUpdateOption(repo eventports.OptionRepository) *UpdateOption {
	return &UpdateOption{repo: repo}
}

func (uc *UpdateOption) Execute(ctx context.Context, input UpdateOptionInput) (*eventdomain.EventOption, error) {
	o, err := uc.repo.FindOptionByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	label := strings.TrimSpace(input.Label)
	if label == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	o.Label = label
	o.MaxCapacity = input.MaxCapacity
	o.SortOrder = input.SortOrder
	if input.CurrentCount != nil {
		o.CurrentCount = *input.CurrentCount
	}
	if err := uc.repo.UpdateOption(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}
