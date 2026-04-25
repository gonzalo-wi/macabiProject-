package mealpersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/pagination"

	"gorm.io/gorm"
)

type BookingModel struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	MealID    string    `gorm:"type:uuid;not null"`
	Meal      MealModel `gorm:"foreignKey:MealID"`
	CreatedAt time.Time
}

func (BookingModel) TableName() string {
	return "meal_bookings"
}

type BookingRepositoryPG struct {
	db *gorm.DB
}

func NewBookingRepositoryPG(db *gorm.DB) *BookingRepositoryPG {
	return &BookingRepositoryPG{db: db}
}

func (r *BookingRepositoryPG) Save(ctx context.Context, booking *mealdomain.Booking) error {
	model := BookingModel{
		UserID: booking.UserID,
		MealID: booking.MealID,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("save booking: %w", err)
	}
	booking.ID = model.ID
	booking.CreatedAt = model.CreatedAt
	return nil
}

func (r *BookingRepositoryPG) FindByID(ctx context.Context, id string) (*mealdomain.Booking, error) {
	var model BookingModel
	err := r.db.WithContext(ctx).Preload("Meal").Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, mealdomain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("find booking by id: %w", err)
	}
	return toDomainBooking(&model), nil
}

func (r *BookingRepositoryPG) FindByUserID(ctx context.Context, userID string, params pagination.Params) ([]mealdomain.Booking, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&BookingModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	var models []BookingModel
	err := r.db.WithContext(ctx).
		Preload("Meal").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find bookings by user: %w", err)
	}

	bookings := make([]mealdomain.Booking, len(models))
	for i := range models {
		bookings[i] = *toDomainBooking(&models[i])
	}
	return bookings, total, nil
}

func (r *BookingRepositoryPG) FindByUserAndMealTypeAndDate(ctx context.Context, userID string, mealType mealdomain.MealType, date time.Time) (*mealdomain.Booking, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var model BookingModel
	err := r.db.WithContext(ctx).
		Select("meal_bookings.*").
		Joins("JOIN meals ON meals.id = meal_bookings.meal_id").
		Where("meal_bookings.user_id = ?", userID).
		Where("meals.type = ?", string(mealType)).
		Where("meals.date >= ? AND meals.date < ?", start, end).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, mealdomain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("find booking by user/type/date: %w", err)
	}
	return toDomainBooking(&model), nil
}

func (r *BookingRepositoryPG) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&BookingModel{}).Error
	if err != nil {
		return fmt.Errorf("delete booking: %w", err)
	}
	return nil
}

func toDomainBooking(model *BookingModel) *mealdomain.Booking {
	b := &mealdomain.Booking{
		ID:        model.ID,
		UserID:    model.UserID,
		MealID:    model.MealID,
		CreatedAt: model.CreatedAt,
	}
	if model.Meal.ID != "" {
		b.Meal = toDomainMeal(&model.Meal)
	}
	return b
}

type dailySummaryRow struct {
	MealID    string
	MealTitle string
	UserName  string
}

func (r *BookingRepositoryPG) GetDailySummary(ctx context.Context, date time.Time) (*mealdomain.DailySummary, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	var rows []dailySummaryRow
	err := r.db.WithContext(ctx).
		Table("meal_bookings b").
		Select("m.id AS meal_id, m.title AS meal_title, u.name AS user_name").
		Joins("JOIN meals m ON m.id = b.meal_id").
		Joins("JOIN users u ON u.id = b.user_id").
		Where("m.date >= ? AND m.date < ?", start, end).
		Order("m.id, u.name").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("get daily summary: %w", err)
	}

	indexByMealID := make(map[string]int)
	summaries := make([]mealdomain.MealDailySummary, 0)
	for _, row := range rows {
		idx, exists := indexByMealID[row.MealID]
		if !exists {
			summaries = append(summaries, mealdomain.MealDailySummary{
				MealID:  row.MealID,
				Title:   row.MealTitle,
				Persons: []mealdomain.PersonSummary{},
			})
			idx = len(summaries) - 1
			indexByMealID[row.MealID] = idx
		}
		summaries[idx].Quantity++
		summaries[idx].Persons = append(summaries[idx].Persons, mealdomain.PersonSummary{
			Name: row.UserName,
		})
	}

	total := 0
	for _, s := range summaries {
		total += s.Quantity
	}

	return &mealdomain.DailySummary{
		Date:          start,
		TotalMenus:    total,
		MealSummaries: summaries,
	}, nil
}
