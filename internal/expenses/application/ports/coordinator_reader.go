package expensesports

import "context"

type ProjectCoordinatorReader interface {
	FindProjectCoordinators(ctx context.Context, projectID string) ([]string, error)
}
