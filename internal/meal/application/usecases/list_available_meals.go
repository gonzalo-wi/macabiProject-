package mealusecases

import (
	"context"
	"fmt"
	"time"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/pagination"
)

type ListAvailableMeals struct {
	repo mealports.MealRepository
}

func NewListAvailableMeals(repo mealports.MealRepository) *ListAvailableMeals {
	return &ListAvailableMeals{repo: repo}
}

func (uc *ListAvailableMeals) Execute(ctx context.Context, date time.Time, params pagination.Params) (pagination.Result[mealdomain.Meal], error) {
	meals, total, err := uc.repo.FindByDate(ctx, date, params)
	if err != nil {
		return pagination.Result[mealdomain.Meal]{}, fmt.Errorf("list meals: %w", err)
	}
	return pagination.NewResult(meals, total, params), nil
}
