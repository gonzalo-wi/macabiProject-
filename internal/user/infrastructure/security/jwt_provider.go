package usersecurity

import (
	"fmt"
	"time"

	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTProvider(secret string, expiration time.Duration) *JWTProvider {
	return &JWTProvider{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (p *JWTProvider) Generate(claims userports.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.UserID,
		"email": claims.Email,
		"role":  string(claims.Role),
		"exp":   time.Now().Add(p.expiration).Unix(),
		"iat":   time.Now().Unix(),
	})
	return token.SignedString(p.secret)
}

func (p *JWTProvider) Validate(tokenStr string) (*userports.TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return p.secret, nil
	})
	if err != nil {
		return nil, userdomain.ErrUnauthorized
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, userdomain.ErrUnauthorized
	}

	userID, ok := mapClaims["sub"].(string)
	if !ok {
		return nil, userdomain.ErrUnauthorized
	}
	email, ok := mapClaims["email"].(string)
	if !ok {
		return nil, userdomain.ErrUnauthorized
	}
	role, ok := mapClaims["role"].(string)
	if !ok {
		return nil, userdomain.ErrUnauthorized
	}

	return &userports.TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   userdomain.Role(role),
	}, nil
}
