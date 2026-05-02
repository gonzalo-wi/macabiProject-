package mealdomain

import "time"

type Booking struct {
	ID              string
	UserID          string
	MealID          string
	GarnishOptionID *string
	Meal            *Meal
	GarnishOption   *GarnishOption
	CreatedAt       time.Time
}

func NewBooking(userID, mealID string, garnishOptionID *string) *Booking {
	return &Booking{
		UserID:          userID,
		MealID:          mealID,
		GarnishOptionID: garnishOptionID,
	}
}
