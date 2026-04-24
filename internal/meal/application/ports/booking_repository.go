package mealports

import (
	"context"
	"time"

	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/pagination"
)

type BookingRepository interface {
	Save(ctx context.Context, booking *mealdomain.Booking) error
	FindByID(ctx context.Context, id string) (*mealdomain.Booking, error)
	FindByUserID(ctx context.Context, userID string, params pagination.Params) ([]mealdomain.Booking, int64, error)
	FindByUserAndMealTypeAndDate(ctx context.Context, userID string, mealType mealdomain.MealType, date time.Time, isPostre bool) (*mealdomain.Booking, error)
	Delete(ctx context.Context, id string) error
}
