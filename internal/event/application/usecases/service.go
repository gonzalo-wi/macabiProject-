package eventusecases

import (
	"context"
	"strings"
	"time"

	eventports "macabi-back/internal/event/application/ports"
	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type Service struct {
	repo eventports.Repository
}

func NewService(repo eventports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEvent(ctx context.Context, title string, typ eventdomain.EventType, startsAt time.Time, deadline *time.Time, status eventdomain.EventStatus) (*eventdomain.EventInstance, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	e := &eventdomain.EventInstance{
		Title:              title,
		Type:               typ,
		StartsAt:           startsAt,
		ResponseDeadlineAt: deadline,
		Status:             status,
	}
	if err := s.repo.CreateInstance(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) ListEvents(ctx context.Context, params pagination.Params) (pagination.Result[eventdomain.EventInstance], error) {
	return s.repo.ListInstances(ctx, params)
}

func (s *Service) GetEventDetail(ctx context.Context, id string) (*eventdomain.EventDetail, error) {
	return s.repo.LoadEventDetail(ctx, id)
}

func (s *Service) UpdateEvent(ctx context.Context, id, title string, typ eventdomain.EventType, startsAt time.Time, deadline *time.Time, status eventdomain.EventStatus) (*eventdomain.EventInstance, error) {
	e, err := s.repo.FindInstanceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	e.Title = title
	e.Type = typ
	e.StartsAt = startsAt
	e.ResponseDeadlineAt = deadline
	e.Status = status
	if err := s.repo.UpdateInstance(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) SetEventProjects(ctx context.Context, eventID string, projectIDs []string) error {
	return s.repo.SetInstanceProjects(ctx, eventID, projectIDs)
}

func (s *Service) CreateModule(ctx context.Context, eventID, title string, typ eventdomain.ModuleType, sort int, required bool) (*eventdomain.EventModule, error) {
	m := &eventdomain.EventModule{
		EventInstanceID: eventID,
		Title:           strings.TrimSpace(title),
		Type:            typ,
		SortOrder:       sort,
		IsRequired:      required,
	}
	if m.Title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	if err := s.repo.CreateModule(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) UpdateModule(ctx context.Context, id, title string, typ eventdomain.ModuleType, sort int, required bool) (*eventdomain.EventModule, error) {
	m, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	m.Title = title
	m.Type = typ
	m.SortOrder = sort
	m.IsRequired = required
	if err := s.repo.UpdateModule(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) DeleteModule(ctx context.Context, id string) error {
	return s.repo.DeleteModule(ctx, id)
}

func (s *Service) SetModuleProjects(ctx context.Context, moduleID string, projectIDs []string) error {
	return s.repo.SetModuleProjects(ctx, moduleID, projectIDs)
}

func (s *Service) CreateOptionGroup(ctx context.Context, moduleID, name string, typ eventdomain.OptionGroupType, sort int, required bool) (*eventdomain.EventOptionGroup, error) {
	g := &eventdomain.EventOptionGroup{
		ModuleID:   moduleID,
		Name:       strings.TrimSpace(name),
		Type:       typ,
		SortOrder:  sort,
		IsRequired: required,
	}
	if g.Name == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	if err := s.repo.CreateOptionGroup(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *Service) UpdateOptionGroup(ctx context.Context, id, name string, typ eventdomain.OptionGroupType, sort int, required bool) (*eventdomain.EventOptionGroup, error) {
	g, err := s.repo.FindOptionGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	g.Name = name
	g.Type = typ
	g.SortOrder = sort
	g.IsRequired = required
	if err := s.repo.UpdateOptionGroup(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *Service) CreateOption(ctx context.Context, groupID, label string, maxCap *int, sort int) (*eventdomain.EventOption, error) {
	o := &eventdomain.EventOption{
		GroupID:      groupID,
		Label:        strings.TrimSpace(label),
		MaxCapacity:  maxCap,
		CurrentCount: 0,
		SortOrder:    sort,
	}
	if o.Label == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	if err := s.repo.CreateOption(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) UpdateOption(ctx context.Context, id, label string, maxCap *int, sort int, currentCount *int) (*eventdomain.EventOption, error) {
	o, err := s.repo.FindOptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, eventdomain.ErrEmptyTitle
	}
	o.Label = label
	o.MaxCapacity = maxCap
	o.SortOrder = sort
	if currentCount != nil {
		o.CurrentCount = *currentCount
	}
	if err := s.repo.UpdateOption(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) SubmitResponse(ctx context.Context, eventID, userID string, projectID *string, answers []eventdomain.AnswerInput) error {
	return s.repo.SubmitResponse(ctx, eventID, userID, projectID, answers)
}

func (s *Service) GetMyResponse(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, []eventdomain.EventResponseAnswer, error) {
	return s.repo.LoadUserResponseDetail(ctx, eventID, userID)
}
