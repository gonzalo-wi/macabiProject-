package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
)

type GetProject struct {
	repo projectports.ProjectRepository
}

func NewGetProject(repo projectports.ProjectRepository) *GetProject {
	return &GetProject{repo: repo}
}

func (uc *GetProject) Execute(ctx context.Context, id string) (*projectdomain.Project, error) {
	return uc.repo.FindByID(ctx, id)
}
