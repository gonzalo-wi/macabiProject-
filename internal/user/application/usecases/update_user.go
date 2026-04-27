package userusecases

import (
	"context"
	"fmt"
	"strings"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type UpdateUser struct {
	repo userports.UserRepository
}

func NewUpdateUser(repo userports.UserRepository) *UpdateUser {
	return &UpdateUser{repo: repo}
}

type UpdateUserInput struct {
	TargetUserID string
	Name         string
	Email        string
}

func (uc *UpdateUser) Execute(ctx context.Context, input UpdateUserInput) (*userdomain.User, error) {
	user, err := uc.repo.FindByID(ctx, input.TargetUserID)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		name := strings.TrimSpace(input.Name)
		if name == "" {
			return nil, userdomain.ErrEmptyName
		}
		user.Name = name
	}

	if input.Email != "" {
		email := strings.TrimSpace(strings.ToLower(input.Email))
		existing, err := uc.repo.FindByEmail(ctx, email)
		if err != nil && err != userdomain.ErrUserNotFound {
			return nil, fmt.Errorf("check email: %w", err)
		}
		if existing != nil && existing.ID != user.ID {
			return nil, userdomain.ErrEmailAlreadyTaken
		}
		user.Email = email
	}

	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return user, nil
}
