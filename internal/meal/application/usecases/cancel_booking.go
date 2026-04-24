package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type CancelBooking struct {
	bookingRepo mealports.BookingRepository
	mealRepo    mealports.MealRepository
}

func NewCancelBooking(bookingRepo mealports.BookingRepository, mealRepo mealports.MealRepository) *CancelBooking {
	return &CancelBooking{bookingRepo: bookingRepo, mealRepo: mealRepo}
}

type CancelBookingInput struct {
	BookingID string
	UserID    string
}

func (uc *CancelBooking) Execute(ctx context.Context, input CancelBookingInput) error {
	booking, err := uc.bookingRepo.FindByID(ctx, input.BookingID)
	if err != nil {
		return err
	}
	if booking.UserID != input.UserID {
		return mealdomain.ErrBookingNotOwned
	}

	meal, err := uc.mealRepo.FindByID(ctx, booking.MealID)
	if err != nil {
		return fmt.Errorf("find meal: %w", err)
	}
	meal.IncrementAvailable()
	if err := uc.mealRepo.Update(ctx, meal); err != nil {
		return fmt.Errorf("update meal: %w", err)
	}
	if err := uc.bookingRepo.Delete(ctx, input.BookingID); err != nil {
		return fmt.Errorf("cancel booking: %w", err)
	}
	return nil
}
