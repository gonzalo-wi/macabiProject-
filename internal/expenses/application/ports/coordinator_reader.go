package expensesports

import "context"

// ProjectCoordinatorReader lists coordinator user IDs for a project.
type ProjectCoordinatorReader interface {
	FindProjectCoordinators(ctx context.Context, projectID string) ([]string, error)
}
