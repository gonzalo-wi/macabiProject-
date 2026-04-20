package userhttp

import (
	"errors"
	"net/http"

	userdomain "macabi-back/internal/user/domain"
)

func httpStatus(err error) int {
	switch {
	case errors.Is(err, userdomain.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, userdomain.ErrEmailAlreadyTaken):
		return http.StatusConflict
	case errors.Is(err, userdomain.ErrInvalidEmail),
		errors.Is(err, userdomain.ErrWeakPassword),
		errors.Is(err, userdomain.ErrEmptyName),
		errors.Is(err, userdomain.ErrInvalidRole):
		return http.StatusBadRequest
	case errors.Is(err, userdomain.ErrInvalidCredentials),
		errors.Is(err, userdomain.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, userdomain.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
