package stockhttp

import (
	"net/http"

	stockdomain "macabi-back/internal/stock/domain"
)

func httpStatus(err error) int {
	switch err {
	case stockdomain.ErrResourceNotFound,
		stockdomain.ErrRequestNotFound,
		stockdomain.ErrNotificationNotFound:
		return http.StatusNotFound
	case stockdomain.ErrForbidden:
		return http.StatusForbidden
	case stockdomain.ErrInsufficientStock,
		stockdomain.ErrInvalidStatusTransition,
		stockdomain.ErrInvalidQuantity,
		stockdomain.ErrEmptyResourceName,
		stockdomain.ErrInvalidStock,
		stockdomain.ErrInvalidResourceType,
		stockdomain.ErrMissingRequiredField,
		stockdomain.ErrNotReturnable,
		stockdomain.ErrResourceInUse:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
