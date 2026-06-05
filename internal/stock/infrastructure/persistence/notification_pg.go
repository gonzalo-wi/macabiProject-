package stockpersistence

import (
	"context"
	"time"

	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type NotificationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null"`
	RequestID string `gorm:"type:uuid;not null"`
	Message   string `gorm:"not null"`
	ReadAt    *time.Time
	CreatedAt time.Time
}

func (NotificationModel) TableName() string { return "stock_notifications" }

func (r *StockRepositoryPG) SaveNotification(ctx context.Context, n *stockdomain.StockNotification) error {
	m := NotificationModel{
		UserID:    n.UserID,
		RequestID: n.RequestID,
		Message:   n.Message,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	return nil
}

func (r *StockRepositoryPG) ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.StockNotification], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&NotificationModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return pagination.Result[stockdomain.StockNotification]{}, err
	}
	var models []NotificationModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error; err != nil {
		return pagination.Result[stockdomain.StockNotification]{}, err
	}
	out := make([]stockdomain.StockNotification, len(models))
	for i, m := range models {
		out[i] = toDomainNotification(m)
	}
	return pagination.NewResult(out, total, params), nil
}

func (r *StockRepositoryPG) MarkNotificationRead(ctx context.Context, id, userID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return stockdomain.ErrNotificationNotFound
	}
	return nil
}

func (r *StockRepositoryPG) MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", now)
	return res.RowsAffected, res.Error
}

func (r *StockRepositoryPG) UnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func toDomainNotification(m NotificationModel) stockdomain.StockNotification {
	return stockdomain.StockNotification{
		ID:        m.ID,
		UserID:    m.UserID,
		RequestID: m.RequestID,
		Message:   m.Message,
		ReadAt:    m.ReadAt,
		CreatedAt: m.CreatedAt,
	}
}
