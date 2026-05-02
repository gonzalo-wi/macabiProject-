package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
	"macabi-back/internal/shared/pagination"
)

type ListProjects struct {
	repo projectports.ProjectRepository
}

func NewListProjects(repo projectports.ProjectRepository) *ListProjects {
	return &ListProjects{repo: repo}
}

func (uc *ListProjects) Execute(ctx context.Context, params pagination.Params) (pagination.Result[projectdomain.Project], error) {
	return uc.repo.FindAll(ctx, params)
}
