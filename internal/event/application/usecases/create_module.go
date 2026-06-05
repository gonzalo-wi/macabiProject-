package eventusecases

import (
	"context"
	"strings"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type CreateModuleInput struct {
	EventInstanceID string
	Title           string
	Type            eventdomain.ModuleType
	SortOrder       int
	IsRequired      bool
}

type CreateModule struct {
	repo eventports.ModuleRepository
}

func NewCreateModule(repo eventports.ModuleRepository) *CreateModule {
	return &CreateModule{repo: repo}
}

func (uc *CreateModule) Execute(ctx context.Context, input CreateModuleInput) (*eventdomain.EventModule, error) {
	m := &eventdomain.EventModule{
		EventInstanceID: input.EventInstanceID,
		Title:           strings.TrimSpace(input.Title),
		Type:            input.Type,
		SortOrder:       input.SortOrder,
		IsRequired:      input.IsRequired,
	}
	if m.Title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	if err := uc.repo.CreateModule(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}
