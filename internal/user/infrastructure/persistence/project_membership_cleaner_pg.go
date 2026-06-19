package userpersistence

import (
	"context"

	projectpersistence "macabi-back/internal/project/infrastructure/persistence"
)

type ProjectMembershipCleanerPG struct {
	repo *projectpersistence.ProjectRepositoryPG
}

func NewProjectMembershipCleanerPG(repo *projectpersistence.ProjectRepositoryPG) *ProjectMembershipCleanerPG {
	return &ProjectMembershipCleanerPG{repo: repo}
}

func (c *ProjectMembershipCleanerPG) RemoveAllByUserID(ctx context.Context, userID string) error {
	return c.repo.RemoveAllMembersByUserID(ctx, userID)
}
