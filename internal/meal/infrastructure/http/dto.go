package mealhttp

import (
	"time"

	mealdomain "macabi-back/internal/meal/domain"
)

type CreateMealRequest struct {
	Title          string    `json:"title" binding:"required"`
	ImageURL       string    `json:"image_url"`
	Description    string    `json:"description"`
	Category       string    `json:"category" binding:"required"`
	Type           string    `json:"type" binding:"required"`
	AvailableCount int       `json:"available_count" binding:"min=0"`
	Date           time.Time `json:"date" binding:"required"`
}

type MealResponse struct {
	ID             string    `json:"id"`
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
	return MealResponse{
		ID:             m.ID,
		Title:          m.Title,
		ImageURL:       m.ImageURL,
		Description:    m.Description,
		Category:       m.Category.String(),
		Type:           m.Type.String(),
		SoldOut:        m.SoldOut,
		AvailableCount: m.AvailableCount,
		Date:           m.Date,
		CreatedAt:      m.CreatedAt,
	}
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
