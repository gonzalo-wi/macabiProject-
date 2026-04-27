package mealusecases

import (
	"context"
	"fmt"
	"time"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type CreateMeal struct {
	repo         mealports.MealRepository
	templateRepo mealports.MealTemplateRepository
}

func NewCreateMeal(repo mealports.MealRepository, templateRepo mealports.MealTemplateRepository) *CreateMeal {
	return &CreateMeal{repo: repo, templateRepo: templateRepo}
}

type CreateMealInput struct {
	TemplateID     string
	AvailableCount int
	Date           time.Time
}

func (uc *CreateMeal) Execute(ctx context.Context, input CreateMealInput) (*mealdomain.Meal, error) {
	tmpl, err := uc.templateRepo.FindByID(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}
	meal, err := mealdomain.NewMeal(input.TemplateID, input.AvailableCount, input.Date)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, meal); err != nil {
		return nil, fmt.Errorf("create meal: %w", err)
	}
	meal.Template = tmpl
	return meal, nil
}
