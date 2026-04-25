package mealpersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	mealdomain "macabi-back/internal/meal/domain"
	"macabi-back/internal/shared/database"
	"macabi-back/internal/shared/pagination"

	"gorm.io/gorm"
)

type MealModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title          string `gorm:"not null"`
	ImageURL       string
	Description    string
	Category       string    `gorm:"not null"`
	Type           string    `gorm:"not null"`
	SoldOut        bool      `gorm:"not null;default:false"`
	AvailableCount int       `gorm:"not null;default:0"`
	Date           time.Time `gorm:"not null;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (MealModel) TableName() string {
	return "meals"
}

type MealRepositoryPG struct {
	db *gorm.DB
}

func NewMealRepositoryPG(db *gorm.DB) *MealRepositoryPG {
	return &MealRepositoryPG{db: db}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&MealModel{}, &BookingModel{})
}

func (r *MealRepositoryPG) tx(ctx context.Context) *gorm.DB {
	return database.TxFromCtx(ctx, r.db).WithContext(ctx)
}

func (r *MealRepositoryPG) Save(ctx context.Context, meal *mealdomain.Meal) error {
	model := toMealModel(meal)
	if err := r.tx(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("save meal: %w", err)
	}
	meal.ID = model.ID
	meal.CreatedAt = model.CreatedAt
	meal.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *MealRepositoryPG) FindByID(ctx context.Context, id string) (*mealdomain.Meal, error) {
	var model MealModel
	err := r.tx(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, mealdomain.ErrMealNotFound
		}
		return nil, fmt.Errorf("find meal by id: %w", err)
	}
	return toDomainMeal(&model), nil
}

func (r *MealRepositoryPG) FindByDate(ctx context.Context, date time.Time, params pagination.Params) ([]mealdomain.Meal, int64, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var total int64
	if err := r.tx(ctx).Model(&MealModel{}).
		Where("date >= ? AND date < ?", start, end).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count meals: %w", err)
	}

	var models []MealModel
	err := r.tx(ctx).
		Where("date >= ? AND date < ?", start, end).
		Order("type ASC, category ASC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find meals by date: %w", err)
	}

	meals := make([]mealdomain.Meal, len(models))
	for i := range models {
		meals[i] = *toDomainMeal(&models[i])
	}
	return meals, total, nil
}

func (r *MealRepositoryPG) Update(ctx context.Context, meal *mealdomain.Meal) error {
	err := r.tx(ctx).Model(&MealModel{}).Where("id = ?", meal.ID).Updates(map[string]interface{}{
		"title":           meal.Title,
		"image_url":       meal.ImageURL,
		"description":     meal.Description,
		"category":        string(meal.Category),
		"type":            string(meal.Type),
		"sold_out":        meal.SoldOut,
		"available_count": meal.AvailableCount,
		"date":            meal.Date,
	}).Error
	if err != nil {
		return fmt.Errorf("update meal: %w", err)
	}
	return nil
}

func (r *MealRepositoryPG) Delete(ctx context.Context, id string) error {
	err := r.tx(ctx).Where("id = ?", id).Delete(&MealModel{}).Error
	if err != nil {
		return fmt.Errorf("delete meal: %w", err)
	}
	return nil
}

func toMealModel(meal *mealdomain.Meal) MealModel {
	return MealModel{
		ID:             meal.ID,
		Title:          meal.Title,
		ImageURL:       meal.ImageURL,
		Description:    meal.Description,
		Category:       string(meal.Category),
		Type:           string(meal.Type),
		SoldOut:        meal.SoldOut,
		AvailableCount: meal.AvailableCount,
		Date:           meal.Date,
	}
}

func toDomainMeal(model *MealModel) *mealdomain.Meal {
	return &mealdomain.Meal{
		ID:             model.ID,
		Title:          model.Title,
		ImageURL:       model.ImageURL,
		Description:    model.Description,
		Category:       mealdomain.Category(model.Category),
		Type:           mealdomain.MealType(model.Type),
		SoldOut:        model.SoldOut,
		AvailableCount: model.AvailableCount,
		Date:           model.Date,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}
