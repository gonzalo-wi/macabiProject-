package userports

import "context"

type UserEmailReader interface {
	FindEmailByID(ctx context.Context, userID string) (string, error)
	FindEmailsByIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}
