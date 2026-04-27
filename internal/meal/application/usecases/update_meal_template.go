package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type UpdateMealTemplate struct {
	repo mealports.MealTemplateRepository
}

func NewUpdateMealTemplate(repo mealports.MealTemplateRepository) *UpdateMealTemplate {
	return &UpdateMealTemplate{repo: repo}
}

type UpdateMealTemplateInput struct {
	ID          string
	Title       string
	ImageURL    string
	Description string
	Category    string
	Type        string
}

func (uc *UpdateMealTemplate) Execute(ctx context.Context, input UpdateMealTemplateInput) (*mealdomain.MealTemplate, error) {
	tmpl, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		tmpl.Title = input.Title
	}
	if input.ImageURL != "" {
		tmpl.ImageURL = input.ImageURL
	}
	if input.Description != "" {
		tmpl.Description = input.Description
	}
	if input.Category != "" {
		cat, err := mealdomain.NewCategory(input.Category)
		if err != nil {
			return nil, err
		}
		tmpl.Category = cat
	}
	if input.Type != "" {
		mt, err := mealdomain.NewMealType(input.Type)
		if err != nil {
			return nil, err
		}
		tmpl.Type = mt
	}

	if err := uc.repo.Update(ctx, tmpl); err != nil {
		return nil, fmt.Errorf("update meal template: %w", err)
	}
	return tmpl, nil
}
