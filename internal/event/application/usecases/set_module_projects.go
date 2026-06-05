package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
)

type SetModuleProjects struct {
	repo eventports.ModuleRepository
}

func NewSetModuleProjects(repo eventports.ModuleRepository) *SetModuleProjects {
	return &SetModuleProjects{repo: repo}
}

func (uc *SetModuleProjects) Execute(ctx context.Context, moduleID string, projectIDs []string) error {
	return uc.repo.SetModuleProjects(ctx, moduleID, projectIDs)
}
