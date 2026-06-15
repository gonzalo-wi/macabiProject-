package stockports

import (
	"context"

	"macabi-back/internal/shared/notifications"
)

type StoredPushSubscription struct {
	ID       string
	UserID   string
	Endpoint string
	P256DH   string
	Auth     string
}

type PushSubscriptionRepository interface {
	Upsert(ctx context.Context, userID, endpoint, p256dh, auth string) error
	DeleteByEndpoint(ctx context.Context, endpoint string) error
	FindByUserID(ctx context.Context, userID string) ([]StoredPushSubscription, error)
}

// UserPushNotifier alias del contrato compartido de Web Push.
type UserPushNotifier = notifications.PushNotifier
