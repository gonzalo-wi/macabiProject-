package mealdomain

import (
	"time"
)

type Meal struct {
	ID             string
	ProjectID      string
	TemplateID     string
	Template       *MealTemplate
	SoldOut        bool
	AvailableCount int
	Date           time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewMeal(projectID, templateID string, availableCount int, date time.Time) (*Meal, error) {
	if availableCount < 0 {
		return nil, ErrInvalidAvailableCount
	}
	if !isValidMealDay(date) {
		return nil, ErrInvalidDate
	}
	return &Meal{
		ProjectID:      projectID,
		TemplateID:     templateID,
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

const (
	BookingDeadlineWeekday = time.Friday
	BookingDeadlineHour    = 11
	BookingDeadlineMinute  = 59
)

func isValidMealDay(date time.Time) bool {
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}

func IsBookingOpen(mealDate time.Time, now time.Time) bool {
	frida := mealDate.AddDate(0, 0, -1)
	deadline := time.Date(frida.Year(), frida.Month(), frida.Day(),
		BookingDeadlineHour, BookingDeadlineMinute, 59, 0, mealDate.Location())
	return now.Before(deadline) || now.Equal(deadline)
}
