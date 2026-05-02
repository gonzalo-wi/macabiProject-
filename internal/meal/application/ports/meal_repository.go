package mealports

import (
	"context"
	"time"

	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/pagination"
)

// MealFilter permite filtrar viandas por criterios opcionales.
// Los campos vacíos no aplican filtro.
type MealFilter struct {
	ProjectID string
}

type MealRepository interface {
	Save(ctx context.Context, meal *mealdomain.Meal) error
	FindByID(ctx context.Context, id string) (*mealdomain.Meal, error)
	FindByDate(ctx context.Context, date time.Time, filter MealFilter, params pagination.Params) ([]mealdomain.Meal, int64, error)
	Update(ctx context.Context, meal *mealdomain.Meal) error
	Delete(ctx context.Context, id string) error
}
