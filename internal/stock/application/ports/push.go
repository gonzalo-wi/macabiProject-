package stockports

import "context"

// StoredPushSubscription is a Web Push API subscription saved for a user.
type StoredPushSubscription struct {
	ID       string
	UserID   string
	Endpoint string
	P256DH   string
	Auth     string
}

// PushSubscriptionRepository persists Web Push subscriptions.
type PushSubscriptionRepository interface {
	Upsert(ctx context.Context, userID, endpoint, p256dh, auth string) error
	DeleteByEndpoint(ctx context.Context, endpoint string) error
	FindByUserID(ctx context.Context, userID string) ([]StoredPushSubscription, error)
}

// UserPushNotifier sends Web Push notifications to all devices of a user.
// Implementations handle subscription lookup and stale-subscription cleanup.
type UserPushNotifier interface {
	Notify(ctx context.Context, userID string, title, body, actionURL string)
}
