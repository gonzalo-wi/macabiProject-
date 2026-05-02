package mealhttp

import (
	"time"

	mealdomain "macabi-back/internal/meal/domain"
)

// ── MealTemplate ─────────────────────────────────────────────────────────────

type CreateMealTemplateRequest struct {
	Title       string `json:"title" binding:"required"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required"`
	Type        string `json:"type" binding:"required"`
}

type UpdateMealTemplateRequest struct {
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Type        string `json:"type"`
}

type GarnishOptionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MealTemplateResponse struct {
	ID             string                  `json:"id"`
	Title          string                  `json:"title"`
	ImageURL       string                  `json:"image_url"`
	Description    string                  `json:"description"`
	Category       string                  `json:"category"`
	Type           string                  `json:"type"`
	GarnishOptions []GarnishOptionResponse `json:"garnish_options"`
	CreatedAt      time.Time               `json:"created_at"`
}

type PaginatedMealTemplateResponse struct {
	Data       []MealTemplateResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

func ToMealTemplateResponse(t *mealdomain.MealTemplate) MealTemplateResponse {
	garnishOptions := make([]GarnishOptionResponse, len(t.GarnishOptions))
	for i, g := range t.GarnishOptions {
		garnishOptions[i] = GarnishOptionResponse{ID: g.ID, Name: g.Name}
	}
	return MealTemplateResponse{
		ID:             t.ID,
		Title:          t.Title,
		ImageURL:       t.ImageURL,
		Description:    t.Description,
		Category:       t.Category.String(),
		Type:           t.Type.String(),
		GarnishOptions: garnishOptions,
		CreatedAt:      t.CreatedAt,
	}
}

type CreateGarnishOptionRequest struct {
	Name string `json:"name" binding:"required"`
}

// ── Meal ─────────────────────────────────────────────────────────────────────

type CreateMealRequest struct {
	ProjectID      string    `json:"project_id" binding:"required"`
	TemplateID     string    `json:"template_id" binding:"required"`
	AvailableCount int       `json:"available_count" binding:"min=0"`
	Date           time.Time `json:"date" binding:"required"`
}

type MealResponse struct {
	ID             string                  `json:"id"`
	ProjectID      string                  `json:"project_id"`
	TemplateID     string                  `json:"template_id"`
	Title          string                  `json:"title"`
	ImageURL       string                  `json:"image_url"`
	Description    string                  `json:"description"`
	Category       string                  `json:"category"`
	Type           string                  `json:"type"`
	GarnishOptions []GarnishOptionResponse `json:"garnish_options"`
	SoldOut        bool                    `json:"sold_out"`
	AvailableCount int                     `json:"available_count"`
	Date           time.Time               `json:"date"`
	CreatedAt      time.Time               `json:"created_at"`
}

func ToMealResponse(m *mealdomain.Meal) MealResponse {
	resp := MealResponse{
		ID:             m.ID,
		ProjectID:      m.ProjectID,
		TemplateID:     m.TemplateID,
		SoldOut:        m.SoldOut,
		AvailableCount: m.AvailableCount,
		Date:           m.Date,
		CreatedAt:      m.CreatedAt,
	}
	if m.Template != nil {
		resp.Title = m.Template.Title
		resp.ImageURL = m.Template.ImageURL
		resp.Description = m.Template.Description
		resp.Category = m.Template.Category.String()
		resp.Type = m.Template.Type.String()
		garnishOptions := make([]GarnishOptionResponse, len(m.Template.GarnishOptions))
		for i, g := range m.Template.GarnishOptions {
			garnishOptions[i] = GarnishOptionResponse{ID: g.ID, Name: g.Name}
		}
		resp.GarnishOptions = garnishOptions
	}
	return resp
}

type PaginatedMealResponse struct {
	Data       []MealResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// ── Booking ───────────────────────────────────────────────────────────────────

type BookMealRequest struct {
	MealID          string  `json:"meal_id" binding:"required"`
	GarnishOptionID *string `json:"garnish_option_id"`
}

type BookingResponse struct {
	ID            string                 `json:"id"`
	MealID        string                 `json:"meal_id"`
	GarnishOption *GarnishOptionResponse `json:"garnish_option,omitempty"`
	Meal          *MealResponse          `json:"meal,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
}

func ToBookingResponse(b *mealdomain.Booking) BookingResponse {
	resp := BookingResponse{
		ID:        b.ID,
		MealID:    b.MealID,
		CreatedAt: b.CreatedAt,
	}
	if b.Meal != nil {
		meal := ToMealResponse(b.Meal)
		resp.Meal = &meal
	}
	if b.GarnishOption != nil {
		resp.GarnishOption = &GarnishOptionResponse{ID: b.GarnishOption.ID, Name: b.GarnishOption.Name}
	}
	return resp
}

type PaginatedBookingResponse struct {
	Data       []BookingResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// ── Daily Summary ─────────────────────────────────────────────────────────────

type MealDailySummaryResponse struct {
	MealID   string   `json:"meal_id"`
	Title    string   `json:"title"`
	Quantity int      `json:"quantity"`
	Persons  []string `json:"persons"`
}

type ProjectDailySummaryResponse struct {
	ProjectID     string                     `json:"project_id"`
	ProjectName   string                     `json:"project_name"`
	TotalMenus    int                        `json:"total_menus"`
	MealSummaries []MealDailySummaryResponse `json:"meal_summaries"`
}

type DailySummaryResponse struct {
	Date       time.Time                     `json:"date"`
	TotalMenus int                           `json:"total_menus"`
	Projects   []ProjectDailySummaryResponse `json:"projects"`
}

func ToDailySummaryResponse(s *mealdomain.DailySummary) DailySummaryResponse {
	projects := make([]ProjectDailySummaryResponse, len(s.Projects))
	for i, p := range s.Projects {
		summaries := make([]MealDailySummaryResponse, len(p.MealSummaries))
		for j, ms := range p.MealSummaries {
			persons := make([]string, len(ms.Persons))
			for k, person := range ms.Persons {
				persons[k] = person.Name
			}
			summaries[j] = MealDailySummaryResponse{
				MealID:   ms.MealID,
				Title:    ms.Title,
				Quantity: ms.Quantity,
				Persons:  persons,
			}
		}
		projects[i] = ProjectDailySummaryResponse{
			ProjectID:     p.ProjectID,
			ProjectName:   p.ProjectName,
			TotalMenus:    p.TotalMenus,
			MealSummaries: summaries,
		}
	}
	return DailySummaryResponse{
		Date:       s.Date,
		TotalMenus: s.TotalMenus,
		Projects:   projects,
	}
}
