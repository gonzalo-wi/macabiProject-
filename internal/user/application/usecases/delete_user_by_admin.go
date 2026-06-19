package userusecases

import (
	"context"

	"macabi-back/internal/shared/database"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type DeleteUserByAdmin struct {
	transactor *database.GORMTransactor
	users      userports.UserRepository
	invites    userports.UserInvitationRepository
	members    userports.ProjectMembershipCleaner
}

func NewDeleteUserByAdmin(
	transactor *database.GORMTransactor,
	users userports.UserRepository,
	invites userports.UserInvitationRepository,
	members userports.ProjectMembershipCleaner,
) *DeleteUserByAdmin {
	return &DeleteUserByAdmin{
		transactor: transactor,
		users:      users,
		invites:    invites,
		members:    members,
	}
}

type DeleteUserByAdminInput struct {
	TargetUserID string
	ActorUserID  string
}

func (uc *DeleteUserByAdmin) Execute(ctx context.Context, in DeleteUserByAdminInput) error {
	if in.TargetUserID == in.ActorUserID {
		return userdomain.ErrCannotDeleteSelf
	}

	user, err := uc.users.FindByID(ctx, in.TargetUserID)
	if err != nil {
		return err
	}
	if user.PasswordSet {
		return userdomain.ErrUserAlreadyJoined
	}

	return uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.invites.InvalidatePendingByEmail(txCtx, user.Email); err != nil {
			return err
		}
		if err := uc.members.RemoveAllByUserID(txCtx, user.ID); err != nil {
			return err
		}
		return uc.users.Delete(txCtx, user.ID)
	})
}
