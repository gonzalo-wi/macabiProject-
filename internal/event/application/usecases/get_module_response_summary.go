package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type GetModuleResponseSummary struct {
	repo eventports.ResponseRepository
}

func NewGetModuleResponseSummary(repo eventports.ResponseRepository) *GetModuleResponseSummary {
	return &GetModuleResponseSummary{repo: repo}
}

func (uc *GetModuleResponseSummary) Execute(ctx context.Context, moduleID string) (*eventdomain.ModuleResponseSummary, error) {
	return uc.repo.LoadModuleResponseSummary(ctx, moduleID)
}
