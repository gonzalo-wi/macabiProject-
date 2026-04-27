package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type CreateMealTemplate struct {
	repo mealports.MealTemplateRepository
}

func NewCreateMealTemplate(repo mealports.MealTemplateRepository) *CreateMealTemplate {
	return &CreateMealTemplate{repo: repo}
}

type CreateMealTemplateInput struct {
	Title       string
	ImageURL    string
	Description string
	Category    string
	Type        string
}

func (uc *CreateMealTemplate) Execute(ctx context.Context, input CreateMealTemplateInput) (*mealdomain.MealTemplate, error) {
	category, err := mealdomain.NewCategory(input.Category)
	if err != nil {
		return nil, err
	}
	mealType, err := mealdomain.NewMealType(input.Type)
	if err != nil {
		return nil, err
	}
	tmpl, err := mealdomain.NewMealTemplate(input.Title, input.ImageURL, input.Description, category, mealType)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, tmpl); err != nil {
		return nil, fmt.Errorf("create meal template: %w", err)
	}
	return tmpl, nil
}
