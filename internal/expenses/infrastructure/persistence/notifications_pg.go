package expensespersistence

import (
	"context"
	"time"

	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type ExpenseNotificationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null;index"`
	ExpenseID string `gorm:"type:uuid;not null"`
	ProjectID string `gorm:"type:uuid;not null"`
	Message   string `gorm:"not null"`
	ReadAt    *time.Time
	CreatedAt time.Time
}

func (ExpenseNotificationModel) TableName() string { return "expense_notifications" }

func (r *ExpenseRepositoryPG) SaveNotification(ctx context.Context, n *expensesdomain.ExpenseNotification) error {
	m := ExpenseNotificationModel{
		UserID:    n.UserID,
		ExpenseID: n.ExpenseID,
		ProjectID: n.ProjectID,
		Message:   n.Message,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	return nil
}

func (r *ExpenseRepositoryPG) ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseNotification], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ExpenseNotificationModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseNotification]{}, err
	}
	var models []ExpenseNotificationModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseNotification]{}, err
	}
	out := make([]expensesdomain.ExpenseNotification, len(models))
	for i, m := range models {
		out[i] = toDomainExpenseNotification(m)
	}
	return pagination.NewResult(out, total, params), nil
}

func (r *ExpenseRepositoryPG) MarkNotificationRead(ctx context.Context, id, userID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&ExpenseNotificationModel{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return expensesdomain.ErrNotificationNotFound
	}
	return nil
}

func (r *ExpenseRepositoryPG) MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&ExpenseNotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", now)
	return res.RowsAffected, res.Error
}

func (r *ExpenseRepositoryPG) UnreadNotificationCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ExpenseNotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *ExpenseRepositoryPG) DeleteNotificationsByExpenseID(ctx context.Context, expenseID string) error {
	return r.db.WithContext(ctx).Where("expense_id = ?", expenseID).Delete(&ExpenseNotificationModel{}).Error
}

func toDomainExpenseNotification(m ExpenseNotificationModel) expensesdomain.ExpenseNotification {
	return expensesdomain.ExpenseNotification{
		ID:        m.ID,
		UserID:    m.UserID,
		ExpenseID: m.ExpenseID,
		ProjectID: m.ProjectID,
		Message:   m.Message,
		ReadAt:    m.ReadAt,
		CreatedAt: m.CreatedAt,
	}
}
