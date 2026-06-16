package userusecases

import (
	"context"

	"macabi-back/internal/shared/pagination"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type ListUsers struct {
	repo userports.UserRepository
}

func NewListUsers(repo userports.UserRepository) *ListUsers {
	return &ListUsers{repo: repo}
}

func (uc *ListUsers) Execute(ctx context.Context, filter userports.UserListFilter, params pagination.Params) (pagination.Result[userdomain.User], error) {
	users, total, err := uc.repo.FindAll(ctx, filter, params)
	if err != nil {
		return pagination.Result[userdomain.User]{}, err
	}
	return pagination.NewResult(users, total, params), nil
}
