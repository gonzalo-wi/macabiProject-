package userports

import "context"

type ProjectMembershipCleaner interface {
	RemoveAllByUserID(ctx context.Context, userID string) error
}
