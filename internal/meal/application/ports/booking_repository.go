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
	// FindByUserAndDate devuelve la reserva del usuario para cualquier tipo de comida en esa fecha.
	// Se usa para prevenir tener almuerzo y cena el mismo día.
	FindByUserAndDate(ctx context.Context, userID string, date time.Time) (*mealdomain.Booking, error)
	Delete(ctx context.Context, id string) error
	GetDailySummary(ctx context.Context, date time.Time, projectID string) (*mealdomain.DailySummary, error)
}
