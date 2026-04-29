package userports

import (
	"context"
	"time"

	userdomain "macabi-back/internal/user/domain"
)

// UserInvitation is a pending invite row used by application services.
type UserInvitation struct {
	ID              string
	Email           string
	Name            string
	Role            userdomain.Role
	InvitedByUserID string
}

// PendingUserInvitation is a not-yet-accepted invite (used_at IS NULL), including expired tokens.
type PendingUserInvitation struct {
	ID              string
	Email           string
	Name            string
	Role            userdomain.Role
	InvitedByUserID string
	ExpiresAt       time.Time
	CreatedAt       time.Time
}

type UserInvitationRepository interface {
	InvalidatePendingByEmail(ctx context.Context, email string) error
	Create(ctx context.Context, email, name string, role userdomain.Role, invitedByUserID, tokenHash string, expiresAt time.Time) error
	FindValidByTokenHash(ctx context.Context, tokenHash string) (*UserInvitation, error)
	FindPendingByID(ctx context.Context, id string) (*PendingUserInvitation, error)
	ListPending(ctx context.Context) ([]PendingUserInvitation, error)
	MarkUsed(ctx context.Context, id string) error
}

type InvitationMailer interface {
	SendInvitationLink(ctx context.Context, toEmail, acceptURL string) error
}
