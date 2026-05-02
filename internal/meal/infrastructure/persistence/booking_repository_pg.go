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

type BookingModel struct {
	ID              string              `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          string              `gorm:"type:uuid;not null;index"`
	MealID          string              `gorm:"type:uuid;not null"`
	GarnishOptionID *string             `gorm:"type:uuid"`
	Meal            MealModel           `gorm:"foreignKey:MealID"`
	GarnishOption   *GarnishOptionModel `gorm:"foreignKey:GarnishOptionID"`
	CreatedAt       time.Time
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

func (r *BookingRepositoryPG) tx(ctx context.Context) *gorm.DB {
	return database.TxFromCtx(ctx, r.db).WithContext(ctx)
}

func (r *BookingRepositoryPG) Save(ctx context.Context, booking *mealdomain.Booking) error {
	model := BookingModel{
		UserID:          booking.UserID,
		MealID:          booking.MealID,
		GarnishOptionID: booking.GarnishOptionID,
	}
	if err := r.tx(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("save booking: %w", err)
	}
	booking.ID = model.ID
	booking.CreatedAt = model.CreatedAt
	return nil
}

func (r *BookingRepositoryPG) FindByID(ctx context.Context, id string) (*mealdomain.Booking, error) {
	var model BookingModel
	err := r.tx(ctx).Preload("Meal.Template").Preload("GarnishOption").Where("id = ?", id).First(&model).Error
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
	if err := r.tx(ctx).Model(&BookingModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	var models []BookingModel
	err := r.tx(ctx).
		Preload("Meal.Template").Preload("GarnishOption").
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

func (r *BookingRepositoryPG) FindByUserAndDate(ctx context.Context, userID string, date time.Time) (*mealdomain.Booking, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var model BookingModel
	err := r.tx(ctx).
		Select("meal_bookings.*").
		Joins("JOIN meals ON meals.id = meal_bookings.meal_id").
		Where("meal_bookings.user_id = ?", userID).
		Where("meals.date >= ? AND meals.date < ?", start, end).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, mealdomain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("find booking by user/date: %w", err)
	}
	return toDomainBooking(&model), nil
}

func (r *BookingRepositoryPG) Delete(ctx context.Context, id string) error {
	err := r.tx(ctx).Where("id = ?", id).Delete(&BookingModel{}).Error
	if err != nil {
		return fmt.Errorf("delete booking: %w", err)
	}
	return nil
}

func toDomainBooking(model *BookingModel) *mealdomain.Booking {
	b := &mealdomain.Booking{
		ID:              model.ID,
		UserID:          model.UserID,
		MealID:          model.MealID,
		GarnishOptionID: model.GarnishOptionID,
		CreatedAt:       model.CreatedAt,
	}
	if model.Meal.ID != "" {
		b.Meal = toDomainMeal(&model.Meal)
	}
	if model.GarnishOption != nil {
		b.GarnishOption = &mealdomain.GarnishOption{
			ID:         model.GarnishOption.ID,
			TemplateID: model.GarnishOption.TemplateID,
			Name:       model.GarnishOption.Name,
		}
	}
	return b
}

type dailySummaryRow struct {
	ProjectID   string
	ProjectName string
	MealID      string
	MealTitle   string
	UserName    string
}

func (r *BookingRepositoryPG) GetDailySummary(ctx context.Context, date time.Time, projectID string) (*mealdomain.DailySummary, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	var rows []dailySummaryRow
	q := r.db.WithContext(ctx).
		Table("meal_bookings b").
		Select("m.project_id AS project_id, p.name AS project_name, m.id AS meal_id, mt.title AS meal_title, u.name AS user_name").
		Joins("JOIN meals m ON m.id = b.meal_id").
		Joins("JOIN projects p ON p.id = m.project_id").
		Joins("JOIN meal_templates mt ON mt.id = m.template_id").
		Joins("JOIN users u ON u.id = b.user_id").
		Where("m.date >= ? AND m.date < ?", start, end)
	if projectID != "" {
		q = q.Where("m.project_id = ?", projectID)
	}
	err := q.Order("p.name, m.id, u.name").Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("get daily summary: %w", err)
	}

	// Agrupar por proyecto, luego por vianda
	projectIndex := make(map[string]int)
	mealIndex := make(map[string]int) // clave: "projectID|mealID"
	projects := make([]mealdomain.ProjectDailySummary, 0)

	for _, row := range rows {
		pIdx, pExists := projectIndex[row.ProjectID]
		if !pExists {
			projects = append(projects, mealdomain.ProjectDailySummary{
				ProjectID:     row.ProjectID,
				ProjectName:   row.ProjectName,
				MealSummaries: []mealdomain.MealDailySummary{},
			})
			pIdx = len(projects) - 1
			projectIndex[row.ProjectID] = pIdx
		}

		mealKey := row.ProjectID + "|" + row.MealID
		mIdx, mExists := mealIndex[mealKey]
		if !mExists {
			projects[pIdx].MealSummaries = append(projects[pIdx].MealSummaries, mealdomain.MealDailySummary{
				MealID:  row.MealID,
				Title:   row.MealTitle,
				Persons: []mealdomain.PersonSummary{},
			})
			mIdx = len(projects[pIdx].MealSummaries) - 1
			mealIndex[mealKey] = mIdx
		}
		projects[pIdx].MealSummaries[mIdx].Quantity++
		projects[pIdx].MealSummaries[mIdx].Persons = append(
			projects[pIdx].MealSummaries[mIdx].Persons,
			mealdomain.PersonSummary{Name: row.UserName},
		)
		projects[pIdx].TotalMenus++
	}

	total := 0
	for _, p := range projects {
		total += p.TotalMenus
	}

	return &mealdomain.DailySummary{
		Date:       start,
		TotalMenus: total,
		Projects:   projects,
	}, nil
}
