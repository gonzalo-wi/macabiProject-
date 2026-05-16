package stockpersistence

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	stockports "macabi-back/internal/stock/application/ports"
)

type PushSubscriptionModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null;index"`
	Endpoint  string `gorm:"type:text;not null;uniqueIndex"`
	P256DH    string `gorm:"column:p256dh;type:text;not null"`
	Auth      string `gorm:"type:text;not null"`
	CreatedAt time.Time
}

func (PushSubscriptionModel) TableName() string { return "push_subscriptions" }

type PushSubscriptionRepositoryPG struct {
	db *gorm.DB
}

var _ stockports.PushSubscriptionRepository = (*PushSubscriptionRepositoryPG)(nil)

func NewPushSubscriptionRepositoryPG(db *gorm.DB) *PushSubscriptionRepositoryPG {
	return &PushSubscriptionRepositoryPG{db: db}
}

func (r *PushSubscriptionRepositoryPG) Upsert(ctx context.Context, userID, endpoint, p256dh, auth string) error {
	m := PushSubscriptionModel{
		UserID:   userID,
		Endpoint: endpoint,
		P256DH:   p256dh,
		Auth:     auth,
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "endpoint"}},
			DoUpdates: clause.AssignmentColumns([]string{"user_id", "p256dh", "auth"}),
		}).
		Create(&m).Error
}

func (r *PushSubscriptionRepositoryPG) DeleteByEndpoint(ctx context.Context, endpoint string) error {
	return r.db.WithContext(ctx).
		Where("endpoint = ?", endpoint).
		Delete(&PushSubscriptionModel{}).Error
}

func (r *PushSubscriptionRepositoryPG) FindByUserID(ctx context.Context, userID string) ([]stockports.StoredPushSubscription, error) {
	var models []PushSubscriptionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]stockports.StoredPushSubscription, len(models))
	for i, m := range models {
		out[i] = stockports.StoredPushSubscription{
			ID:       m.ID,
			UserID:   m.UserID,
			Endpoint: m.Endpoint,
			P256DH:   m.P256DH,
			Auth:     m.Auth,
		}
	}
	return out, nil
}
