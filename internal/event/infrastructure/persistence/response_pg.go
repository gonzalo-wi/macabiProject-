package eventpersistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

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

func (r *RepositoryPG) LoadModuleResponseSummary(ctx context.Context, moduleID string) (*eventdomain.ModuleResponseSummary, error) {
	var modRow EventModuleModel
	if err := r.db.WithContext(ctx).First(&modRow, "id = ?", moduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, eventdomain.ErrNotFound
		}
		return nil, err
	}

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

	var optRows []EventOptionModel
	if err := r.db.WithContext(ctx).Where("group_id IN ?", groupIDs).Order("sort_order ASC").Find(&optRows).Error; err != nil {
		return nil, err
	}
	optsByGroup := make(map[string][]EventOptionModel)
	for _, o := range optRows {
		optsByGroup[o.GroupID] = append(optsByGroup[o.GroupID], o)
	}

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

	type optKey struct{ groupID, optionID string }
	optUsers := make(map[optKey][]eventdomain.SummaryParticipant)
	textAnswers := make(map[string][]eventdomain.TextAnswerSummary)

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
