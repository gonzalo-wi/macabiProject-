package eventpersistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	eventdomain "macabi-back/internal/event/domain"
)

type EventOptionGroupModel struct {
	ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ModuleID   string `gorm:"type:uuid;not null"`
	Name       string `gorm:"not null"`
	Type       string `gorm:"not null"`
	IsRequired bool   `gorm:"default:false"`
	SortOrder  int    `gorm:"default:0"`
}

func (EventOptionGroupModel) TableName() string { return "event_option_groups" }

type EventOptionModel struct {
	ID           string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	GroupID      string `gorm:"type:uuid;not null"`
	Label        string `gorm:"not null"`
	MaxCapacity  *int
	CurrentCount int `gorm:"default:0"`
	SortOrder    int `gorm:"default:0"`
}

func (EventOptionModel) TableName() string { return "event_options" }

func (r *RepositoryPG) CreateOptionGroup(ctx context.Context, g *eventdomain.EventOptionGroup) error {
	row := toGroupModel(g)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	g.ID = row.ID
	return nil
}

func (r *RepositoryPG) FindOptionGroupByID(ctx context.Context, id string) (*eventdomain.EventOptionGroup, error) {
	var m EventOptionGroupModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, eventdomain.ErrNotFound
		}
		return nil, err
	}
	return toDomainGroup(&m), nil
}

func (r *RepositoryPG) UpdateOptionGroup(ctx context.Context, g *eventdomain.EventOptionGroup) error {
	return r.db.WithContext(ctx).Model(&EventOptionGroupModel{}).Where("id = ?", g.ID).Updates(map[string]interface{}{
		"name":        g.Name,
		"type":        string(g.Type),
		"sort_order":  g.SortOrder,
		"is_required": g.IsRequired,
	}).Error
}

func (r *RepositoryPG) DeleteOptionGroup(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&EventOptionGroupModel{}, "id = ?", id).Error
}

func (r *RepositoryPG) CreateOption(ctx context.Context, o *eventdomain.EventOption) error {
	row := toOptModel(o)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	o.ID = row.ID
	return nil
}

func (r *RepositoryPG) UpdateOption(ctx context.Context, o *eventdomain.EventOption) error {
	return r.db.WithContext(ctx).Model(&EventOptionModel{}).Where("id = ?", o.ID).Updates(map[string]interface{}{
		"label":         o.Label,
		"max_capacity":  o.MaxCapacity,
		"sort_order":    o.SortOrder,
		"current_count": o.CurrentCount,
	}).Error
}

func (r *RepositoryPG) DeleteOption(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&EventOptionModel{}, "id = ?", id).Error
}

func (r *RepositoryPG) FindOptionByID(ctx context.Context, id string) (*eventdomain.EventOption, error) {
	var row EventOptionModel
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, eventdomain.ErrNotFound
		}
		return nil, err
	}
	return toDomainOpt(&row), nil
}

func toGroupModel(g *eventdomain.EventOptionGroup) EventOptionGroupModel {
	return EventOptionGroupModel{
		ID:         g.ID,
		ModuleID:   g.ModuleID,
		Name:       g.Name,
		Type:       string(g.Type),
		IsRequired: g.IsRequired,
		SortOrder:  g.SortOrder,
	}
}

func toDomainGroup(m *EventOptionGroupModel) *eventdomain.EventOptionGroup {
	return &eventdomain.EventOptionGroup{
		ID:         m.ID,
		ModuleID:   m.ModuleID,
		Name:       m.Name,
		Type:       eventdomain.OptionGroupType(m.Type),
		IsRequired: m.IsRequired,
		SortOrder:  m.SortOrder,
	}
}

func toOptModel(o *eventdomain.EventOption) EventOptionModel {
	return EventOptionModel{
		ID:           o.ID,
		GroupID:      o.GroupID,
		Label:        o.Label,
		MaxCapacity:  o.MaxCapacity,
		CurrentCount: o.CurrentCount,
		SortOrder:    o.SortOrder,
	}
}

func toDomainOpt(m *EventOptionModel) *eventdomain.EventOption {
	return &eventdomain.EventOption{
		ID:           m.ID,
		GroupID:      m.GroupID,
		Label:        m.Label,
		MaxCapacity:  m.MaxCapacity,
		CurrentCount: m.CurrentCount,
		SortOrder:    m.SortOrder,
	}
}
