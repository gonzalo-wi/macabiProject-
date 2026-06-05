package eventports

import (
	"context"

	eventdomain "macabi-back/internal/event/domain"
)

type ModuleRepository interface {
	CreateModule(ctx context.Context, m *eventdomain.EventModule) error
	UpdateModule(ctx context.Context, m *eventdomain.EventModule) error
	DeleteModule(ctx context.Context, id string) error
	FindModuleByID(ctx context.Context, id string) (*eventdomain.EventModule, error)
	SetModuleProjects(ctx context.Context, moduleID string, projectIDs []string) error
}
