package userdomain

import (
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID          string
	Name        string
	Email       string
	Password    string
	Role        Role
	Active      bool
	PasswordSet bool
	CreatedAt   time.Time
}

func NewUser(name, email, hashedPassword string) (*User, error) {
	return NewUserWithRole(name, email, hashedPassword, RoleUser)
}

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
		Name:        name,
		Email:       email,
		Password:    hashedPassword,
		Role:        role,
		Active:      true,
		PasswordSet: true,
	}, nil
}

func NewDraftUser(name, email, placeholderPasswordHash string, role Role) (*User, error) {
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
		Name:        name,
		Email:       email,
		Password:    placeholderPasswordHash,
		Role:        role,
		Active:      false,
		PasswordSet: false,
	}, nil
}

func (u *User) SetPassword(hashedPassword string) {
	u.Password = hashedPassword
	u.PasswordSet = true
	u.Active = true
}

func InvitationStatus(u *User, hasPendingInvitation bool) string {
	if u.PasswordSet {
		if u.Active {
			return "active"
		}
		return "inactive"
	}
	if hasPendingInvitation {
		return "invited"
	}
	return "draft"
}

func ValidateRawPassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

func (u *User) ChangeRole(newRole Role, changedBy *User) error {
	if !changedBy.Role.IsAtLeast(RoleAdmin) {
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
