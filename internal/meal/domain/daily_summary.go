package mealdomain

import "time"

type PersonSummary struct {
	Name string
}

type MealDailySummary struct {
	MealID   string
	Title    string
	Quantity int
	Persons  []PersonSummary
}

type DailySummary struct {
	Date          time.Time
	TotalMenus    int
	MealSummaries []MealDailySummary
}
