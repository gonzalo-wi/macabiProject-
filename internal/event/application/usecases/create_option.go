package eventusecases

import (
	"context"
	"strings"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type CreateOptionInput struct {
	GroupID     string
	Label       string
	MaxCapacity *int
	SortOrder   int
}

type CreateOption struct {
	repo eventports.OptionRepository
}

func NewCreateOption(repo eventports.OptionRepository) *CreateOption {
	return &CreateOption{repo: repo}
}

func (uc *CreateOption) Execute(ctx context.Context, input CreateOptionInput) (*eventdomain.EventOption, error) {
	o := &eventdomain.EventOption{
		GroupID:      input.GroupID,
		Label:        strings.TrimSpace(input.Label),
		MaxCapacity:  input.MaxCapacity,
		CurrentCount: 0,
		SortOrder:    input.SortOrder,
	}
	if o.Label == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	if err := uc.repo.CreateOption(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}
