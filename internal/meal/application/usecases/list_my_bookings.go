package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/pagination"
)

type ListMyBookings struct {
	repo mealports.BookingRepository
}

func NewListMyBookings(repo mealports.BookingRepository) *ListMyBookings {
	return &ListMyBookings{repo: repo}
}

func (uc *ListMyBookings) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[mealdomain.Booking], error) {
	bookings, total, err := uc.repo.FindByUserID(ctx, userID, params)
	if err != nil {
		return pagination.Result[mealdomain.Booking]{}, fmt.Errorf("list bookings: %w", err)
	}
	return pagination.NewResult(bookings, total, params), nil
}
