package eventpersistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	eventdomain "macabi-back/internal/event/domain"
)

type EventModuleModel struct {
	ID              string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EventInstanceID string `gorm:"type:uuid;not null"`
	Title           string `gorm:"not null"`
	Type            string `gorm:"not null"`
	SortOrder       int    `gorm:"default:0"`
	IsRequired      bool   `gorm:"default:false"`
	CreatedAt       time.Time
}

func (EventModuleModel) TableName() string { return "event_modules" }

type EventModuleProjectModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ModuleID  string `gorm:"type:uuid;not null"`
	ProjectID string `gorm:"type:uuid;not null"`
}

func (EventModuleProjectModel) TableName() string { return "event_module_projects" }

func (r *RepositoryPG) CreateModule(ctx context.Context, m *eventdomain.EventModule) error {
	row := toModModel(m)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	m.ID = row.ID
	m.CreatedAt = row.CreatedAt
	return nil
}

func (r *RepositoryPG) UpdateModule(ctx context.Context, m *eventdomain.EventModule) error {
	return r.db.WithContext(ctx).Model(&EventModuleModel{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"title":       m.Title,
		"type":        string(m.Type),
		"sort_order":  m.SortOrder,
		"is_required": m.IsRequired,
	}).Error
}

func (r *RepositoryPG) DeleteModule(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&EventModuleModel{}, "id = ?", id).Error
}

func (r *RepositoryPG) SetModuleProjects(ctx context.Context, moduleID string, projectIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("module_id = ?", moduleID).Delete(&EventModuleProjectModel{}).Error; err != nil {
			return err
		}
		for _, pid := range projectIDs {
			row := EventModuleProjectModel{ModuleID: moduleID, ProjectID: pid}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RepositoryPG) FindModuleByID(ctx context.Context, id string) (*eventdomain.EventModule, error) {
	var m EventModuleModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, eventdomain.ErrNotFound
		}
		return nil, err
	}
	return toDomainMod(&m), nil
}

func toModModel(m *eventdomain.EventModule) EventModuleModel {
	return EventModuleModel{
		ID:              m.ID,
		EventInstanceID: m.EventInstanceID,
		Title:           m.Title,
		Type:            string(m.Type),
		SortOrder:       m.SortOrder,
		IsRequired:      m.IsRequired,
	}
}

func toDomainMod(m *EventModuleModel) *eventdomain.EventModule {
	return &eventdomain.EventModule{
		ID:              m.ID,
		EventInstanceID: m.EventInstanceID,
		Title:           m.Title,
		Type:            eventdomain.ModuleType(m.Type),
		SortOrder:       m.SortOrder,
		IsRequired:      m.IsRequired,
		CreatedAt:       m.CreatedAt,
	}
}
