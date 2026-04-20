package userusecases

import (
	"context"
	"fmt"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type RegisterUser struct {
	repo   userports.UserRepository
	hasher userports.PasswordHasher
}

func NewRegisterUser(repo userports.UserRepository, hasher userports.PasswordHasher) *RegisterUser {
	return &RegisterUser{repo: repo, hasher: hasher}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

func (uc *RegisterUser) Execute(ctx context.Context, input RegisterInput) (*userdomain.User, error) {
	if err := userdomain.ValidateRawPassword(input.Password); err != nil {
		return nil, err
	}
	existing, _ := uc.repo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, userdomain.ErrEmailAlreadyTaken
	}
	hashed, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user, err := userdomain.NewUser(input.Name, input.Email, hashed)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	return user, nil
}
