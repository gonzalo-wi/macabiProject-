package mealports

import (
	"context"

	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/pagination"
)

type MealTemplateRepository interface {
	Save(ctx context.Context, template *mealdomain.MealTemplate) error
	FindByID(ctx context.Context, id string) (*mealdomain.MealTemplate, error)
	FindAll(ctx context.Context, params pagination.Params) ([]mealdomain.MealTemplate, int64, error)
	Update(ctx context.Context, template *mealdomain.MealTemplate) error
	Delete(ctx context.Context, id string) error

	SaveGarnishOption(ctx context.Context, option *mealdomain.GarnishOption) error
	FindGarnishOptionByID(ctx context.Context, id string) (*mealdomain.GarnishOption, error)
	DeleteGarnishOption(ctx context.Context, id string) error
}
