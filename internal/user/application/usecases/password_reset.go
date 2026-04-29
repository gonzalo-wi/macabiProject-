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

type RequestPasswordReset struct {
	users   userports.UserRepository
	tokens  userports.PasswordResetTokenRepository
	mailer  userports.PasswordResetMailer
	baseURL string
	ttl     time.Duration
}

func NewRequestPasswordReset(
	users userports.UserRepository,
	tokens userports.PasswordResetTokenRepository,
	mailer userports.PasswordResetMailer,
	frontendBaseURL string,
	ttl time.Duration,
) *RequestPasswordReset {
	return &RequestPasswordReset{
		users:   users,
		tokens:  tokens,
		mailer:  mailer,
		baseURL: strings.TrimRight(strings.TrimSpace(frontendBaseURL), "/"),
		ttl:     ttl,
	}
}

func (uc *RequestPasswordReset) Execute(ctx context.Context, email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return userdomain.ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return userdomain.ErrInvalidEmail
	}

	user, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, userdomain.ErrUserNotFound) {
			return nil
		}
		return err
	}
	if !user.Active {
		return nil
	}

	if err := uc.tokens.InvalidateUnusedByUserID(ctx, user.ID); err != nil {
		return err
	}

	raw, hash, err := generateOpaqueResetToken()
	if err != nil {
		return err
	}
	expires := time.Now().Add(uc.ttl)
	if err := uc.tokens.Create(ctx, user.ID, hash, expires); err != nil {
		return err
	}

	resetURL := fmt.Sprintf("%s/restablecer-contrasena?token=%s", uc.baseURL, url.QueryEscape(raw))
	return uc.mailer.SendResetLink(ctx, user.Email, resetURL)
}

type ResetPassword struct {
	transactor *database.GORMTransactor
	users      userports.UserRepository
	tokens     userports.PasswordResetTokenRepository
	hasher     userports.PasswordHasher
}

func NewResetPassword(
	transactor *database.GORMTransactor,
	users userports.UserRepository,
	tokens userports.PasswordResetTokenRepository,
	hasher userports.PasswordHasher,
) *ResetPassword {
	return &ResetPassword{
		transactor: transactor,
		users:      users,
		tokens:     tokens,
		hasher:     hasher,
	}
}

func (uc *ResetPassword) Execute(ctx context.Context, rawToken, newPassword string) error {
	if err := userdomain.ValidateRawPassword(newPassword); err != nil {
		return err
	}
	if strings.TrimSpace(rawToken) == "" {
		return userdomain.ErrInvalidOrExpiredResetToken
	}

	hash := hashResetToken(rawToken)

	return uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		rec, err := uc.tokens.FindValidByTokenHash(txCtx, hash)
		if err != nil {
			return err
		}
		if rec == nil {
			return userdomain.ErrInvalidOrExpiredResetToken
		}

		user, err := uc.users.FindByID(txCtx, rec.UserID)
		if err != nil {
			return err
		}

		hashed, err := uc.hasher.Hash(newPassword)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		user.Password = hashed

		if err := uc.users.Update(txCtx, user); err != nil {
			return err
		}
		return uc.tokens.MarkUsed(txCtx, rec.ID)
	})
}

func generateOpaqueResetToken() (raw string, hashHex string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("rand: %w", err)
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	return raw, hex.EncodeToString(sum[:]), nil
}

func hashResetToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
