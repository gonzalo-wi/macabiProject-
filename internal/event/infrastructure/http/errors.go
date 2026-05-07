package eventhttp

import (
	"errors"
	"net/http"

	eventdomain "macabi-back/internal/event/domain"
)

func httpStatus(err error) int {
	switch {
	case errors.Is(err, eventdomain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, eventdomain.ErrEmptyTitle),
		errors.Is(err, eventdomain.ErrInvalidEventType),
		errors.Is(err, eventdomain.ErrInvalidEventStatus),
		errors.Is(err, eventdomain.ErrInvalidModuleType),
		errors.Is(err, eventdomain.ErrInvalidOptionGroupType):
		return http.StatusBadRequest
	case errors.Is(err, eventdomain.ErrOptionCapacityReached):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
