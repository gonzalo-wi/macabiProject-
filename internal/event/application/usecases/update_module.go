package eventusecases

import (
	"context"
	"strings"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type UpdateModuleInput struct {
	ID         string
	Title      string
	Type       eventdomain.ModuleType
	SortOrder  int
	IsRequired bool
}

type UpdateModule struct {
	repo eventports.ModuleRepository
}

func NewUpdateModule(repo eventports.ModuleRepository) *UpdateModule {
	return &UpdateModule{repo: repo}
}

func (uc *UpdateModule) Execute(ctx context.Context, input UpdateModuleInput) (*eventdomain.EventModule, error) {
	m, err := uc.repo.FindModuleByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	m.Title = title
	m.Type = input.Type
	m.SortOrder = input.SortOrder
	m.IsRequired = input.IsRequired
	if err := uc.repo.UpdateModule(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}
