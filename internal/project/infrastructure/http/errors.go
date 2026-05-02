package projecthttp

import (
	"errors"
	"net/http"

	projectdomain "macabi-back/internal/project/domain"
)

func httpStatus(err error) int {
	switch {
	case errors.Is(err, projectdomain.ErrProjectNotFound):
		return http.StatusNotFound
	case errors.Is(err, projectdomain.ErrEmptyName),
		errors.Is(err, projectdomain.ErrMissingAdmin):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
