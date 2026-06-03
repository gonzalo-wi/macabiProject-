package eventusecases

import (
	"context"
	"time"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
)

type SubmitResponseInput struct {
	EventID   string
	UserID    string
	ProjectID *string
	Answers   []eventdomain.AnswerInput
}

type SubmitResponse struct {
	repo eventports.Repository
}

func NewSubmitResponse(repo eventports.Repository) *SubmitResponse {
	return &SubmitResponse{repo: repo}
}

func (uc *SubmitResponse) Execute(ctx context.Context, input SubmitResponseInput) error {
	inst, err := uc.repo.FindInstanceByID(ctx, input.EventID)
	if err != nil {
		return err
	}
	if inst.Status != eventdomain.EventStatusOpen {
		return eventdomain.ErrEventNotOpen
	}
	if inst.ResponseDeadlineAt != nil && inst.ResponseDeadlineAt.Before(time.Now()) {
		return eventdomain.ErrResponseDeadlinePassed
	}
	return uc.repo.SubmitResponse(ctx, input.EventID, input.UserID, input.ProjectID, input.Answers)
}
