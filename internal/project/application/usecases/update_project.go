package projectusecases

import (
	"context"
	"strings"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
)

type UpdateProjectInput struct {
	ID          string
	Name        string
	Description string
	AdminUserID string
	Capacity    int
}

type UpdateProject struct {
	repo projectports.ProjectRepository
}

func NewUpdateProject(repo projectports.ProjectRepository) *UpdateProject {
	return &UpdateProject{repo: repo}
}

func (uc *UpdateProject) Execute(ctx context.Context, input UpdateProjectInput) (*projectdomain.Project, error) {
	p, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, projectdomain.ErrEmptyName
	}
	if input.AdminUserID == "" {
		return nil, projectdomain.ErrMissingAdmin
	}

	p.Name = name
	p.Description = input.Description
	p.AdminUserID = input.AdminUserID
	p.Capacity = input.Capacity

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
