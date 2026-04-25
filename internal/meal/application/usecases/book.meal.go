package mealusecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type BookMeal struct {
	mealRepo    mealports.MealRepository
	bookingRepo mealports.BookingRepository
}

func NewBookMeal(mealRepo mealports.MealRepository, bookingRepo mealports.BookingRepository) *BookMeal {
	return &BookMeal{mealRepo: mealRepo, bookingRepo: bookingRepo}
}

type BookMealInput struct {
	UserID string
	MealID string
}

func (uc *BookMeal) Execute(ctx context.Context, input BookMealInput) (*mealdomain.Booking, error) {
	meal, err := uc.mealRepo.FindByID(ctx, input.MealID)
	if err != nil {
		return nil, err
	}

	if !mealdomain.IsBookingOpen(meal.Date, time.Now()) {
		return nil, mealdomain.ErrBookingDeadlinePassed
	}

	if meal.SoldOut {
		return nil, mealdomain.ErrMealSoldOut
	}
	existing, err := uc.bookingRepo.FindByUserAndMealTypeAndDate(ctx, input.UserID, meal.Type, meal.Date)
	if err == nil && existing != nil {
		prev, prevErr := uc.mealRepo.FindByID(ctx, existing.MealID)
		if prevErr == nil {
			prev.IncrementAvailable()
			_ = uc.mealRepo.Update(ctx, prev)
		}
		_ = uc.bookingRepo.Delete(ctx, existing.ID)
	} else if err != nil && !errors.Is(err, mealdomain.ErrBookingNotFound) {
		return nil, fmt.Errorf("check existing booking: %w", err)
	}

	if err := meal.DecrementAvailable(); err != nil {
		return nil, err
	}
	if err := uc.mealRepo.Update(ctx, meal); err != nil {
		return nil, fmt.Errorf("update meal: %w", err)
	}

	booking := mealdomain.NewBooking(input.UserID, input.MealID)
	if err := uc.bookingRepo.Save(ctx, booking); err != nil {
		return nil, fmt.Errorf("save booking: %w", err)
	}
	booking.Meal = meal
	return booking, nil
}
