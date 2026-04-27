package mealhttp

import (
	"time"

	mealdomain "macabi-back/internal/meal/domain"
)

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

type MealTemplateResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginatedMealTemplateResponse struct {
	Data       []MealTemplateResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

func ToMealTemplateResponse(t *mealdomain.MealTemplate) MealTemplateResponse {
	return MealTemplateResponse{
		ID:          t.ID,
		Title:       t.Title,
		ImageURL:    t.ImageURL,
		Description: t.Description,
		Category:    t.Category.String(),
		Type:        t.Type.String(),
		CreatedAt:   t.CreatedAt,
	}
}

type CreateMealRequest struct {
	TemplateID     string    `json:"template_id" binding:"required"`
	AvailableCount int       `json:"available_count" binding:"min=0"`
	Date           time.Time `json:"date" binding:"required"`
}

type MealResponse struct {
	ID             string    `json:"id"`
	TemplateID     string    `json:"template_id"`
	Title          string    `json:"title"`
	ImageURL       string    `json:"image_url"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	Type           string    `json:"type"`
	SoldOut        bool      `json:"sold_out"`
	AvailableCount int       `json:"available_count"`
	Date           time.Time `json:"date"`
	CreatedAt      time.Time `json:"created_at"`
}

type BookMealRequest struct {
	MealID string `json:"meal_id" binding:"required"`
}

type BookingResponse struct {
	ID        string        `json:"id"`
	MealID    string        `json:"meal_id"`
	Meal      *MealResponse `json:"meal,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

type PaginatedMealResponse struct {
	Data       []MealResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type PaginatedBookingResponse struct {
	Data       []BookingResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

func ToMealResponse(m *mealdomain.Meal) MealResponse {
	r := MealResponse{
		ID:             m.ID,
		TemplateID:     m.TemplateID,
		SoldOut:        m.SoldOut,
		AvailableCount: m.AvailableCount,
		Date:           m.Date,
		CreatedAt:      m.CreatedAt,
	}
	if m.Template != nil {
		r.Title = m.Template.Title
		r.ImageURL = m.Template.ImageURL
		r.Description = m.Template.Description
		r.Category = m.Template.Category.String()
		r.Type = m.Template.Type.String()
	}
	return r
}

func ToBookingResponse(b *mealdomain.Booking) BookingResponse {
	resp := BookingResponse{
		ID:        b.ID,
		MealID:    b.MealID,
		CreatedAt: b.CreatedAt,
	}
	if b.Meal != nil {
		mr := ToMealResponse(b.Meal)
		resp.Meal = &mr
	}
	return resp
}

// --- Daily summary (admin) ---

type PersonSummaryResponse struct {
	Nombre string `json:"nombre"`
}

type MealDailySummaryResponse struct {
	MenuID   string                  `json:"menuId"`
	Nombre   string                  `json:"nombre"`
	Cantidad int                     `json:"cantidad"`
	Personas []PersonSummaryResponse `json:"personas"`
}

type DailySummaryResponse struct {
	Fecha      string                     `json:"fecha"`
	TotalMenus int                        `json:"totalMenus"`
	PorMenu    []MealDailySummaryResponse `json:"porMenu"`
}

func ToDailySummaryResponse(s *mealdomain.DailySummary) DailySummaryResponse {
	porMenu := make([]MealDailySummaryResponse, len(s.MealSummaries))
	for i, ms := range s.MealSummaries {
		personas := make([]PersonSummaryResponse, len(ms.Persons))
		for j, p := range ms.Persons {
			personas[j] = PersonSummaryResponse{Nombre: p.Name}
		}
		porMenu[i] = MealDailySummaryResponse{
			MenuID:   ms.MealID,
			Nombre:   ms.Title,
			Cantidad: ms.Quantity,
			Personas: personas,
		}
	}
	return DailySummaryResponse{
		Fecha:      s.Date.Format("2006-01-02"),
		TotalMenus: s.TotalMenus,
		PorMenu:    porMenu,
	}
}
