package stockports

import "context"

type ProjectMemberReader interface {
	FindProjectCoordinators(ctx context.Context, projectID string) ([]string, error)
	IsProjectCoordinator(ctx context.Context, projectID, userID string) (bool, error)
	IsProjectMember(ctx context.Context, projectID, userID string) (bool, error)
}
