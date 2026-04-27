package mealdomain

import (
	"strings"
	"time"
)

type MealTemplate struct {
	ID          string
	Title       string
	ImageURL    string
	Description string
	Category    Category
	Type        MealType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewMealTemplate(title, imageURL, description string, category Category, mealType MealType) (*MealTemplate, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrEmptyTitle
	}
	return &MealTemplate{
		Title:       title,
		ImageURL:    imageURL,
		Description: description,
		Category:    category,
		Type:        mealType,
	}, nil
}
