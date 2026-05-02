package mealusecases

import (
	"context"
	"fmt"
	"time"

	mealports "macabi-back/internal/meal/application/ports"
	mealdomain "macabi-back/internal/meal/domain"
)

type GetDailySummary struct {
	repo mealports.BookingRepository
}

func NewGetDailySummary(repo mealports.BookingRepository) *GetDailySummary {
	return &GetDailySummary{repo: repo}
}

type GetDailySummaryInput struct {
	Date      time.Time
	ProjectID string // opcional; vacío = todos los proyectos
}

func (uc *GetDailySummary) Execute(ctx context.Context, input GetDailySummaryInput) (*mealdomain.DailySummary, error) {
	summary, err := uc.repo.GetDailySummary(ctx, input.Date, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get daily summary: %w", err)
	}
	return summary, nil
}
