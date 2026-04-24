package mealdomain

import "time"

type Booking struct {
	ID        string
	UserID    string
	MealID    string
	Meal      *Meal
	CreatedAt time.Time
}

func NewBooking(userID, mealID string) *Booking {
	return &Booking{
		UserID: userID,
		MealID: mealID,
	}
}
