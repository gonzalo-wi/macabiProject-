package eventhttp

import (
	"time"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

// ── Request types ──────────────────────────────────────────────────────────────

type createEventRequest struct {
	Title              string     `json:"title" binding:"required"`
	Type               string     `json:"type"`
	StartsAt           time.Time  `json:"starts_at" binding:"required"`
	ResponseDeadlineAt *time.Time `json:"response_deadline_at"`
	Status             string     `json:"status"`
}

type patchEventRequest struct {
	Title              string     `json:"title" binding:"required"`
	Type               string     `json:"type"`
	StartsAt           time.Time  `json:"starts_at" binding:"required"`
	ResponseDeadlineAt *time.Time `json:"response_deadline_at"`
	Status             string     `json:"status"`
}

type setProjectsRequest struct {
	ProjectIDs []string `json:"project_ids" binding:"required"`
}

type createModuleRequest struct {
	EventInstanceID string `json:"event_instance_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Type            string `json:"type" binding:"required"`
	SortOrder       int    `json:"sort_order"`
	IsRequired      bool   `json:"is_required"`
}

type patchModuleRequest struct {
	Title      string `json:"title" binding:"required"`
	Type       string `json:"type" binding:"required"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type setModuleProjectsRequest struct {
	ProjectIDs []string `json:"project_ids" binding:"required"`
}

type createOptionGroupRequest struct {
	ModuleID   string `json:"module_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type patchOptionGroupRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type createOptionRequest struct {
	GroupID     string `json:"group_id" binding:"required"`
	Label       string `json:"label" binding:"required"`
	MaxCapacity *int   `json:"max_capacity"`
	SortOrder   int    `json:"sort_order"`
}

type patchOptionRequest struct {
	Label        string `json:"label" binding:"required"`
	MaxCapacity  *int   `json:"max_capacity"`
	SortOrder    int    `json:"sort_order"`
	CurrentCount *int   `json:"current_count"`
}

type submitResponseRequest struct {
	ProjectID *string                   `json:"project_id"`
	Answers   []eventdomain.AnswerInput `json:"answers" binding:"required"`
}

// ── Input parsers ──────────────────────────────────────────────────────────────

func parseEventType(s string) (eventdomain.EventType, error) {
	switch s {
	case "", string(eventdomain.EventTypeActivity):
		return eventdomain.EventTypeActivity, nil
	case string(eventdomain.EventTypeCustom):
		return eventdomain.EventTypeCustom, nil
	default:
		return "", eventdomain.ErrInvalidEventType
	}
}

func parseEventStatus(s string) (eventdomain.EventStatus, error) {
	switch s {
	case "", string(eventdomain.EventStatusDraft):
		return eventdomain.EventStatusDraft, nil
	case string(eventdomain.EventStatusOpen):
		return eventdomain.EventStatusOpen, nil
	case string(eventdomain.EventStatusClosed):
		return eventdomain.EventStatusClosed, nil
	case string(eventdomain.EventStatusCancelled):
		return eventdomain.EventStatusCancelled, nil
	default:
		return "", eventdomain.ErrInvalidEventStatus
	}
}

func parseModuleType(s string) (eventdomain.ModuleType, error) {
	switch eventdomain.ModuleType(s) {
	case eventdomain.ModuleAttendance, eventdomain.ModuleMeal, eventdomain.ModuleTransport,
		eventdomain.ModuleMaterials, eventdomain.ModuleCustom:
		return eventdomain.ModuleType(s), nil
	default:
		return "", eventdomain.ErrInvalidModuleType
	}
}

func parseOptionGroupType(s string) (eventdomain.OptionGroupType, error) {
	switch eventdomain.OptionGroupType(s) {
	case eventdomain.OptionGroupSingleChoice, eventdomain.OptionGroupMultipleChoice,
		eventdomain.OptionGroupText, eventdomain.OptionGroupNumber:
		return eventdomain.OptionGroupType(s), nil
	default:
		return "", eventdomain.ErrInvalidOptionGroupType
	}
}

// ── Response types ─────────────────────────────────────────────────────────────

type eventInstanceJSON struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Type               string     `json:"type"`
	StartsAt           time.Time  `json:"starts_at"`
	ResponseDeadlineAt *time.Time `json:"response_deadline_at,omitempty"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
}

type moduleResponse struct {
	ID              string    `json:"id"`
	EventInstanceID string    `json:"event_instance_id"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	SortOrder       int       `json:"sort_order"`
	IsRequired      bool      `json:"is_required"`
	CreatedAt       time.Time `json:"created_at"`
}

type optionGroupResponse struct {
	ID         string `json:"id"`
	ModuleID   string `json:"module_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type optionResponse struct {
	ID           string `json:"id"`
	GroupID      string `json:"group_id"`
	Label        string `json:"label"`
	MaxCapacity  *int   `json:"max_capacity,omitempty"`
	CurrentCount int    `json:"current_count"`
	SortOrder    int    `json:"sort_order"`
}

type eventResponseJSON struct {
	ID              string    `json:"id"`
	EventInstanceID string    `json:"event_instance_id"`
	UserID          string    `json:"user_id"`
	ProjectID       *string   `json:"project_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type answerResponse struct {
	ID         string  `json:"id"`
	ResponseID string  `json:"response_id"`
	GroupID    *string `json:"group_id,omitempty"`
	OptionID   *string `json:"option_id,omitempty"`
	TextValue  *string `json:"text_value,omitempty"`
}

type optionGroupDetailResponse struct {
	Group   optionGroupResponse `json:"group"`
	Options []optionResponse    `json:"options"`
}

type moduleDetailResponse struct {
	Module       moduleResponse              `json:"module"`
	ProjectIDs   []string                    `json:"project_ids"`
	OptionGroups []optionGroupDetailResponse `json:"option_groups"`
}

type eventDetailResponse struct {
	Instance   eventInstanceJSON      `json:"instance"`
	ProjectIDs []string               `json:"project_ids"`
	Modules    []moduleDetailResponse `json:"modules"`
}

type myResponseResult struct {
	Response *eventResponseJSON `json:"response"`
	Answers  []answerResponse   `json:"answers"`
}

type adminEventResponseJSON struct {
	eventResponseJSON
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

type adminResponseRow struct {
	Response adminEventResponseJSON `json:"response"`
	Answers  []answerResponse       `json:"answers"`
}

type adminResponsesResult struct {
	Data []adminResponseRow `json:"data"`
}

// ── Response mappers ───────────────────────────────────────────────────────────

func toEventInstanceJSON(e *eventdomain.EventInstance) eventInstanceJSON {
	return eventInstanceJSON{
		ID:                 e.ID,
		Title:              e.Title,
		Type:               string(e.Type),
		StartsAt:           e.StartsAt,
		ResponseDeadlineAt: e.ResponseDeadlineAt,
		Status:             string(e.Status),
		CreatedAt:          e.CreatedAt,
	}
}

func eventListResponse(res pagination.Result[eventdomain.EventInstance]) pagination.Result[eventInstanceJSON] {
	items := make([]eventInstanceJSON, len(res.Data))
	for i, e := range res.Data {
		items[i] = toEventInstanceJSON(&e)
	}
	return pagination.Result[eventInstanceJSON]{
		Data:       items,
		Total:      res.Total,
		Page:       res.Page,
		PageSize:   res.PageSize,
		TotalPages: res.TotalPages,
	}
}

func toModuleResponse(m *eventdomain.EventModule) moduleResponse {
	return moduleResponse{
		ID:              m.ID,
		EventInstanceID: m.EventInstanceID,
		Title:           m.Title,
		Type:            string(m.Type),
		SortOrder:       m.SortOrder,
		IsRequired:      m.IsRequired,
		CreatedAt:       m.CreatedAt,
	}
}

func toOptionGroupResponse(g *eventdomain.EventOptionGroup) optionGroupResponse {
	return optionGroupResponse{
		ID:         g.ID,
		ModuleID:   g.ModuleID,
		Name:       g.Name,
		Type:       string(g.Type),
		SortOrder:  g.SortOrder,
		IsRequired: g.IsRequired,
	}
}

func toOptionResponse(o *eventdomain.EventOption) optionResponse {
	return optionResponse{
		ID:           o.ID,
		GroupID:      o.GroupID,
		Label:        o.Label,
		MaxCapacity:  o.MaxCapacity,
		CurrentCount: o.CurrentCount,
		SortOrder:    o.SortOrder,
	}
}

func toEventResponseJSON(r *eventdomain.EventResponse) eventResponseJSON {
	return eventResponseJSON{
		ID:              r.ID,
		EventInstanceID: r.EventInstanceID,
		UserID:          r.UserID,
		ProjectID:       r.ProjectID,
		CreatedAt:       r.CreatedAt,
	}
}

func toAnswersResponse(aa []eventdomain.EventResponseAnswer) []answerResponse {
	out := make([]answerResponse, len(aa))
	for i, a := range aa {
		out[i] = answerResponse{
			ID:         a.ID,
			ResponseID: a.ResponseID,
			GroupID:    a.GroupID,
			OptionID:   a.OptionID,
			TextValue:  a.TextValue,
		}
	}
	return out
}

func toEventDetailResponse(d *eventdomain.EventDetail) eventDetailResponse {
	mods := make([]moduleDetailResponse, len(d.Modules))
	for i, md := range d.Modules {
		groups := make([]optionGroupDetailResponse, len(md.OptionGroups))
		for j, gd := range md.OptionGroups {
			opts := make([]optionResponse, len(gd.Options))
			for k, o := range gd.Options {
				opts[k] = toOptionResponse(&o)
			}
			groups[j] = optionGroupDetailResponse{
				Group:   toOptionGroupResponse(&gd.Group),
				Options: opts,
			}
		}
		mods[i] = moduleDetailResponse{
			Module:       toModuleResponse(&md.Module),
			ProjectIDs:   md.ProjectIDs,
			OptionGroups: groups,
		}
	}
	return eventDetailResponse{
		Instance:   toEventInstanceJSON(&d.Instance),
		ProjectIDs: d.ProjectIDs,
		Modules:    mods,
	}
}

func toMyResponseResult(r *eventdomain.EventResponse, aa []eventdomain.EventResponseAnswer) myResponseResult {
	if r == nil {
		return myResponseResult{Response: nil, Answers: []answerResponse{}}
	}
	resp := toEventResponseJSON(r)
	return myResponseResult{Response: &resp, Answers: toAnswersResponse(aa)}
}

func toAdminResponseRows(list []eventdomain.EventResponseWithParticipant) []adminResponseRow {
	rows := make([]adminResponseRow, len(list))
	for i, row := range list {
		rows[i] = adminResponseRow{
			Response: adminEventResponseJSON{
				eventResponseJSON: toEventResponseJSON(&row.Response),
				UserName:          row.UserName,
				UserEmail:         row.UserEmail,
			},
			Answers: toAnswersResponse(row.Answers),
		}
	}
	return rows
}
