package stockports

import "context"

// ProjectMemberReader provides read-only access to project membership data.
// It is intentionally kept narrow: only the queries the stock module needs.
type ProjectMemberReader interface {
	FindProjectCoordinators(ctx context.Context, projectID string) ([]string, error)
	IsProjectCoordinator(ctx context.Context, projectID, userID string) (bool, error)
}
