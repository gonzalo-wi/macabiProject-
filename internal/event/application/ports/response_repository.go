package eventports

import (
	"context"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type ResponseRepository interface {
	FindResponseByEventAndUser(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, error)
	SaveResponse(ctx context.Context, r *eventdomain.EventResponse) error
	LoadUserResponseDetail(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, []eventdomain.EventResponseAnswer, error)
	SubmitResponse(ctx context.Context, eventID, userID string, projectID *string, answers []eventdomain.AnswerInput) error
	ListResponsesForEvent(ctx context.Context, eventID string, params pagination.Params) (pagination.Result[eventdomain.EventResponseWithParticipant], error)
	LoadModuleResponseSummary(ctx context.Context, moduleID string) (*eventdomain.ModuleResponseSummary, error)
}
