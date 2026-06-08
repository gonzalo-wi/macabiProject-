package userusecases

import (
	"context"
	"strings"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type RevokeUserInvitation struct {
	invites userports.UserInvitationRepository
}

func NewRevokeUserInvitation(invites userports.UserInvitationRepository) *RevokeUserInvitation {
	return &RevokeUserInvitation{invites: invites}
}

func (uc *RevokeUserInvitation) Execute(ctx context.Context, invitationID string) error {
	id := strings.TrimSpace(invitationID)
	if id == "" {
		return userdomain.ErrInvitationNotFound
	}
	inv, err := uc.invites.FindPendingByID(ctx, id)
	if err != nil {
		return err
	}
	if inv == nil {
		return userdomain.ErrInvitationNotFound
	}
	return uc.invites.MarkUsed(ctx, id)
}
