package userusecases

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type CreateUserInvitation struct {
	users   userports.UserRepository
	invites userports.UserInvitationRepository
	mailer  userports.InvitationMailer
	baseURL string
	ttl     time.Duration
}

func NewCreateUserInvitation(
	users userports.UserRepository,
	invites userports.UserInvitationRepository,
	mailer userports.InvitationMailer,
	frontendBaseURL string,
	ttl time.Duration,
) *CreateUserInvitation {
	return &CreateUserInvitation{
		users:   users,
		invites: invites,
		mailer:  mailer,
		baseURL: strings.TrimRight(strings.TrimSpace(frontendBaseURL), "/"),
		ttl:     ttl,
	}
}

type CreateUserInvitationInput struct {
	Email         string
	Name          string
	RequestedRole string
	InviterID     string
	InviterRole   userdomain.Role
}

func (uc *CreateUserInvitation) Execute(ctx context.Context, in CreateUserInvitationInput) error {
	name := strings.TrimSpace(in.Name)
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if name == "" {
		return userdomain.ErrEmptyName
	}
	if email == "" {
		return userdomain.ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return userdomain.ErrInvalidEmail
	}

	role, err := resolveInvitationRole(in.InviterRole, in.RequestedRole)
	if err != nil {
		return err
	}

	if _, err := uc.users.FindByEmail(ctx, email); err == nil {
		return userdomain.ErrEmailAlreadyTaken
	} else if !errors.Is(err, userdomain.ErrUserNotFound) {
		return err
	}

	if err := uc.invites.InvalidatePendingByEmail(ctx, email); err != nil {
		return err
	}

	raw, hash, err := generateInviteToken()
	if err != nil {
		return err
	}
	expires := time.Now().Add(uc.ttl)
	if err := uc.invites.Create(ctx, email, name, role, hash, expires); err != nil {
		return err
	}

	acceptURL := fmt.Sprintf("%s/aceptar-invitacion?token=%s", uc.baseURL, url.QueryEscape(raw))
	return uc.mailer.SendInvitationLink(ctx, email, acceptURL)
}

func resolveInvitationRole(inviter userdomain.Role, requested string) (userdomain.Role, error) {
	if !inviter.IsAtLeast(userdomain.RoleAdmin) {
		return "", userdomain.ErrForbidden
	}
	r := strings.TrimSpace(strings.ToLower(requested))
	if r == "" || r == "user" {
		return userdomain.RoleUser, nil
	}
	if r == "admin" {
		return userdomain.RoleAdmin, nil
	}
	return "", userdomain.ErrInvalidRole
}
