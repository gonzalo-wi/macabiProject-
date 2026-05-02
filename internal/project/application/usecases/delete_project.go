package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
)

type DeleteProject struct {
	repo projectports.ProjectRepository
}

func NewDeleteProject(repo projectports.ProjectRepository) *DeleteProject {
	return &DeleteProject{repo: repo}
}

func (uc *DeleteProject) Execute(ctx context.Context, id string) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
