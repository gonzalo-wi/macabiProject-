package mealdomain

import (
	"strings"
	"time"
)

type Meal struct {
	ID             string
	Title          string
	ImageURL       string
	Description    string
	Category       Category
	Type           MealType
	SoldOut        bool
	AvailableCount int
	Date           time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewMeal(title, imageURL, description string, category Category, mealType MealType, availableCount int, date time.Time) (*Meal, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if availableCount < 0 {
		return nil, ErrInvalidAvailableCount
	}
	if !isSaturday(date) {
		return nil, ErrInvalidDate
	}
	return &Meal{
		Title:          title,
		ImageURL:       imageURL,
		Description:    description,
		Category:       category,
		Type:           mealType,
		SoldOut:        availableCount == 0,
		AvailableCount: availableCount,
		Date:           date,
	}, nil
}

func (m *Meal) DecrementAvailable() error {
	if m.SoldOut || m.AvailableCount == 0 {
		return ErrMealSoldOut
	}
	m.AvailableCount--
	if m.AvailableCount == 0 {
		m.SoldOut = true
	}
	return nil
}

func (m *Meal) IncrementAvailable() {
	m.AvailableCount++
	m.SoldOut = false
}

// BookingDeadlineWeekday and BookingDeadlineHour define the cutoff for reservations.
// Change these values to adjust the deadline (e.g. Friday 23:59).
const (
	BookingDeadlineWeekday = time.Friday
	BookingDeadlineHour    = 23
	BookingDeadlineMinute  = 59
)

func isSaturday(date time.Time) bool {
	return date.Weekday() == time.Saturday
}

// IsBookingOpen reports whether reservations are still accepted for the given meal date.
func IsBookingOpen(mealDate time.Time, now time.Time) bool {
	// Find the Friday before the Saturday meal date
	frida := mealDate.AddDate(0, 0, -1)
	deadline := time.Date(frida.Year(), frida.Month(), frida.Day(),
		BookingDeadlineHour, BookingDeadlineMinute, 59, 0, mealDate.Location())
	return now.Before(deadline) || now.Equal(deadline)
}
