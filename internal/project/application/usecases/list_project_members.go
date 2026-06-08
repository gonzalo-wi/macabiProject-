package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
)

type ListProjectMembers struct {
	repo projectports.ProjectRepository
}

func NewListProjectMembers(repo projectports.ProjectRepository) *ListProjectMembers {
	return &ListProjectMembers{repo: repo}
}

func (uc *ListProjectMembers) Execute(ctx context.Context, projectID string) ([]projectdomain.ProjectMember, error) {
	if _, err := uc.repo.FindByID(ctx, projectID); err != nil {
		return nil, err
	}
	return uc.repo.ListMembers(ctx, projectID)
}
