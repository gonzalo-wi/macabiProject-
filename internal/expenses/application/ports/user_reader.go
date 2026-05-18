package expensesports

import "context"

// UserEmailReader provides read access to user emails for expense notifications.
type UserEmailReader interface {
	FindEmailByID(ctx context.Context, userID string) (string, error)
	FindEmailsByIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}
