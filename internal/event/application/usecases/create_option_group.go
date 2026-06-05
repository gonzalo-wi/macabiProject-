package eventusecases

import (
	"context"
	"strings"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type CreateOptionGroupInput struct {
	ModuleID   string
	Name       string
	Type       eventdomain.OptionGroupType
	SortOrder  int
	IsRequired bool
}

type CreateOptionGroup struct {
	repo eventports.OptionRepository
}

func NewCreateOptionGroup(repo eventports.OptionRepository) *CreateOptionGroup {
	return &CreateOptionGroup{repo: repo}
}

func (uc *CreateOptionGroup) Execute(ctx context.Context, input CreateOptionGroupInput) (*eventdomain.EventOptionGroup, error) {
	g := &eventdomain.EventOptionGroup{
		ModuleID:   input.ModuleID,
		Name:       strings.TrimSpace(input.Name),
		Type:       input.Type,
		SortOrder:  input.SortOrder,
		IsRequired: input.IsRequired,
	}
	if g.Name == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	if err := uc.repo.CreateOptionGroup(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}
