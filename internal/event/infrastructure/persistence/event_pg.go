package eventpersistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type EventInstanceModel struct {
	ID                 string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title              string `gorm:"not null"`
	Type               string `gorm:"not null;default:'activity'"`
	StartsAt           time.Time
	ResponseDeadlineAt *time.Time
	Status             string `gorm:"default:'draft'"`
	CreatedAt          time.Time
}

func (EventInstanceModel) TableName() string { return "event_instances" }

type EventInstanceProjectModel struct {
	ID              string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EventInstanceID string `gorm:"type:uuid;not null"`
	ProjectID       string `gorm:"type:uuid;not null"`
}

func (EventInstanceProjectModel) TableName() string { return "event_instance_projects" }

func (r *RepositoryPG) CreateInstance(ctx context.Context, e *eventdomain.EventInstance) error {
	m := toInstModel(e)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	e.ID = m.ID
	e.CreatedAt = m.CreatedAt
	return nil
}

func (r *RepositoryPG) UpdateInstance(ctx context.Context, e *eventdomain.EventInstance) error {
	return r.db.WithContext(ctx).Model(&EventInstanceModel{}).Where("id = ?", e.ID).Updates(map[string]interface{}{
		"title":                e.Title,
		"type":                 string(e.Type),
		"starts_at":            e.StartsAt,
		"response_deadline_at": e.ResponseDeadlineAt,
		"status":               string(e.Status),
	}).Error
}

func (r *RepositoryPG) DeleteInstance(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var respIDs []string
		if err := tx.Model(&EventResponseModel{}).Where("event_instance_id = ?", id).Pluck("id", &respIDs).Error; err != nil {
			return err
		}
		if len(respIDs) > 0 {
			if err := tx.Where("response_id IN ?", respIDs).Delete(&EventResponseAnswerModel{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("event_instance_id = ?", id).Delete(&EventResponseModel{}).Error; err != nil {
			return err
		}

		var moduleIDs []string
		if err := tx.Model(&EventModuleModel{}).Where("event_instance_id = ?", id).Pluck("id", &moduleIDs).Error; err != nil {
			return err
		}
		for _, mid := range moduleIDs {
			if err := tx.Where("module_id = ?", mid).Delete(&EventModuleProjectModel{}).Error; err != nil {
				return err
			}
		}
		if len(moduleIDs) > 0 {
			var groupIDs []string
			if err := tx.Model(&EventOptionGroupModel{}).Where("module_id IN ?", moduleIDs).Pluck("id", &groupIDs).Error; err != nil {
				return err
			}
			if len(groupIDs) > 0 {
				if err := tx.Where("group_id IN ?", groupIDs).Delete(&EventOptionModel{}).Error; err != nil {
					return err
				}
			}
			if err := tx.Where("module_id IN ?", moduleIDs).Delete(&EventOptionGroupModel{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("event_instance_id = ?", id).Delete(&EventModuleModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("event_instance_id = ?", id).Delete(&EventInstanceProjectModel{}).Error; err != nil {
			return err
		}

		res := tx.Where("id = ?", id).Delete(&EventInstanceModel{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return eventdomain.ErrNotFound
		}
		return nil
	})
}

func (r *RepositoryPG) FindInstanceByID(ctx context.Context, id string) (*eventdomain.EventInstance, error) {
	var m EventInstanceModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, eventdomain.ErrNotFound
		}
		return nil, err
	}
	return toDomainInst(&m), nil
}

func (r *RepositoryPG) ListInstances(ctx context.Context, params pagination.Params) (pagination.Result[eventdomain.EventInstance], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&EventInstanceModel{}).Count(&total).Error; err != nil {
		return pagination.Result[eventdomain.EventInstance]{}, err
	}
	var rows []EventInstanceModel
	if err := r.db.WithContext(ctx).Order("starts_at DESC").
		Offset(params.Offset()).Limit(params.PageSize).Find(&rows).Error; err != nil {
		return pagination.Result[eventdomain.EventInstance]{}, err
	}
	items := make([]eventdomain.EventInstance, len(rows))
	for i := range rows {
		items[i] = *toDomainInst(&rows[i])
	}
	return pagination.NewResult(items, total, params), nil
}

func (r *RepositoryPG) SetInstanceProjects(ctx context.Context, eventID string, projectIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("event_instance_id = ?", eventID).Delete(&EventInstanceProjectModel{}).Error; err != nil {
			return err
		}
		for _, pid := range projectIDs {
			row := EventInstanceProjectModel{EventInstanceID: eventID, ProjectID: pid}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RepositoryPG) LoadEventDetail(ctx context.Context, id string) (*eventdomain.EventDetail, error) {
	var inst EventInstanceModel
	if err := r.db.WithContext(ctx).First(&inst, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, eventdomain.ErrNotFound
		}
		return nil, err
	}
	var ipRows []EventInstanceProjectModel
	if err := r.db.WithContext(ctx).Where("event_instance_id = ?", id).Find(&ipRows).Error; err != nil {
		return nil, err
	}
	pids := make([]string, len(ipRows))
	for i, row := range ipRows {
		pids[i] = row.ProjectID
	}
	var modRows []EventModuleModel
	if err := r.db.WithContext(ctx).Where("event_instance_id = ?", id).
		Order("sort_order ASC, created_at ASC").Find(&modRows).Error; err != nil {
		return nil, err
	}
	modules := make([]eventdomain.ModuleDetail, len(modRows))
	for i, m := range modRows {
		var mpRows []EventModuleProjectModel
		_ = r.db.WithContext(ctx).Where("module_id = ?", m.ID).Find(&mpRows).Error
		mpids := make([]string, len(mpRows))
		for j, row := range mpRows {
			mpids[j] = row.ProjectID
		}
		var grpRows []EventOptionGroupModel
		if err := r.db.WithContext(ctx).Where("module_id = ?", m.ID).Order("sort_order ASC").Find(&grpRows).Error; err != nil {
			return nil, err
		}
		groups := make([]eventdomain.GroupDetail, len(grpRows))
		for j, g := range grpRows {
			var optRows []EventOptionModel
			if err := r.db.WithContext(ctx).Where("group_id = ?", g.ID).Order("sort_order ASC").Find(&optRows).Error; err != nil {
				return nil, err
			}
			opts := make([]eventdomain.EventOption, len(optRows))
			for k, o := range optRows {
				opts[k] = *toDomainOpt(&o)
			}
			groups[j] = eventdomain.GroupDetail{
				Group:   *toDomainGroup(&g),
				Options: opts,
			}
		}
		modules[i] = eventdomain.ModuleDetail{
			Module:       *toDomainMod(&m),
			ProjectIDs:   mpids,
			OptionGroups: groups,
		}
	}
	return &eventdomain.EventDetail{
		Instance:   *toDomainInst(&inst),
		ProjectIDs: pids,
		Modules:    modules,
	}, nil
}

func toInstModel(e *eventdomain.EventInstance) EventInstanceModel {
	return EventInstanceModel{
		ID:                 e.ID,
		Title:              e.Title,
		Type:               string(e.Type),
		StartsAt:           e.StartsAt,
		ResponseDeadlineAt: e.ResponseDeadlineAt,
		Status:             string(e.Status),
	}
}

func toDomainInst(m *EventInstanceModel) *eventdomain.EventInstance {
	return &eventdomain.EventInstance{
		ID:                 m.ID,
		Title:              m.Title,
		Type:               eventdomain.EventType(m.Type),
		StartsAt:           m.StartsAt,
		ResponseDeadlineAt: m.ResponseDeadlineAt,
		Status:             eventdomain.EventStatus(m.Status),
		CreatedAt:          m.CreatedAt,
	}
}
