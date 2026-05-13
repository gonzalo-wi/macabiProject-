package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type GetEventDetail struct {
	repo eventports.Repository
}

func NewGetEventDetail(repo eventports.Repository) *GetEventDetail {
	return &GetEventDetail{repo: repo}
}

func (uc *GetEventDetail) Execute(ctx context.Context, id string) (*eventdomain.EventDetail, error) {
	return uc.repo.LoadEventDetail(ctx, id)
}
