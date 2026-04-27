package userusecases

import (
	"context"
	"fmt"

	userports "macabi-back/internal/user/application/ports"
)

type SetUserStatus struct {
	repo userports.UserRepository
}

func NewSetUserStatus(repo userports.UserRepository) *SetUserStatus {
	return &SetUserStatus{repo: repo}
}

type SetUserStatusInput struct {
	TargetUserID string
	Active       bool
}

func (uc *SetUserStatus) Execute(ctx context.Context, input SetUserStatusInput) error {
	user, err := uc.repo.FindByID(ctx, input.TargetUserID)
	if err != nil {
		return err
	}

	if input.Active {
		user.Activate()
	} else {
		user.Deactivate()
	}

	if err := uc.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("set user status: %w", err)
	}
	return nil
}
