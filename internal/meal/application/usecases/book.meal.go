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
	UserID          string
	MealID          string
	GarnishOptionID *string
}

func (uc *BookMeal) Execute(ctx context.Context, input BookMealInput) (*mealdomain.Booking, error) {
	meal, err := uc.mealRepo.FindByID(ctx, input.MealID)
	if err != nil {
		return nil, err
	}

	if !mealdomain.IsBookingOpen(meal.Date, time.Now()) {
		return nil, mealdomain.ErrBookingDeadlinePassed
	}

	if err := validateGarnish(meal, input.GarnishOptionID); err != nil {
		return nil, err
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

		existing, err := uc.bookingRepo.FindByUserAndDate(ctx, input.UserID, meal.Date)
		if err != nil && !errors.Is(err, mealdomain.ErrBookingNotFound) {
			return fmt.Errorf("check existing booking: %w", err)
		}

		if existing != nil {
			prev, err := uc.mealRepo.FindByID(ctx, existing.MealID)
			if err != nil {
				return fmt.Errorf("find previous meal: %w", err)
			}
			// Un usuario solo puede asistir a un proyecto por día.
			if prev.ProjectID != meal.ProjectID {
				return mealdomain.ErrProjectConflict
			}
			// No se puede tener almuerzo y cena el mismo día.
			if prev.Template.Type != meal.Template.Type {
				return mealdomain.ErrAlreadyBooked
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

		booking = mealdomain.NewBooking(input.UserID, input.MealID, input.GarnishOptionID)
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

// validateGarnish verifica que la guarnición sea válida para la vianda:
// - Si la vianda tiene opciones de guarnición, se debe elegir una.
// - La guarnición elegida debe pertenecer a la plantilla de la vianda.
func validateGarnish(meal *mealdomain.Meal, garnishOptionID *string) error {
	if len(meal.Template.GarnishOptions) == 0 {
		return nil
	}
	if garnishOptionID == nil {
		return mealdomain.ErrGarnishRequired
	}
	for _, g := range meal.Template.GarnishOptions {
		if g.ID == *garnishOptionID {
			return nil
		}
	}
	return mealdomain.ErrGarnishNotForMeal
}
