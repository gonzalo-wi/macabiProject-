package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
)

type CreateProjectInput struct {
	Name        string
	Description string
	AdminUserID string
}

type CreateProject struct {
	repo projectports.ProjectRepository
}

func NewCreateProject(repo projectports.ProjectRepository) *CreateProject {
	return &CreateProject{repo: repo}
}

func (uc *CreateProject) Execute(ctx context.Context, input CreateProjectInput) (*projectdomain.Project, error) {
	p, err := projectdomain.NewProject(input.Name, input.Description, input.AdminUserID)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
