package eventusecases

import (
	"context"
	"strings"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type UpdateOptionGroupInput struct {
	ID         string
	Name       string
	Type       eventdomain.OptionGroupType
	SortOrder  int
	IsRequired bool
}

type UpdateOptionGroup struct {
	repo eventports.Repository
}

func NewUpdateOptionGroup(repo eventports.Repository) *UpdateOptionGroup {
	return &UpdateOptionGroup{repo: repo}
}

func (uc *UpdateOptionGroup) Execute(ctx context.Context, input UpdateOptionGroupInput) (*eventdomain.EventOptionGroup, error) {
	g, err := uc.repo.FindOptionGroupByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	g.Name = name
	g.Type = input.Type
	g.SortOrder = input.SortOrder
	g.IsRequired = input.IsRequired
	if err := uc.repo.UpdateOptionGroup(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}
