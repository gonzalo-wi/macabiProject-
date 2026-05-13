package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type ListEvents struct {
	repo eventports.Repository
}

func NewListEvents(repo eventports.Repository) *ListEvents {
	return &ListEvents{repo: repo}
}

func (uc *ListEvents) Execute(ctx context.Context, params pagination.Params) (pagination.Result[eventdomain.EventInstance], error) {
	return uc.repo.ListInstances(ctx, params)
}
