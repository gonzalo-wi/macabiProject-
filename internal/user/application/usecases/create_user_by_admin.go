package userusecases

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type CreateUserByAdmin struct {
	repo   userports.UserRepository
	hasher userports.PasswordHasher
}

func NewCreateUserByAdmin(repo userports.UserRepository, hasher userports.PasswordHasher) *CreateUserByAdmin {
	return &CreateUserByAdmin{repo: repo, hasher: hasher}
}

type CreateUserByAdminInput struct {
	Name          string
	Email         string
	RequestedRole string
	InviterRole   userdomain.Role
}

func (uc *CreateUserByAdmin) Execute(ctx context.Context, in CreateUserByAdminInput) (*userdomain.User, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if name == "" {
		return nil, userdomain.ErrEmptyName
	}
	if email == "" {
		return nil, userdomain.ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, userdomain.ErrInvalidEmail
	}

	role, err := resolveInvitationRole(in.InviterRole, in.RequestedRole)
	if err != nil {
		return nil, err
	}

	if _, err := uc.repo.FindByEmail(ctx, email); err == nil {
		return nil, userdomain.ErrEmailAlreadyTaken
	} else if !errors.Is(err, userdomain.ErrUserNotFound) {
		return nil, err
	}

	placeholder, err := generatePlaceholderPassword()
	if err != nil {
		return nil, err
	}
	hashed, err := uc.hasher.Hash(placeholder)
	if err != nil {
		return nil, fmt.Errorf("hash placeholder password: %w", err)
	}

	user, err := userdomain.NewDraftUser(name, email, hashed, role)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("create draft user: %w", err)
	}
	return user, nil
}

func generatePlaceholderPassword() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
