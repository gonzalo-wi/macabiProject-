package userports

import (
	"context"

	"macabi-back/internal/shared/pagination"
	userdomain "macabi-back/internal/user/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *userdomain.User) error
	FindByEmail(ctx context.Context, email string) (*userdomain.User, error)
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
	FindAll(ctx context.Context, params pagination.Params) ([]userdomain.User, int64, error)
	Update(ctx context.Context, user *userdomain.User) error
}
