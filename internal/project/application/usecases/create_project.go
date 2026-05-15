package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
)

type CreateProjectInput struct {
	Name          string
	Description   string
	CoordinatorID string
}

type CreateProject struct {
	repo projectports.ProjectRepository
}

func NewCreateProject(repo projectports.ProjectRepository) *CreateProject {
	return &CreateProject{repo: repo}
}

func (uc *CreateProject) Execute(ctx context.Context, input CreateProjectInput) (*projectdomain.Project, error) {
	if input.CoordinatorID == "" {
		return nil, projectdomain.ErrMissingCoordinator
	}
	p, err := projectdomain.NewProject(input.Name, input.Description)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	member := &projectdomain.ProjectMember{
		ProjectID: p.ID,
		UserID:    input.CoordinatorID,
		Role:      projectdomain.MemberRoleCoordinator,
	}
	if err := uc.repo.AddMember(ctx, member); err != nil {
		return nil, err
	}
	return p, nil
}
