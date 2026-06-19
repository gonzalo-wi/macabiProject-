package newshttp

import (
	"errors"
	"net/http"

	newsdomain "macabi-back/internal/news/domain"
)

func httpStatus(err error) int {
	switch {
	case errors.Is(err, newsdomain.ErrNewsNotFound),
		errors.Is(err, newsdomain.ErrNotificationNotFound),
		errors.Is(err, newsdomain.ErrNoImageAttached):
		return http.StatusNotFound
	case errors.Is(err, newsdomain.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, newsdomain.ErrMissingRequiredField),
		errors.Is(err, newsdomain.ErrInvalidStatusTransition),
		errors.Is(err, newsdomain.ErrInvalidMimeType),
		errors.Is(err, newsdomain.ErrImageTooLarge),
		errors.Is(err, newsdomain.ErrImageAttachFailed),
		errors.Is(err, newsdomain.ErrInvalidImagePath):
		return http.StatusBadRequest
	case errors.Is(err, newsdomain.ErrImageStorageUnavailable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
