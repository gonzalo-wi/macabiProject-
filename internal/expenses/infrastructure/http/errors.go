package expenseshttp

import (
	"errors"
	"net/http"

	expensesdomain "macabi-back/internal/expenses/domain"
)

func httpStatus(err error) int {
	switch {
	case errors.Is(err, expensesdomain.ErrExpenseNotFound):
		return http.StatusNotFound
	case errors.Is(err, expensesdomain.ErrNoReceiptAttached),
		errors.Is(err, expensesdomain.ErrNotificationNotFound):
		return http.StatusNotFound
	case errors.Is(err, expensesdomain.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, expensesdomain.ErrInvalidReceiptPath),
		errors.Is(err, expensesdomain.ErrInvalidAmount),
		errors.Is(err, expensesdomain.ErrMissingRequiredField),
		errors.Is(err, expensesdomain.ErrInvalidStatusTransition),
		errors.Is(err, expensesdomain.ErrInvalidMimeType),
		errors.Is(err, expensesdomain.ErrReceiptTooLarge),
		errors.Is(err, expensesdomain.ErrReceiptAttachFailed):
		return http.StatusBadRequest
	case errors.Is(err, expensesdomain.ErrReceiptStorageUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, expensesdomain.ErrCategoryNotFound):
		return http.StatusNotFound
	case errors.Is(err, expensesdomain.ErrCategoryInUse):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
