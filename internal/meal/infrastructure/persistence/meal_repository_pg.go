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

// ── MealTemplate ────────────────────────────────────────────────────────────

type MealTemplateModel struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string `gorm:"not null"`
	ImageURL    string
	Description string
	Category    string `gorm:"not null"`
	Type        string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (MealTemplateModel) TableName() string {
	return "meal_templates"
}

type MealTemplateRepositoryPG struct {
	db *gorm.DB
}

func NewMealTemplateRepositoryPG(db *gorm.DB) *MealTemplateRepositoryPG {
	return &MealTemplateRepositoryPG{db: db}
}

func (r *MealTemplateRepositoryPG) Save(ctx context.Context, tmpl *mealdomain.MealTemplate) error {
	model := toMealTemplateModel(tmpl)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("save meal template: %w", err)
	}
	tmpl.ID = model.ID
	tmpl.CreatedAt = model.CreatedAt
	tmpl.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *MealTemplateRepositoryPG) FindByID(ctx context.Context, id string) (*mealdomain.MealTemplate, error) {
	var model MealTemplateModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, mealdomain.ErrTemplateNotFound
		}
		return nil, fmt.Errorf("find meal template by id: %w", err)
	}
	return toDomainMealTemplate(&model), nil
}

func (r *MealTemplateRepositoryPG) FindAll(ctx context.Context, params pagination.Params) ([]mealdomain.MealTemplate, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&MealTemplateModel{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count meal templates: %w", err)
	}

	var models []MealTemplateModel
	err := r.db.WithContext(ctx).
		Order("title ASC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find all meal templates: %w", err)
	}

	templates := make([]mealdomain.MealTemplate, len(models))
	for i := range models {
		templates[i] = *toDomainMealTemplate(&models[i])
	}
	return templates, total, nil
}

func (r *MealTemplateRepositoryPG) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&MealTemplateModel{}).Error; err != nil {
		return fmt.Errorf("delete meal template: %w", err)
	}
	return nil
}

func (r *MealTemplateRepositoryPG) Update(ctx context.Context, tmpl *mealdomain.MealTemplate) error {
	err := r.db.WithContext(ctx).Model(&MealTemplateModel{}).Where("id = ?", tmpl.ID).Updates(map[string]interface{}{
		"title":       tmpl.Title,
		"image_url":   tmpl.ImageURL,
		"description": tmpl.Description,
		"category":    string(tmpl.Category),
		"type":        string(tmpl.Type),
	}).Error
	if err != nil {
		return fmt.Errorf("update meal template: %w", err)
	}
	return nil
}

func toMealTemplateModel(tmpl *mealdomain.MealTemplate) MealTemplateModel {
	return MealTemplateModel{
		ID:          tmpl.ID,
		Title:       tmpl.Title,
		ImageURL:    tmpl.ImageURL,
		Description: tmpl.Description,
		Category:    string(tmpl.Category),
		Type:        string(tmpl.Type),
	}
}

func toDomainMealTemplate(model *MealTemplateModel) *mealdomain.MealTemplate {
	return &mealdomain.MealTemplate{
		ID:          model.ID,
		Title:       model.Title,
		ImageURL:    model.ImageURL,
		Description: model.Description,
		Category:    mealdomain.Category(model.Category),
		Type:        mealdomain.MealType(model.Type),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

// ── Meal ─────────────────────────────────────────────────────────────────────

type MealModel struct {
	ID             string            `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TemplateID     string            `gorm:"type:uuid;not null"`
	Template       MealTemplateModel `gorm:"foreignKey:TemplateID"`
	SoldOut        bool              `gorm:"not null;default:false"`
	AvailableCount int               `gorm:"not null;default:0"`
	Date           time.Time         `gorm:"not null;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (MealModel) TableName() string {
	return "meals"
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&MealTemplateModel{}, &MealModel{}, &BookingModel{})
}

type MealRepositoryPG struct {
	db *gorm.DB
}

func NewMealRepositoryPG(db *gorm.DB) *MealRepositoryPG {
	return &MealRepositoryPG{db: db}
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
	err := r.tx(ctx).Preload("Template").Where("id = ?", id).First(&model).Error
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
		Preload("Template").
		Where("date >= ? AND date < ?", start, end).
		Order("date ASC").
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
		"sold_out":        meal.SoldOut,
		"available_count": meal.AvailableCount,
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
		TemplateID:     meal.TemplateID,
		SoldOut:        meal.SoldOut,
		AvailableCount: meal.AvailableCount,
		Date:           meal.Date,
	}
}

func toDomainMeal(model *MealModel) *mealdomain.Meal {
	m := &mealdomain.Meal{
		ID:             model.ID,
		TemplateID:     model.TemplateID,
		SoldOut:        model.SoldOut,
		AvailableCount: model.AvailableCount,
		Date:           model.Date,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
	if model.Template.ID != "" {
		m.Template = toDomainMealTemplate(&model.Template)
	}
	return m
}
