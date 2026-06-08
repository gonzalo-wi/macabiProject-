package userusecases

import (
	"context"
	"errors"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type ListPendingInvitations struct {
	invites userports.UserInvitationRepository
	users   userports.UserRepository
}

func NewListPendingInvitations(
	invites userports.UserInvitationRepository,
	users userports.UserRepository,
) *ListPendingInvitations {
	return &ListPendingInvitations{invites: invites, users: users}
}

func (uc *ListPendingInvitations) Execute(ctx context.Context) ([]userports.PendingUserInvitation, error) {
	all, err := uc.invites.ListPending(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]userports.PendingUserInvitation, 0, len(all))
	for _, inv := range all {
		if _, err := uc.users.FindByEmail(ctx, inv.Email); err == nil {
			continue
		} else if !errors.Is(err, userdomain.ErrUserNotFound) {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, nil
}
