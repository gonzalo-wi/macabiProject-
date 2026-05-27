package eventpersistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	eventports "macabi-back/internal/event/application/ports"
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

type EventResponseModel struct {
	ID              string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EventInstanceID string  `gorm:"type:uuid;not null"`
	UserID          string  `gorm:"type:uuid;not null"`
	ProjectID       *string `gorm:"type:uuid"`
	CreatedAt       time.Time
}

func (EventResponseModel) TableName() string { return "event_responses" }

type EventResponseAnswerModel struct {
	ID         string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ResponseID string  `gorm:"type:uuid;not null"`
	GroupID    *string `gorm:"type:uuid"`
	OptionID   *string `gorm:"type:uuid"`
	TextValue  *string
}

func (EventResponseAnswerModel) TableName() string { return "event_response_answers" }

type RepositoryPG struct {
	db *gorm.DB
}

func NewRepositoryPG(db *gorm.DB) *RepositoryPG {
	return &RepositoryPG{db: db}
}

var _ eventports.Repository = (*RepositoryPG)(nil)

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

func (r *RepositoryPG) DeleteOptionGroup(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&EventOptionGroupModel{}, "id = ?", id).Error
}

func (r *RepositoryPG) UpdateOptionGroup(ctx context.Context, g *eventdomain.EventOptionGroup) error {
	return r.db.WithContext(ctx).Model(&EventOptionGroupModel{}).Where("id = ?", g.ID).Updates(map[string]interface{}{
		"name":        g.Name,
		"type":        string(g.Type),
		"sort_order":  g.SortOrder,
		"is_required": g.IsRequired,
	}).Error
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

func (r *RepositoryPG) FindResponseByEventAndUser(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, error) {
	var row EventResponseModel
	err := r.db.WithContext(ctx).Where("event_instance_id = ? AND user_id = ?", eventID, userID).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainResp(&row), nil
}

func (r *RepositoryPG) SaveResponse(ctx context.Context, resp *eventdomain.EventResponse) error {
	row := EventResponseModel{
		EventInstanceID: resp.EventInstanceID,
		UserID:          resp.UserID,
		ProjectID:       resp.ProjectID,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	resp.ID = row.ID
	resp.CreatedAt = row.CreatedAt
	return nil
}

func (r *RepositoryPG) LoadUserResponseDetail(ctx context.Context, eventID, userID string) (*eventdomain.EventResponse, []eventdomain.EventResponseAnswer, error) {
	resp, err := r.FindResponseByEventAndUser(ctx, eventID, userID)
	if err != nil || resp == nil {
		return resp, nil, err
	}
	var ans []EventResponseAnswerModel
	if err := r.db.WithContext(ctx).Where("response_id = ?", resp.ID).Find(&ans).Error; err != nil {
		return nil, nil, err
	}
	out := make([]eventdomain.EventResponseAnswer, len(ans))
	for i, a := range ans {
		out[i] = *toDomainAns(&a)
	}
	return resp, out, nil
}

func (r *RepositoryPG) ListResponsesForEvent(ctx context.Context, eventID string, params pagination.Params) (pagination.Result[eventdomain.EventResponseWithParticipant], error) {
	type listRow struct {
		ID              string         `gorm:"column:id"`
		EventInstanceID string         `gorm:"column:event_instance_id"`
		UserID          string         `gorm:"column:user_id"`
		ProjectID       *string        `gorm:"column:project_id"`
		CreatedAt       time.Time      `gorm:"column:created_at"`
		UserName        sql.NullString `gorm:"column:user_name"`
		UserEmail       sql.NullString `gorm:"column:user_email"`
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&EventResponseModel{}).Where("event_instance_id = ?", eventID).Count(&total).Error; err != nil {
		return pagination.Result[eventdomain.EventResponseWithParticipant]{}, err
	}

	var rows []listRow
	err := r.db.WithContext(ctx).Raw(`
SELECT er.id, er.event_instance_id, er.user_id, er.project_id, er.created_at,
       u.name AS user_name, u.email AS user_email
FROM event_responses er
LEFT JOIN users u ON u.id = er.user_id
WHERE er.event_instance_id = ?
ORDER BY er.created_at ASC
LIMIT ? OFFSET ?`, eventID, params.PageSize, params.Offset()).Scan(&rows).Error
	if err != nil {
		return pagination.Result[eventdomain.EventResponseWithParticipant]{}, err
	}
	if len(rows) == 0 {
		return pagination.NewResult([]eventdomain.EventResponseWithParticipant{}, total, params), nil
	}
	respIDs := make([]string, len(rows))
	for i := range rows {
		respIDs[i] = rows[i].ID
	}
	var ansRows []EventResponseAnswerModel
	if err := r.db.WithContext(ctx).Where("response_id IN ?", respIDs).Find(&ansRows).Error; err != nil {
		return pagination.Result[eventdomain.EventResponseWithParticipant]{}, err
	}
	byResp := make(map[string][]eventdomain.EventResponseAnswer)
	for i := range ansRows {
		da := *toDomainAns(&ansRows[i])
		byResp[ansRows[i].ResponseID] = append(byResp[ansRows[i].ResponseID], da)
	}
	out := make([]eventdomain.EventResponseWithParticipant, len(rows))
	for i, row := range rows {
		m := EventResponseModel{
			ID: row.ID, EventInstanceID: row.EventInstanceID, UserID: row.UserID,
			ProjectID: row.ProjectID, CreatedAt: row.CreatedAt,
		}
		uname, uemail := "", ""
		if row.UserName.Valid {
			uname = row.UserName.String
		}
		if row.UserEmail.Valid {
			uemail = row.UserEmail.String
		}
		resp := toDomainResp(&m)
		out[i] = eventdomain.EventResponseWithParticipant{
			Response:  *resp,
			UserName:  uname,
			UserEmail: uemail,
			Answers:   byResp[resp.ID],
		}
	}
	return pagination.NewResult(out, total, params), nil
}

func (r *RepositoryPG) SubmitResponse(ctx context.Context, eventID, userID string, projectID *string, answers []eventdomain.AnswerInput) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var resp EventResponseModel
		err := tx.Where("event_instance_id = ? AND user_id = ?", eventID, userID).First(&resp).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp = EventResponseModel{
				EventInstanceID: eventID,
				UserID:          userID,
				ProjectID:       projectID,
			}
			if err := tx.Create(&resp).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if projectID != nil {
				if err := tx.Model(&EventResponseModel{}).Where("id = ?", resp.ID).Update("project_id", projectID).Error; err != nil {
					return err
				}
			}
		}
		var oldAnswers []EventResponseAnswerModel
		if err := tx.Where("response_id = ?", resp.ID).Find(&oldAnswers).Error; err != nil {
			return err
		}
		for _, a := range oldAnswers {
			if a.OptionID != nil && *a.OptionID != "" {
				var opt EventOptionModel
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&opt, "id = ?", *a.OptionID).Error; err != nil {
					return err
				}
				if opt.CurrentCount > 0 {
					if err := tx.Model(&EventOptionModel{}).Where("id = ?", opt.ID).UpdateColumn("current_count", gorm.Expr("current_count - ?", 1)).Error; err != nil {
						return err
					}
				}
			}
		}
		if err := tx.Where("response_id = ?", resp.ID).Delete(&EventResponseAnswerModel{}).Error; err != nil {
			return err
		}
		for _, in := range answers {
			gid := in.GroupID
			row := EventResponseAnswerModel{
				ResponseID: resp.ID,
				GroupID:    &gid,
				OptionID:   in.OptionID,
				TextValue:  in.TextValue,
			}
			if in.OptionID != nil && *in.OptionID != "" {
				var opt EventOptionModel
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&opt, "id = ?", *in.OptionID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return eventdomain.ErrNotFound
					}
					return err
				}
				if opt.MaxCapacity != nil && opt.CurrentCount >= *opt.MaxCapacity {
					return eventdomain.ErrOptionCapacityReached
				}
				if err := tx.Model(&EventOptionModel{}).Where("id = ?", opt.ID).UpdateColumn("current_count", gorm.Expr("current_count + ?", 1)).Error; err != nil {
					return err
				}
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
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

func toDomainResp(m *EventResponseModel) *eventdomain.EventResponse {
	return &eventdomain.EventResponse{
		ID:              m.ID,
		EventInstanceID: m.EventInstanceID,
		UserID:          m.UserID,
		ProjectID:       m.ProjectID,
		CreatedAt:       m.CreatedAt,
	}
}

func toDomainAns(m *EventResponseAnswerModel) *eventdomain.EventResponseAnswer {
	return &eventdomain.EventResponseAnswer{
		ID:         m.ID,
		ResponseID: m.ResponseID,
		GroupID:    m.GroupID,
		OptionID:   m.OptionID,
		TextValue:  m.TextValue,
	}
}

func (r *RepositoryPG) LoadModuleResponseSummary(ctx context.Context, moduleID string) (*eventdomain.ModuleResponseSummary, error) {
	// Load module
	var modRow EventModuleModel
	if err := r.db.WithContext(ctx).First(&modRow, "id = ?", moduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, eventdomain.ErrNotFound
		}
		return nil, err
	}

	// Load groups for this module
	var groupRows []EventOptionGroupModel
	if err := r.db.WithContext(ctx).Where("module_id = ?", moduleID).Order("sort_order ASC").Find(&groupRows).Error; err != nil {
		return nil, err
	}

	if len(groupRows) == 0 {
		return &eventdomain.ModuleResponseSummary{Module: *toDomainMod(&modRow)}, nil
	}

	groupIDs := make([]string, len(groupRows))
	for i, g := range groupRows {
		groupIDs[i] = g.ID
	}

	// Load options for all groups
	var optRows []EventOptionModel
	if err := r.db.WithContext(ctx).Where("group_id IN ?", groupIDs).Order("sort_order ASC").Find(&optRows).Error; err != nil {
		return nil, err
	}
	optsByGroup := make(map[string][]EventOptionModel)
	for _, o := range optRows {
		optsByGroup[o.GroupID] = append(optsByGroup[o.GroupID], o)
	}

	// Load all answers for groups of this module, joined with user info
	type answerRow struct {
		ResponseID string         `gorm:"column:response_id"`
		GroupID    string         `gorm:"column:group_id"`
		OptionID   *string        `gorm:"column:option_id"`
		TextValue  *string        `gorm:"column:text_value"`
		UserID     string         `gorm:"column:user_id"`
		UserName   sql.NullString `gorm:"column:user_name"`
		UserEmail  sql.NullString `gorm:"column:user_email"`
		ProjectID  *string        `gorm:"column:project_id"`
	}
	var ansRows []answerRow
	err := r.db.WithContext(ctx).Raw(`
SELECT era.response_id, era.group_id, era.option_id, era.text_value,
       er.user_id, u.name AS user_name, u.email AS user_email, er.project_id
FROM event_response_answers era
JOIN event_responses er ON er.id = era.response_id
LEFT JOIN users u ON u.id = er.user_id
WHERE era.group_id IN ?
ORDER BY era.group_id, era.option_id, u.name`, groupIDs).Scan(&ansRows).Error
	if err != nil {
		return nil, err
	}

	// Index answers: group -> option -> []participants  (nil optionID = text/number)
	type optKey struct{ groupID, optionID string }
	optUsers := make(map[optKey][]eventdomain.SummaryParticipant)
	textAnswers := make(map[string][]eventdomain.TextAnswerSummary) // groupID -> []

	for _, a := range ansRows {
		p := eventdomain.SummaryParticipant{
			UserID:    a.UserID,
			ProjectID: a.ProjectID,
		}
		if a.UserName.Valid {
			p.UserName = a.UserName.String
		}
		if a.UserEmail.Valid {
			p.UserEmail = a.UserEmail.String
		}
		if a.OptionID != nil {
			k := optKey{a.GroupID, *a.OptionID}
			optUsers[k] = append(optUsers[k], p)
		} else {
			val := ""
			if a.TextValue != nil {
				val = *a.TextValue
			}
			textAnswers[a.GroupID] = append(textAnswers[a.GroupID], eventdomain.TextAnswerSummary{Value: val, User: p})
		}
	}

	// Build summary
	groups := make([]eventdomain.GroupSummary, 0, len(groupRows))
	for _, g := range groupRows {
		dg := toDomainGroup(&g)
		opts := optsByGroup[g.ID]
		optSummaries := make([]eventdomain.OptionSummary, 0, len(opts))
		for _, o := range opts {
			do := toDomainOpt(&o)
			k := optKey{g.ID, o.ID}
			optSummaries = append(optSummaries, eventdomain.OptionSummary{
				Option: *do,
				Users:  optUsers[k],
			})
		}
		groups = append(groups, eventdomain.GroupSummary{
			Group:       *dg,
			Options:     optSummaries,
			TextAnswers: textAnswers[g.ID],
		})
	}

	return &eventdomain.ModuleResponseSummary{
		Module: *toDomainMod(&modRow),
		Groups: groups,
	}, nil
}
