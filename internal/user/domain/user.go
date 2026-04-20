package userdomain

import (
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email, hashedPassword string) (*User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}
	email = strings.TrimSpace(strings.ToLower(email))
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	return &User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     RoleUser,
	}, nil
}

func ValidateRawPassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

func (u *User) ChangeRole(newRole Role, changedBy *User) error {
	if !changedBy.Role.IsAtLeast(RoleSuperAdmin) {
		return ErrForbidden
	}
	u.Role = newRole
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
