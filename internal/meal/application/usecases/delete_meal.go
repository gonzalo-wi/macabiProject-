package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
)

type DeleteMeal struct {
	repo mealports.MealRepository
}

func NewDeleteMeal(repo mealports.MealRepository) *DeleteMeal {
	return &DeleteMeal{repo: repo}
}

func (uc *DeleteMeal) Execute(ctx context.Context, id string) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete meal: %w", err)
	}
	return nil
}
