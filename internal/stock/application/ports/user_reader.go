package stockports

import "context"

// UserEmailReader provides narrow read access to user emails needed by the stock module.
type UserEmailReader interface {
	FindEmailByID(ctx context.Context, userID string) (string, error)
	FindEmailsByIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}
