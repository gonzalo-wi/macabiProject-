package userusecases

import (
	"context"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type GetCurrentUser struct {
	repo userports.UserRepository
}

func NewGetCurrentUser(repo userports.UserRepository) *GetCurrentUser {
	return &GetCurrentUser{repo: repo}
}

func (uc *GetCurrentUser) Execute(ctx context.Context, userID string) (*userdomain.User, error) {
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, userdomain.ErrUserNotFound
	}
	return user, nil
}
