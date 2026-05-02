package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type AddGarnishOption struct {
	repo mealports.MealTemplateRepository
}

func NewAddGarnishOption(repo mealports.MealTemplateRepository) *AddGarnishOption {
	return &AddGarnishOption{repo: repo}
}

type AddGarnishOptionInput struct {
	TemplateID string
	Name       string
}

func (uc *AddGarnishOption) Execute(ctx context.Context, input AddGarnishOptionInput) (*mealdomain.GarnishOption, error) {
	if _, err := uc.repo.FindByID(ctx, input.TemplateID); err != nil {
		return nil, err
	}
	option, err := mealdomain.NewGarnishOption(input.TemplateID, input.Name)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.SaveGarnishOption(ctx, option); err != nil {
		return nil, fmt.Errorf("save garnish option: %w", err)
	}
	return option, nil
}
