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
	transactor  mealports.Transactor
}

func NewBookMeal(mealRepo mealports.MealRepository, bookingRepo mealports.BookingRepository, transactor mealports.Transactor) *BookMeal {
	return &BookMeal{mealRepo: mealRepo, bookingRepo: bookingRepo, transactor: transactor}
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

	var booking *mealdomain.Booking
	err = uc.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		// Re-fetch inside the transaction to get the latest inventory state.
		meal, err = uc.mealRepo.FindByID(ctx, input.MealID)
		if err != nil {
			return err
		}
		if meal.SoldOut {
			return mealdomain.ErrMealSoldOut
		}

		existing, err := uc.bookingRepo.FindByUserAndMealTypeAndDate(ctx, input.UserID, meal.Template.Type, meal.Date)
		if err != nil && !errors.Is(err, mealdomain.ErrBookingNotFound) {
			return fmt.Errorf("check existing booking: %w", err)
		}

		if existing != nil {
			prev, err := uc.mealRepo.FindByID(ctx, existing.MealID)
			if err != nil {
				return fmt.Errorf("find previous meal: %w", err)
			}
			prev.IncrementAvailable()
			if err := uc.mealRepo.Update(ctx, prev); err != nil {
				return fmt.Errorf("restore previous meal inventory: %w", err)
			}
			if err := uc.bookingRepo.Delete(ctx, existing.ID); err != nil {
				return fmt.Errorf("delete previous booking: %w", err)
			}
		}

		if err := meal.DecrementAvailable(); err != nil {
			return err
		}
		if err := uc.mealRepo.Update(ctx, meal); err != nil {
			return fmt.Errorf("update meal: %w", err)
		}

		booking = mealdomain.NewBooking(input.UserID, input.MealID)
		if err := uc.bookingRepo.Save(ctx, booking); err != nil {
			return fmt.Errorf("save booking: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	booking.Meal = meal
	return booking, nil
}
