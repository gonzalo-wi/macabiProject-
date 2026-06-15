package eventpersistence

import (
	"context"
	"time"

	eventdomain "macabi-back/internal/event/domain"
	"macabi-back/internal/shared/pagination"
)

type EventNotificationModel struct {
	ID              string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          string `gorm:"type:uuid;not null;index"`
	EventInstanceID string `gorm:"type:uuid;not null;index:idx_event_notif_user_event"`
	Message         string `gorm:"not null"`
	ReadAt          *time.Time
	CreatedAt       time.Time
}

func (EventNotificationModel) TableName() string { return "event_notifications" }

func (r *RepositoryPG) SaveNotification(ctx context.Context, n *eventdomain.EventNotification) error {
	m := EventNotificationModel{
		UserID:          n.UserID,
		EventInstanceID: n.EventInstanceID,
		Message:         n.Message,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	return nil
}

func (r *RepositoryPG) ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[eventdomain.EventNotification], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&EventNotificationModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return pagination.Result[eventdomain.EventNotification]{}, err
	}
	var models []EventNotificationModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error; err != nil {
		return pagination.Result[eventdomain.EventNotification]{}, err
	}
	out := make([]eventdomain.EventNotification, len(models))
	for i, m := range models {
		out[i] = toDomainEventNotification(m)
	}
	return pagination.NewResult(out, total, params), nil
}

func (r *RepositoryPG) MarkNotificationRead(ctx context.Context, id, userID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&EventNotificationModel{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return eventdomain.ErrNotificationNotFound
	}
	return nil
}

func (r *RepositoryPG) MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&EventNotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", now)
	return res.RowsAffected, res.Error
}

func (r *RepositoryPG) UnreadNotificationCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&EventNotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *RepositoryPG) HasReminderForEvent(ctx context.Context, userID, eventInstanceID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&EventNotificationModel{}).
		Where("user_id = ? AND event_instance_id = ?", userID, eventInstanceID).
		Count(&count).Error
	return count > 0, err
}

func toDomainEventNotification(m EventNotificationModel) eventdomain.EventNotification {
	return eventdomain.EventNotification{
		ID:              m.ID,
		UserID:          m.UserID,
		EventInstanceID: m.EventInstanceID,
		Message:         m.Message,
		ReadAt:          m.ReadAt,
		CreatedAt:       m.CreatedAt,
	}
}
