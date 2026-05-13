package eventusecases

import (
	"context"

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
	return uc.repo.SubmitResponse(ctx, input.EventID, input.UserID, input.ProjectID, input.Answers)
}
