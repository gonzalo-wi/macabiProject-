package eventdto

import (
	"time"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type CreateEventRequest struct {
	Title              string     `json:"title" binding:"required"`
	Type               string     `json:"type"`
	StartsAt           time.Time  `json:"starts_at" binding:"required"`
	ResponseDeadlineAt *time.Time `json:"response_deadline_at"`
	Status             string     `json:"status"`
}

type PatchEventRequest struct {
	Title              string     `json:"title" binding:"required"`
	Type               string     `json:"type"`
	StartsAt           time.Time  `json:"starts_at" binding:"required"`
	ResponseDeadlineAt *time.Time `json:"response_deadline_at"`
	Status             string     `json:"status"`
}

type SetProjectsRequest struct {
	ProjectIDs []string `json:"project_ids" binding:"required"`
}

type EventInstanceJSON struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Type               string     `json:"type"`
	StartsAt           time.Time  `json:"starts_at"`
	ResponseDeadlineAt *time.Time `json:"response_deadline_at,omitempty"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
}

type EventDetailResponse struct {
	Instance   EventInstanceJSON      `json:"instance"`
	ProjectIDs []string               `json:"project_ids"`
	Modules    []ModuleDetailResponse `json:"modules"`
}

func ParseEventType(s string) (eventdomain.EventType, error) {
	switch s {
	case "", string(eventdomain.EventTypeActivity):
		return eventdomain.EventTypeActivity, nil
	case string(eventdomain.EventTypeCustom):
		return eventdomain.EventTypeCustom, nil
	default:
		return "", eventdomain.ErrInvalidEventType
	}
}

func ParseEventStatus(s string) (eventdomain.EventStatus, error) {
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

func ToEventInstanceJSON(e *eventdomain.EventInstance) EventInstanceJSON {
	return EventInstanceJSON{
		ID:                 e.ID,
		Title:              e.Title,
		Type:               string(e.Type),
		StartsAt:           e.StartsAt,
		ResponseDeadlineAt: e.ResponseDeadlineAt,
		Status:             string(e.Status),
		CreatedAt:          e.CreatedAt,
	}
}

func EventListResponse(res pagination.Result[eventdomain.EventInstance]) pagination.Result[EventInstanceJSON] {
	items := make([]EventInstanceJSON, len(res.Data))
	for i, e := range res.Data {
		items[i] = ToEventInstanceJSON(&e)
	}
	return pagination.Result[EventInstanceJSON]{
		Data:       items,
		Total:      res.Total,
		Page:       res.Page,
		PageSize:   res.PageSize,
		TotalPages: res.TotalPages,
	}
}

func ToEventDetailResponse(d *eventdomain.EventDetail) EventDetailResponse {
	mods := make([]ModuleDetailResponse, len(d.Modules))
	for i, md := range d.Modules {
		groups := make([]OptionGroupDetailResponse, len(md.OptionGroups))
		for j, gd := range md.OptionGroups {
			opts := make([]OptionResponse, len(gd.Options))
			for k, o := range gd.Options {
				opts[k] = ToOptionResponse(&o)
			}
			groups[j] = OptionGroupDetailResponse{
				Group:   ToOptionGroupResponse(&gd.Group),
				Options: opts,
			}
		}
		mods[i] = ModuleDetailResponse{
			Module:       ToModuleResponse(&md.Module),
			ProjectIDs:   md.ProjectIDs,
			OptionGroups: groups,
		}
	}
	return EventDetailResponse{
		Instance:   ToEventInstanceJSON(&d.Instance),
		ProjectIDs: d.ProjectIDs,
		Modules:    mods,
	}
}
