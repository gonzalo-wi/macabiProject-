package expensespersistence

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProjectExpenseBudgetModel struct {
	ProjectID     string          `gorm:"type:uuid;primaryKey"`
	MonthlyAmount decimal.Decimal `gorm:"type:numeric(16,4);not null"`
	UpdatedAt     time.Time
}

func (ProjectExpenseBudgetModel) TableName() string { return "project_expense_budgets" }

func (r *ExpenseRepositoryPG) GetBudget(ctx context.Context, projectID string) (*decimal.Decimal, error) {
	var m ProjectExpenseBudgetModel
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).First(&m).Error
	if err != nil {
	 if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	amt := m.MonthlyAmount
	return &amt, nil
}

func (r *ExpenseRepositoryPG) SetBudget(ctx context.Context, projectID string, amount *decimal.Decimal) error {
	if amount == nil {
		return r.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&ProjectExpenseBudgetModel{}).Error
	}
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&ProjectExpenseBudgetModel{}).
		Where("project_id = ?", projectID).
		Updates(map[string]any{"monthly_amount": *amount, "updated_at": now})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return r.db.WithContext(ctx).Create(&ProjectExpenseBudgetModel{
			ProjectID: projectID, MonthlyAmount: *amount, UpdatedAt: now,
		}).Error
	}
	return nil
}
