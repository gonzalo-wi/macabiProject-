package userusecases

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"macabi-back/internal/shared/database"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type CreateUserInvitation struct {
	users    userports.UserRepository
	invites  userports.UserInvitationRepository
	mailer   userports.InvitationMailer
	baseURL  string
	ttl      time.Duration
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
	Email          string
	Name           string
	RequestedRole  string
	InviterID      string
	InviterRole    userdomain.Role
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
	if err := uc.invites.Create(ctx, email, name, role, in.InviterID, hash, expires); err != nil {
		return err
	}

	acceptURL := fmt.Sprintf("%s/aceptar-invitacion?token=%s", uc.baseURL, url.QueryEscape(raw))
	return uc.mailer.SendInvitationLink(ctx, email, acceptURL)
}

func resolveInvitationRole(inviter userdomain.Role, requested string) (userdomain.Role, error) {
	if !inviter.IsAtLeast(userdomain.RoleAdmin) {
		return "", userdomain.ErrForbidden
	}
	if inviter == userdomain.RoleAdmin {
		return userdomain.RoleUser, nil
	}
	// super_admin
	r := strings.TrimSpace(strings.ToLower(requested))
	if r == "" || r == "user" {
		return userdomain.RoleUser, nil
	}
	if r == "admin" {
		return userdomain.RoleAdmin, nil
	}
	return "", userdomain.ErrInvalidRole
}

type AcceptInvitation struct {
	transactor *database.GORMTransactor
	users      userports.UserRepository
	invites    userports.UserInvitationRepository
	hasher     userports.PasswordHasher
}

func NewAcceptInvitation(
	transactor *database.GORMTransactor,
	users userports.UserRepository,
	invites userports.UserInvitationRepository,
	hasher userports.PasswordHasher,
) *AcceptInvitation {
	return &AcceptInvitation{
		transactor: transactor,
		users:      users,
		invites:    invites,
		hasher:     hasher,
	}
}

func (uc *AcceptInvitation) Execute(ctx context.Context, rawToken, password string) error {
	if err := userdomain.ValidateRawPassword(password); err != nil {
		return err
	}
	if strings.TrimSpace(rawToken) == "" {
		return userdomain.ErrInvalidOrExpiredInvitation
	}
	hashHex := hashInviteToken(rawToken)

	return uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		inv, err := uc.invites.FindValidByTokenHash(txCtx, hashHex)
		if err != nil {
			return err
		}
		if inv == nil {
			return userdomain.ErrInvalidOrExpiredInvitation
		}

		if _, err := uc.users.FindByEmail(txCtx, inv.Email); err == nil {
			return userdomain.ErrEmailAlreadyTaken
		} else if !errors.Is(err, userdomain.ErrUserNotFound) {
			return err
		}

		hashed, err := uc.hasher.Hash(password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		user, err := userdomain.NewUserWithRole(inv.Name, inv.Email, hashed, inv.Role)
		if err != nil {
			return err
		}
		if err := uc.users.Save(txCtx, user); err != nil {
			return err
		}
		return uc.invites.MarkUsed(txCtx, inv.ID)
	})
}

func generateInviteToken() (raw string, hashHex string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("rand: %w", err)
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	return raw, hex.EncodeToString(sum[:]), nil
}

func hashInviteToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

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
	if err := uc.invites.Create(ctx, inv.Email, inv.Name, inv.Role, inv.InvitedByUserID, hash, expires); err != nil {
		return err
	}
	acceptURL := fmt.Sprintf("%s/aceptar-invitacion?token=%s", uc.baseURL, url.QueryEscape(raw))
	return uc.mailer.SendInvitationLink(ctx, inv.Email, acceptURL)
}

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
