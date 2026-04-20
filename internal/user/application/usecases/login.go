package userusecases

import (
	"context"
	"fmt"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type Login struct {
	repo     userports.UserRepository
	hasher   userports.PasswordHasher
	tokenPrv userports.TokenProvider
}

func NewLogin(repo userports.UserRepository, hasher userports.PasswordHasher, tokenPrv userports.TokenProvider) *Login {
	return &Login{repo: repo, hasher: hasher, tokenPrv: tokenPrv}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
	User  *userdomain.User
}

func (uc *Login) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, userdomain.ErrInvalidCredentials
	}
	if err := uc.hasher.Compare(user.Password, input.Password); err != nil {
		return nil, userdomain.ErrInvalidCredentials
	}
	token, err := uc.tokenPrv.Generate(userports.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	return &LoginOutput{Token: token, User: user}, nil
}
