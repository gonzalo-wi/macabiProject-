package eventusecases

import (
	"context"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type GetMyResponse struct {
	repo eventports.ResponseRepository
}

func NewGetMyResponse(repo eventports.ResponseRepository) *GetMyResponse {
	return &GetMyResponse{repo: repo}
}

func (uc *GetMyResponse) Execute(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, []eventdomain.EventResponseAnswer, error) {
	return uc.repo.LoadUserResponseDetail(ctx, eventID, userID)
}
