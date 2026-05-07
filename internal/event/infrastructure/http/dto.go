package eventhttp

import (
	"time"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

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

type eventInstanceJSON struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Type               string     `json:"type"`
	StartsAt           time.Time  `json:"starts_at"`
	ResponseDeadlineAt *time.Time `json:"response_deadline_at,omitempty"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
}

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
