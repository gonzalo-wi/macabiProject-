package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
)

type DeleteMealTemplate struct {
	repo mealports.MealTemplateRepository
}

func NewDeleteMealTemplate(repo mealports.MealTemplateRepository) *DeleteMealTemplate {
	return &DeleteMealTemplate{repo: repo}
}

func (uc *DeleteMealTemplate) Execute(ctx context.Context, id string) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete meal template: %w", err)
	}
	return nil
}
