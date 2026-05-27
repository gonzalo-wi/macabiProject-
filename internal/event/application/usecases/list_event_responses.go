package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type ListEventResponses struct {
	repo eventports.Repository
}

func NewListEventResponses(repo eventports.Repository) *ListEventResponses {
	return &ListEventResponses{repo: repo}
}

func (uc *ListEventResponses) Execute(ctx context.Context, eventID string, params pagination.Params) (pagination.Result[eventdomain.EventResponseWithParticipant], error) {
	if _, err := uc.repo.FindInstanceByID(ctx, eventID); err != nil {
		return pagination.Result[eventdomain.EventResponseWithParticipant]{}, err
	}
	return uc.repo.ListResponsesForEvent(ctx, eventID, params)
}
