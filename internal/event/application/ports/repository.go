package eventports

import (
	"context"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type Repository interface {
	CreateInstance(ctx context.Context, e *eventdomain.EventInstance) error
	UpdateInstance(ctx context.Context, e *eventdomain.EventInstance) error
	DeleteInstance(ctx context.Context, id string) error
	FindInstanceByID(ctx context.Context, id string) (*eventdomain.EventInstance, error)
	ListInstances(ctx context.Context, params pagination.Params) (pagination.Result[eventdomain.EventInstance], error)
	LoadEventDetail(ctx context.Context, id string) (*eventdomain.EventDetail, error)
	SetInstanceProjects(ctx context.Context, eventID string, projectIDs []string) error

	CreateModule(ctx context.Context, m *eventdomain.EventModule) error
	UpdateModule(ctx context.Context, m *eventdomain.EventModule) error
	DeleteModule(ctx context.Context, id string) error
	SetModuleProjects(ctx context.Context, moduleID string, projectIDs []string) error

	FindModuleByID(ctx context.Context, id string) (*eventdomain.EventModule, error)
	CreateOptionGroup(ctx context.Context, g *eventdomain.EventOptionGroup) error
	FindOptionGroupByID(ctx context.Context, id string) (*eventdomain.EventOptionGroup, error)
	UpdateOptionGroup(ctx context.Context, g *eventdomain.EventOptionGroup) error
	DeleteOptionGroup(ctx context.Context, id string) error
	CreateOption(ctx context.Context, o *eventdomain.EventOption) error
	UpdateOption(ctx context.Context, o *eventdomain.EventOption) error
	FindOptionByID(ctx context.Context, id string) (*eventdomain.EventOption, error)
	DeleteOption(ctx context.Context, id string) error

	FindResponseByEventAndUser(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, error)
	SaveResponse(ctx context.Context, r *eventdomain.EventResponse) error
	LoadUserResponseDetail(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, []eventdomain.EventResponseAnswer, error)
	SubmitResponse(ctx context.Context, eventID, userID string, projectID *string, answers []eventdomain.AnswerInput) error
	ListResponsesForEvent(ctx context.Context, eventID string) ([]eventdomain.EventResponseWithParticipant, error)
}
