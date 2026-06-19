package newspersistence

import (
	"context"
	"time"

	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"
)

type NewsNotificationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null;index"`
	NewsID    string `gorm:"type:uuid;not null;index"`
	Message   string `gorm:"not null"`
	ReadAt    *time.Time
	CreatedAt time.Time
}

func (NewsNotificationModel) TableName() string { return "news_notifications" }

func (r *NewsRepositoryPG) SaveNotifications(ctx context.Context, ns []*newsdomain.NewsNotification) error {
	if len(ns) == 0 {
		return nil
	}
	models := make([]NewsNotificationModel, len(ns))
	for i, n := range ns {
		models[i] = NewsNotificationModel{
			UserID:  n.UserID,
			NewsID:  n.NewsID,
			Message: n.Message,
		}
	}
	return r.db.WithContext(ctx).CreateInBatches(&models, 500).Error
}

func (r *NewsRepositoryPG) ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[newsdomain.NewsNotification], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&NewsNotificationModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return pagination.Result[newsdomain.NewsNotification]{}, err
	}
	var models []NewsNotificationModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error; err != nil {
		return pagination.Result[newsdomain.NewsNotification]{}, err
	}
	out := make([]newsdomain.NewsNotification, len(models))
	for i, m := range models {
		out[i] = toDomainNewsNotification(m)
	}
	return pagination.NewResult(out, total, params), nil
}

func (r *NewsRepositoryPG) MarkNotificationRead(ctx context.Context, id, userID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&NewsNotificationModel{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return newsdomain.ErrNotificationNotFound
	}
	return nil
}

func (r *NewsRepositoryPG) MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&NewsNotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", now)
	return res.RowsAffected, res.Error
}

func (r *NewsRepositoryPG) UnreadNotificationCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&NewsNotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *NewsRepositoryPG) DeleteNotificationsByNewsID(ctx context.Context, newsID string) error {
	return r.db.WithContext(ctx).Where("news_id = ?", newsID).Delete(&NewsNotificationModel{}).Error
}

func toDomainNewsNotification(m NewsNotificationModel) newsdomain.NewsNotification {
	return newsdomain.NewsNotification{
		ID:        m.ID,
		UserID:    m.UserID,
		NewsID:    m.NewsID,
		Message:   m.Message,
		ReadAt:    m.ReadAt,
		CreatedAt: m.CreatedAt,
	}
}
