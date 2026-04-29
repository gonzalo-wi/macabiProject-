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
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email, hashedPassword string) (*User, error) {
	return NewUserWithRole(name, email, hashedPassword, RoleUser)
}

// NewUserWithRole validates name/email and builds a user with the given role (e.g. from an invitation).
func NewUserWithRole(name, email, hashedPassword string, role Role) (*User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}
	email = strings.TrimSpace(strings.ToLower(email))
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	if !validRoles[role] {
		return nil, ErrInvalidRole
	}
	return &User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
		Active:   true,
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

func (u *User) Activate()   { u.Active = true }
func (u *User) Deactivate() { u.Active = false }

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
