package userusecases

import (
	"context"
	"fmt"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type ChangePassword struct {
	repo   userports.UserRepository
	hasher userports.PasswordHasher
}

func NewChangePassword(repo userports.UserRepository, hasher userports.PasswordHasher) *ChangePassword {
	return &ChangePassword{repo: repo, hasher: hasher}
}

type ChangePasswordInput struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

func (uc *ChangePassword) Execute(ctx context.Context, input ChangePasswordInput) error {
	user, err := uc.repo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if err := uc.hasher.Compare(user.Password, input.CurrentPassword); err != nil {
		return userdomain.ErrWrongPassword
	}

	if err := userdomain.ValidateRawPassword(input.NewPassword); err != nil {
		return err
	}

	hashed, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashed

	if err := uc.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("change password: %w", err)
	}
	return nil
}
