package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type ListEventResponses struct {
	eventRepo    eventports.EventRepository
	responseRepo eventports.ResponseRepository
}

func NewListEventResponses(eventRepo eventports.EventRepository, responseRepo eventports.ResponseRepository) *ListEventResponses {
	return &ListEventResponses{eventRepo: eventRepo, responseRepo: responseRepo}
}

func (uc *ListEventResponses) Execute(ctx context.Context, eventID string, params pagination.Params) (pagination.Result[eventdomain.EventResponseWithParticipant], error) {
	if _, err := uc.eventRepo.FindInstanceByID(ctx, eventID); err != nil {
		return pagination.Result[eventdomain.EventResponseWithParticipant]{}, err
	}
	return uc.responseRepo.ListResponsesForEvent(ctx, eventID, params)
}
