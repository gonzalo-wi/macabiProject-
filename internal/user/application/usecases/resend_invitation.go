package userusecases

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type ResendUserInvitation struct {
	users   userports.UserRepository
	invites userports.UserInvitationRepository
	mailer  userports.InvitationMailer
	baseURL string
	ttl     time.Duration
}

func NewResendUserInvitation(
	users userports.UserRepository,
	invites userports.UserInvitationRepository,
	mailer userports.InvitationMailer,
	frontendBaseURL string,
	ttl time.Duration,
) *ResendUserInvitation {
	return &ResendUserInvitation{
		users:   users,
		invites: invites,
		mailer:  mailer,
		baseURL: strings.TrimRight(strings.TrimSpace(frontendBaseURL), "/"),
		ttl:     ttl,
	}
}

func (uc *ResendUserInvitation) Execute(ctx context.Context, invitationID string) error {
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
	if _, err := uc.users.FindByEmail(ctx, inv.Email); err == nil {
		return userdomain.ErrEmailAlreadyTaken
	} else if !errors.Is(err, userdomain.ErrUserNotFound) {
		return err
	}
	if err := uc.invites.InvalidatePendingByEmail(ctx, inv.Email); err != nil {
		return err
	}
	raw, hash, err := generateInviteToken()
	if err != nil {
		return err
	}
	expires := time.Now().Add(uc.ttl)
	if err := uc.invites.Create(ctx, inv.Email, inv.Name, inv.Role, hash, expires); err != nil {
		return err
	}
	acceptURL := fmt.Sprintf("%s/aceptar-invitacion?token=%s", uc.baseURL, url.QueryEscape(raw))
	return uc.mailer.SendInvitationLink(ctx, inv.Email, acceptURL)
}
