package mealusecases

import (
	"context"
	"fmt"
	"time"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type CreateMeal struct {
	repo mealports.MealRepository
}

func NewCreateMeal(repo mealports.MealRepository) *CreateMeal {
	return &CreateMeal{repo: repo}
}

type CreateMealInput struct {
	Title          string
	ImageURL       string
	Description    string
	Category       string
	Type           string
	AvailableCount int
	Date           time.Time
}

func (uc *CreateMeal) Execute(ctx context.Context, input CreateMealInput) (*mealdomain.Meal, error) {
	category, err := mealdomain.NewCategory(input.Category)
	if err != nil {
		return nil, err
	}
	mealType, err := mealdomain.NewMealType(input.Type)
	if err != nil {
		return nil, err
	}
	meal, err := mealdomain.NewMeal(input.Title, input.ImageURL, input.Description, category, mealType, input.AvailableCount, input.Date)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, meal); err != nil {
		return nil, fmt.Errorf("create meal: %w", err)
	}
	return meal, nil
}
