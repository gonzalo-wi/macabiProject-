package eventdto

import (
	"time"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type SubmitResponseRequest struct {
	ProjectID *string                   `json:"project_id"`
	Answers   []eventdomain.AnswerInput `json:"answers" binding:"required"`
}

type EventResponseJSON struct {
	ID              string    `json:"id"`
	EventInstanceID string    `json:"event_instance_id"`
	UserID          string    `json:"user_id"`
	ProjectID       *string   `json:"project_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type AnswerResponse struct {
	ID         string  `json:"id"`
	ResponseID string  `json:"response_id"`
	GroupID    *string `json:"group_id,omitempty"`
	OptionID   *string `json:"option_id,omitempty"`
	TextValue  *string `json:"text_value,omitempty"`
}

type MyResponseResult struct {
	Response *EventResponseJSON `json:"response"`
	Answers  []AnswerResponse   `json:"answers"`
}

type AdminEventResponseJSON struct {
	EventResponseJSON
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

type AdminResponseRow struct {
	Response AdminEventResponseJSON `json:"response"`
	Answers  []AnswerResponse       `json:"answers"`
}

type SummaryParticipantJSON struct {
	UserID    string  `json:"user_id"`
	UserName  string  `json:"user_name"`
	UserEmail string  `json:"user_email"`
	ProjectID *string `json:"project_id"`
}

type OptionSummaryJSON struct {
	ID           string                   `json:"id"`
	Label        string                   `json:"label"`
	MaxCapacity  *int                     `json:"max_capacity"`
	CurrentCount int                      `json:"current_count"`
	Count        int                      `json:"count"`
	Users        []SummaryParticipantJSON `json:"users"`
}

type TextAnswerSummaryJSON struct {
	Value string                 `json:"value"`
	User  SummaryParticipantJSON `json:"user"`
}

type GroupSummaryJSON struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Type        string                  `json:"type"`
	IsRequired  bool                    `json:"is_required"`
	Options     []OptionSummaryJSON     `json:"options"`
	TextAnswers []TextAnswerSummaryJSON `json:"text_answers"`
}

type ModuleResponseSummaryJSON struct {
	Module struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"module"`
	Groups []GroupSummaryJSON `json:"groups"`
}

func ToEventResponseJSON(r *eventdomain.EventResponse) EventResponseJSON {
	return EventResponseJSON{
		ID:              r.ID,
		EventInstanceID: r.EventInstanceID,
		UserID:          r.UserID,
		ProjectID:       r.ProjectID,
		CreatedAt:       r.CreatedAt,
	}
}

func ToAnswersResponse(aa []eventdomain.EventResponseAnswer) []AnswerResponse {
	out := make([]AnswerResponse, len(aa))
	for i, a := range aa {
		out[i] = AnswerResponse{
			ID:         a.ID,
			ResponseID: a.ResponseID,
			GroupID:    a.GroupID,
			OptionID:   a.OptionID,
			TextValue:  a.TextValue,
		}
	}
	return out
}

func ToMyResponseResult(r *eventdomain.EventResponse, aa []eventdomain.EventResponseAnswer) MyResponseResult {
	if r == nil {
		return MyResponseResult{Response: nil, Answers: []AnswerResponse{}}
	}
	resp := ToEventResponseJSON(r)
	return MyResponseResult{Response: &resp, Answers: ToAnswersResponse(aa)}
}

func ToAdminResponseRows(list []eventdomain.EventResponseWithParticipant) []AdminResponseRow {
	rows := make([]AdminResponseRow, len(list))
	for i, row := range list {
		rows[i] = AdminResponseRow{
			Response: AdminEventResponseJSON{
				EventResponseJSON: ToEventResponseJSON(&row.Response),
				UserName:          row.UserName,
				UserEmail:         row.UserEmail,
			},
			Answers: ToAnswersResponse(row.Answers),
		}
	}
	return rows
}

func ToAdminResponseListResponse(result pagination.Result[eventdomain.EventResponseWithParticipant]) pagination.Result[AdminResponseRow] {
	return pagination.Result[AdminResponseRow]{
		Data:       ToAdminResponseRows(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

func ToModuleResponseSummaryJSON(s *eventdomain.ModuleResponseSummary) ModuleResponseSummaryJSON {
	groups := make([]GroupSummaryJSON, len(s.Groups))
	for i, g := range s.Groups {
		opts := make([]OptionSummaryJSON, len(g.Options))
		for j, o := range g.Options {
			users := make([]SummaryParticipantJSON, len(o.Users))
			for k, u := range o.Users {
				users[k] = SummaryParticipantJSON{UserID: u.UserID, UserName: u.UserName, UserEmail: u.UserEmail, ProjectID: u.ProjectID}
			}
			opts[j] = OptionSummaryJSON{
				ID:           o.Option.ID,
				Label:        o.Option.Label,
				MaxCapacity:  o.Option.MaxCapacity,
				CurrentCount: o.Option.CurrentCount,
				Count:        len(users),
				Users:        users,
			}
		}
		texts := make([]TextAnswerSummaryJSON, len(g.TextAnswers))
		for j, t := range g.TextAnswers {
			texts[j] = TextAnswerSummaryJSON{
				Value: t.Value,
				User:  SummaryParticipantJSON{UserID: t.User.UserID, UserName: t.User.UserName, UserEmail: t.User.UserEmail, ProjectID: t.User.ProjectID},
			}
		}
		groups[i] = GroupSummaryJSON{
			ID:          g.Group.ID,
			Name:        g.Group.Name,
			Type:        string(g.Group.Type),
			IsRequired:  g.Group.IsRequired,
			Options:     opts,
			TextAnswers: texts,
		}
	}
	var out ModuleResponseSummaryJSON
	out.Module.ID = s.Module.ID
	out.Module.Title = s.Module.Title
	out.Module.Type = string(s.Module.Type)
	out.Groups = groups
	return out
}
