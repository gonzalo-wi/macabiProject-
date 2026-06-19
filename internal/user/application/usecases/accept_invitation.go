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

		hashed, err := uc.hasher.Hash(password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		existing, err := uc.users.FindByEmail(txCtx, inv.Email)
		if err == nil {
			if existing.PasswordSet {
				return userdomain.ErrEmailAlreadyTaken
			}
			existing.Name = inv.Name
			existing.Role = inv.Role
			existing.SetPassword(hashed)
			if err := uc.users.Update(txCtx, existing); err != nil {
				return err
			}
			return uc.invites.MarkUsed(txCtx, inv.ID)
		}
		if !errors.Is(err, userdomain.ErrUserNotFound) {
			return err
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
