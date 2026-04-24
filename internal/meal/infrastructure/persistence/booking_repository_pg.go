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

func (r *BookingRepositoryPG) FindByUserAndMealTypeAndDate(ctx context.Context, userID string, mealType mealdomain.MealType, date time.Time, isPostre bool) (*mealdomain.Booking, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var model BookingModel
	query := r.db.WithContext(ctx).
		Select("meal_bookings.*").
		Joins("JOIN meals ON meals.id = meal_bookings.meal_id").
		Where("meal_bookings.user_id = ?", userID).
		Where("meals.type = ?", string(mealType)).
		Where("meals.date >= ? AND meals.date < ?", start, end)

	if isPostre {
		query = query.Where("meals.category = ?", string(mealdomain.CategoryPostre))
	} else {
		query = query.Where("meals.category != ?", string(mealdomain.CategoryPostre))
	}

	err := query.First(&model).Error
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
