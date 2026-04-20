package userusecases

import (
	"context"
	"fmt"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type ChangeRole struct {
	repo userports.UserRepository
}

func NewChangeRole(repo userports.UserRepository) *ChangeRole {
	return &ChangeRole{repo: repo}
}

type ChangeRoleInput struct {
	TargetUserID string
	NewRole      string
	ChangedByID  string
}

func (uc *ChangeRole) Execute(ctx context.Context, input ChangeRoleInput) error {
	changedBy, err := uc.repo.FindByID(ctx, input.ChangedByID)
	if err != nil {
		return userdomain.ErrUnauthorized
	}
	newRole, err := userdomain.NewRole(input.NewRole)
	if err != nil {
		return err
	}
	target, err := uc.repo.FindByID(ctx, input.TargetUserID)
	if err != nil {
		return userdomain.ErrUserNotFound
	}
	if err := target.ChangeRole(newRole, changedBy); err != nil {
		return err
	}
	if err := uc.repo.Update(ctx, target); err != nil {
		return fmt.Errorf("change role: %w", err)
	}
	return nil
}
