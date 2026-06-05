package stockpush

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	webpush "github.com/SherClockHolmes/webpush-go"

	stockports "macabi-back/internal/stock/application/ports"
)

type pushPayload struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	ActionURL string `json:"url"`
}

type WebPushNotifier struct {
	subRepo         stockports.PushSubscriptionRepository
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string
}

func NewWebPushNotifier(
	subRepo stockports.PushSubscriptionRepository,
	vapidPublicKey, vapidPrivateKey, vapidSubject string,
) *WebPushNotifier {
	return &WebPushNotifier{
		subRepo:         subRepo,
		vapidPublicKey:  vapidPublicKey,
		vapidPrivateKey: vapidPrivateKey,
		vapidSubject:    vapidSubject,
	}
}

func (n *WebPushNotifier) Notify(ctx context.Context, userID string, title, body, actionURL string) {
	subs, err := n.subRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Printf("webpush: error fetching subscriptions for user=%s: %v", userID, err)
		return
	}
	if len(subs) == 0 {
		log.Printf("webpush: no subscriptions for user=%s, skipping push", userID)
		return
	}
	log.Printf("webpush: sending push to user=%s (%d subscription(s))", userID, len(subs))

	payload, err := json.Marshal(pushPayload{Title: title, Body: body, ActionURL: actionURL})
	if err != nil {
		return
	}

	for _, sub := range subs {
		resp, err := webpush.SendNotificationWithContext(ctx, payload, &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256DH,
				Auth:   sub.Auth,
			},
		}, &webpush.Options{
			Subscriber:      n.vapidSubject,
			VAPIDPublicKey:  n.vapidPublicKey,
			VAPIDPrivateKey: n.vapidPrivateKey,
			TTL:             86400,
		})
		if err != nil {
			log.Printf("webpush: send error (user=%s endpoint=%.40s...): %v", userID, sub.Endpoint, err)
			continue
		}
		resp.Body.Close()
		log.Printf("webpush: push sent to user=%s, status=%d", userID, resp.StatusCode)
		// 410 Gone = subscription expired; clean it up.
		if resp.StatusCode == http.StatusGone {
			_ = n.subRepo.DeleteByEndpoint(ctx, sub.Endpoint)
		}
	}
}

type NoOpPushNotifier struct{}

func (NoOpPushNotifier) Notify(_ context.Context, _ string, _, _, _ string) {}
