package eventdomain

import "time"

type EventType string

const (
	EventTypeActivity EventType = "activity"
	EventTypeCustom   EventType = "custom"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusOpen      EventStatus = "open"
	EventStatusClosed    EventStatus = "closed"
	EventStatusCancelled EventStatus = "cancelled"
)

type ModuleType string

const (
	ModuleAttendance ModuleType = "attendance"
	ModuleMeal       ModuleType = "meal"
	ModuleTransport  ModuleType = "transport"
	ModuleMaterials  ModuleType = "materials"
	ModuleCustom     ModuleType = "custom"
)

type OptionGroupType string

const (
	OptionGroupSingleChoice   OptionGroupType = "single_choice"
	OptionGroupMultipleChoice OptionGroupType = "multiple_choice"
	OptionGroupText           OptionGroupType = "text"
	OptionGroupNumber         OptionGroupType = "number"
)

type EventInstance struct {
	ID                 string
	Title              string
	Type               EventType
	StartsAt           time.Time
	ResponseDeadlineAt *time.Time
	Status             EventStatus
	CreatedAt          time.Time
}

type EventModule struct {
	ID              string
	EventInstanceID string
	Title           string
	Type            ModuleType
	SortOrder       int
	IsRequired      bool
	CreatedAt       time.Time
}

type EventOptionGroup struct {
	ID         string
	ModuleID   string
	Name       string
	Type       OptionGroupType
	IsRequired bool
	SortOrder  int
}

type EventOption struct {
	ID           string
	GroupID      string
	Label        string
	MaxCapacity  *int
	CurrentCount int
	SortOrder    int
}

type EventResponse struct {
	ID              string
	EventInstanceID string
	UserID          string
	ProjectID       *string
	CreatedAt       time.Time
}

type EventResponseAnswer struct {
	ID         string
	ResponseID string
	GroupID    *string
	OptionID   *string
	TextValue  *string
}

type EventResponseWithParticipant struct {
	Response  EventResponse
	UserName  string
	UserEmail string
	Answers   []EventResponseAnswer
}

type EventDetail struct {
	Instance   EventInstance
	ProjectIDs []string
	Modules    []ModuleDetail
}

type ModuleDetail struct {
	Module       EventModule
	ProjectIDs   []string
	OptionGroups []GroupDetail
}

type GroupDetail struct {
	Group   EventOptionGroup
	Options []EventOption
}

type AnswerInput struct {
	GroupID   string  `json:"group_id"`
	OptionID  *string `json:"option_id"`
	TextValue *string `json:"text_value"`
}

type SummaryParticipant struct {
	UserID    string
	UserName  string
	UserEmail string
	ProjectID *string
}

type OptionSummary struct {
	Option EventOption
	Users  []SummaryParticipant
}

type TextAnswerSummary struct {
	Value string
	User  SummaryParticipant
}

type GroupSummary struct {
	Group       EventOptionGroup
	Options     []OptionSummary
	TextAnswers []TextAnswerSummary
}

type ModuleResponseSummary struct {
	Module EventModule
	Groups []GroupSummary
}

type EventNotification struct {
	ID              string
	UserID          string
	EventInstanceID string
	Message         string
	ReadAt          *time.Time
	CreatedAt       time.Time
}
