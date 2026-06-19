package userports

import (
	"context"
	"time"

	userdomain "macabi-back/internal/user/domain"
)

type UserInvitation struct {
	ID    string
	Email string
	Name  string
	Role  userdomain.Role
}

type PendingUserInvitation struct {
	ID        string
	Email     string
	Name      string
	Role      userdomain.Role
	ExpiresAt time.Time
	CreatedAt time.Time
}

type UserInvitationRepository interface {
	InvalidatePendingByEmail(ctx context.Context, email string) error
	Create(ctx context.Context, email, name string, role userdomain.Role, tokenHash string, expiresAt time.Time) error
	FindValidByTokenHash(ctx context.Context, tokenHash string) (*UserInvitation, error)
	FindPendingByID(ctx context.Context, id string) (*PendingUserInvitation, error)
	FindPendingByEmail(ctx context.Context, email string) (*PendingUserInvitation, error)
	ListPending(ctx context.Context) ([]PendingUserInvitation, error)
	MarkUsed(ctx context.Context, id string) error
}

type InvitationMailer interface {
	SendInvitationLink(ctx context.Context, toEmail, acceptURL string) error
}
