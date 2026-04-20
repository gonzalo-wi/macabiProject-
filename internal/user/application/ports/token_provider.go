package userports

import userdomain "macabi-back/internal/user/domain"

type TokenClaims struct {
	UserID string
	Email  string
	Role   userdomain.Role
}

type TokenProvider interface {
	Generate(claims TokenClaims) (string, error)
	Validate(token string) (*TokenClaims, error)
}
