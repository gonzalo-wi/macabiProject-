package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
)

type RemoveProjectMember struct {
	repo projectports.ProjectRepository
}

func NewRemoveProjectMember(repo projectports.ProjectRepository) *RemoveProjectMember {
	return &RemoveProjectMember{repo: repo}
}

func (uc *RemoveProjectMember) Execute(ctx context.Context, projectID, userID string) error {
	return uc.repo.RemoveMember(ctx, projectID, userID)
}
