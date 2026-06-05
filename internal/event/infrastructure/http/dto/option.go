package eventdto

import (
	eventdomain "macabi-back/internal/event/domain"
)

type CreateOptionGroupRequest struct {
	ModuleID   string `json:"module_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type PatchOptionGroupRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type CreateOptionRequest struct {
	GroupID     string `json:"group_id" binding:"required"`
	Label       string `json:"label" binding:"required"`
	MaxCapacity *int   `json:"max_capacity"`
	SortOrder   int    `json:"sort_order"`
}

type PatchOptionRequest struct {
	Label        string `json:"label" binding:"required"`
	MaxCapacity  *int   `json:"max_capacity"`
	SortOrder    int    `json:"sort_order"`
	CurrentCount *int   `json:"current_count"`
}

type OptionGroupResponse struct {
	ID         string `json:"id"`
	ModuleID   string `json:"module_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	SortOrder  int    `json:"sort_order"`
	IsRequired bool   `json:"is_required"`
}

type OptionResponse struct {
	ID           string `json:"id"`
	GroupID      string `json:"group_id"`
	Label        string `json:"label"`
	MaxCapacity  *int   `json:"max_capacity,omitempty"`
	CurrentCount int    `json:"current_count"`
	SortOrder    int    `json:"sort_order"`
}

type OptionGroupDetailResponse struct {
	Group   OptionGroupResponse `json:"group"`
	Options []OptionResponse    `json:"options"`
}

func ParseOptionGroupType(s string) (eventdomain.OptionGroupType, error) {
	switch eventdomain.OptionGroupType(s) {
	case eventdomain.OptionGroupSingleChoice, eventdomain.OptionGroupMultipleChoice,
		eventdomain.OptionGroupText, eventdomain.OptionGroupNumber:
		return eventdomain.OptionGroupType(s), nil
	default:
		return "", eventdomain.ErrInvalidOptionGroupType
	}
}

func ToOptionGroupResponse(g *eventdomain.EventOptionGroup) OptionGroupResponse {
	return OptionGroupResponse{
		ID:         g.ID,
		ModuleID:   g.ModuleID,
		Name:       g.Name,
		Type:       string(g.Type),
		SortOrder:  g.SortOrder,
		IsRequired: g.IsRequired,
	}
}

func ToOptionResponse(o *eventdomain.EventOption) OptionResponse {
	return OptionResponse{
		ID:           o.ID,
		GroupID:      o.GroupID,
		Label:        o.Label,
		MaxCapacity:  o.MaxCapacity,
		CurrentCount: o.CurrentCount,
		SortOrder:    o.SortOrder,
	}
}
