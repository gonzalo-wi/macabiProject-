package projectports

import "context"

type ProjectMembership interface {
	IsProjectCoordinator(ctx context.Context, projectID, userID string) (bool, error)
	IsProjectMember(ctx context.Context, projectID, userID string) (bool, error)
}

type ProjectCoordinatorReader interface {
	FindProjectCoordinators(ctx context.Context, projectID string) ([]string, error)
}

type ProjectMemberReader interface {
	ProjectMembership
	ProjectCoordinatorReader
}
