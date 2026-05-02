package mealhttp

import (
	"errors"
	"net/http"

	mealdomain "macabi-back/internal/meal/domain"
)

func httpStatus(err error) int {
	switch {
	case errors.Is(err, mealdomain.ErrMealNotFound),
		errors.Is(err, mealdomain.ErrBookingNotFound),
		errors.Is(err, mealdomain.ErrTemplateNotFound),
		errors.Is(err, mealdomain.ErrGarnishNotFound):
		return http.StatusNotFound
	case errors.Is(err, mealdomain.ErrEmptyTitle),
		errors.Is(err, mealdomain.ErrInvalidMealType),
		errors.Is(err, mealdomain.ErrInvalidCategory),
		errors.Is(err, mealdomain.ErrInvalidDate),
		errors.Is(err, mealdomain.ErrInvalidAvailableCount),
		errors.Is(err, mealdomain.ErrEmptyGarnishName),
		errors.Is(err, mealdomain.ErrGarnishRequired),
		errors.Is(err, mealdomain.ErrGarnishNotForMeal):
		return http.StatusBadRequest
	case errors.Is(err, mealdomain.ErrBookingDeadlinePassed):
		return http.StatusUnprocessableEntity
	case errors.Is(err, mealdomain.ErrMealSoldOut),
		errors.Is(err, mealdomain.ErrAlreadyBooked),
		errors.Is(err, mealdomain.ErrProjectConflict):
		return http.StatusConflict
	case errors.Is(err, mealdomain.ErrBookingNotOwned):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
