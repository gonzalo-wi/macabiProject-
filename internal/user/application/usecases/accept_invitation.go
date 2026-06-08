package userusecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"macabi-back/internal/shared/database"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

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
