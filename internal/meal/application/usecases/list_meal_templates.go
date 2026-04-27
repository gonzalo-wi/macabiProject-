package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/pagination"
)

type ListMealTemplates struct {
	repo mealports.MealTemplateRepository
}

func NewListMealTemplates(repo mealports.MealTemplateRepository) *ListMealTemplates {
	return &ListMealTemplates{repo: repo}
}

func (uc *ListMealTemplates) Execute(ctx context.Context, params pagination.Params) (pagination.Result[mealdomain.MealTemplate], error) {
	templates, total, err := uc.repo.FindAll(ctx, params)
	if err != nil {
		return pagination.Result[mealdomain.MealTemplate]{}, fmt.Errorf("list meal templates: %w", err)
	}
	return pagination.NewResult(templates, total, params), nil
}
