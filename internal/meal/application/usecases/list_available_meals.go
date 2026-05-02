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

type ListAvailableMealsInput struct {
	Date      time.Time
	ProjectID string // opcional; vacío = sin filtro
}

func (uc *ListAvailableMeals) Execute(ctx context.Context, input ListAvailableMealsInput, params pagination.Params) (pagination.Result[mealdomain.Meal], error) {
	filter := mealports.MealFilter{ProjectID: input.ProjectID}
	meals, total, err := uc.repo.FindByDate(ctx, input.Date, filter, params)
	if err != nil {
		return pagination.Result[mealdomain.Meal]{}, fmt.Errorf("list meals: %w", err)
	}
	return pagination.NewResult(meals, total, params), nil
}
