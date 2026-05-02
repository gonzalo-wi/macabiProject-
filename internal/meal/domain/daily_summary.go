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

type ProjectDailySummary struct {
	ProjectID     string
	ProjectName   string
	TotalMenus    int
	MealSummaries []MealDailySummary
}

type DailySummary struct {
	Date       time.Time
	TotalMenus int
	Projects   []ProjectDailySummary
}
