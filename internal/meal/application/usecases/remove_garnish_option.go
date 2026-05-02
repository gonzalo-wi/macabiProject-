package mealusecases

import (
	"context"
	"fmt"

	mealports "macabi-back/internal/meal/application/ports"
)

type RemoveGarnishOption struct {
	repo mealports.MealTemplateRepository
}

func NewRemoveGarnishOption(repo mealports.MealTemplateRepository) *RemoveGarnishOption {
	return &RemoveGarnishOption{repo: repo}
}

type RemoveGarnishOptionInput struct {
	TemplateID string
	OptionID   string
}

func (uc *RemoveGarnishOption) Execute(ctx context.Context, input RemoveGarnishOptionInput) error {
	option, err := uc.repo.FindGarnishOptionByID(ctx, input.OptionID)
	if err != nil {
		return err
	}
	if option.TemplateID != input.TemplateID {
		return fmt.Errorf("garnish option does not belong to this template")
	}
	return uc.repo.DeleteGarnishOption(ctx, input.OptionID)
}
