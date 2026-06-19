package userusecases

import (
	"context"

	"macabi-back/internal/shared/pagination"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type AdminUserRow struct {
	User                userdomain.User
	PendingInvitationID string
	InvitationStatus    string
}

type ListUsers struct {
	repo    userports.UserRepository
	invites userports.UserInvitationRepository
}

func NewListUsers(repo userports.UserRepository, invites userports.UserInvitationRepository) *ListUsers {
	return &ListUsers{repo: repo, invites: invites}
}

func (uc *ListUsers) Execute(ctx context.Context, filter userports.UserListFilter, params pagination.Params) (pagination.Result[AdminUserRow], error) {
	users, total, err := uc.repo.FindAll(ctx, filter, params)
	if err != nil {
		return pagination.Result[AdminUserRow]{}, err
	}

	rows := make([]AdminUserRow, len(users))
	for i := range users {
		pendingID := ""
		if !users[i].PasswordSet {
			inv, err := uc.invites.FindPendingByEmail(ctx, users[i].Email)
			if err != nil {
				return pagination.Result[AdminUserRow]{}, err
			}
			if inv != nil {
				pendingID = inv.ID
			}
		}
		rows[i] = AdminUserRow{
			User:                users[i],
			PendingInvitationID: pendingID,
			InvitationStatus:    userdomain.InvitationStatus(&users[i], pendingID != ""),
		}
	}
	return pagination.NewResult(rows, total, params), nil
}
