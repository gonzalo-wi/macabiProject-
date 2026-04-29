package userports

import (
	"context"
	"time"
)

type PasswordResetToken struct {
	ID     string
	UserID string
}

type PasswordResetTokenRepository interface {
	InvalidateUnusedByUserID(ctx context.Context, userID string) error
	Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	FindValidByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
}

type PasswordResetMailer interface {
	SendResetLink(ctx context.Context, toEmail, resetURL string) error
}
