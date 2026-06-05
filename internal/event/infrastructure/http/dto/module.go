package eventdto

import (
	"time"

	eventdomain "macabi-back/internal/event/domain"
)

type CreateModuleRequest struct {
	EventInstanceID string `json:"event_instance_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Type            string `json:"type" binding:"required"`
	SortOrder       int    `json:"sort_order"`
	IsRequired      bool   `json:"is_required"`
}

type PatchModuleRequest struct {
	Title      string `json:"title" binding:"required"`
	Type       string `json:"type" binding:"required"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type SetModuleProjectsRequest struct {
	ProjectIDs []string `json:"project_ids" binding:"required"`
}

type ModuleResponse struct {
	ID              string    `json:"id"`
	EventInstanceID string    `json:"event_instance_id"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	SortOrder       int       `json:"sort_order"`
	IsRequired      bool      `json:"is_required"`
	CreatedAt       time.Time `json:"created_at"`
}

type ModuleDetailResponse struct {
	Module       ModuleResponse              `json:"module"`
	ProjectIDs   []string                    `json:"project_ids"`
	OptionGroups []OptionGroupDetailResponse `json:"option_groups"`
}

func ParseModuleType(s string) (eventdomain.ModuleType, error) {
	switch eventdomain.ModuleType(s) {
	case eventdomain.ModuleAttendance, eventdomain.ModuleMeal, eventdomain.ModuleTransport,
		eventdomain.ModuleMaterials, eventdomain.ModuleCustom:
		return eventdomain.ModuleType(s), nil
	default:
		return "", eventdomain.ErrInvalidModuleType
	}
}

func ToModuleResponse(m *eventdomain.EventModule) ModuleResponse {
	return ModuleResponse{
		ID:              m.ID,
		EventInstanceID: m.EventInstanceID,
		Title:           m.Title,
		Type:            string(m.Type),
		SortOrder:       m.SortOrder,
		IsRequired:      m.IsRequired,
		CreatedAt:       m.CreatedAt,
	}
}
