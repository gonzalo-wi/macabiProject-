package expensespersistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	expensesdomain "macabi-back/internal/expenses/domain"
)

type ExpenseCategoryModel struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	CreatedAt time.Time
}

func (ExpenseCategoryModel) TableName() string { return "expense_categories" }

func (r *ExpenseRepositoryPG) SaveCategory(ctx context.Context, c *expensesdomain.ExpenseCategory) error {
	m := ExpenseCategoryModel{Name: strings.TrimSpace(c.Name)}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	c.ID = m.ID
	c.CreatedAt = m.CreatedAt
	return nil
}

func (r *ExpenseRepositoryPG) ListCategories(ctx context.Context) ([]expensesdomain.ExpenseCategory, error) {
	var models []ExpenseCategoryModel
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]expensesdomain.ExpenseCategory, len(models))
	for i, m := range models {
		out[i] = toDomainCategory(m)
	}
	return out, nil
}

func (r *ExpenseRepositoryPG) FindCategoryByID(ctx context.Context, id string) (*expensesdomain.ExpenseCategory, error) {
	var m ExpenseCategoryModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, expensesdomain.ErrCategoryNotFound
		}
		return nil, err
	}
	cat := toDomainCategory(m)
	return &cat, nil
}

func (r *ExpenseRepositoryPG) DeleteCategory(ctx context.Context, id string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ExpenseModel{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return expensesdomain.ErrCategoryInUse
	}
	res := r.db.WithContext(ctx).Delete(&ExpenseCategoryModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return expensesdomain.ErrCategoryNotFound
	}
	return nil
}

func toDomainCategory(m ExpenseCategoryModel) expensesdomain.ExpenseCategory {
	return expensesdomain.ExpenseCategory{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
	}
}
