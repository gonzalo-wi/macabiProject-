package stockports

import "context"

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

type UserPushNotifier interface {
	Notify(ctx context.Context, userID string, title, body, actionURL string)
}
