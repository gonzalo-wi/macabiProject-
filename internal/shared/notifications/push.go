package notifications

import "context"

// PushNotifier envía Web Push al navegador del usuario (VAPID). NoOp si no hay keys configuradas.
type PushNotifier interface {
	Notify(ctx context.Context, userID string, title, body, actionURL string)
}

type NoOpPush struct{}

func (NoOpPush) Notify(_ context.Context, _ string, _, _, _ string) {}

// PushToUser envía push nativa; ignora si notifier es nil.
func PushToUser(ctx context.Context, notifier PushNotifier, userID, title, body, actionURL string) {
	if notifier == nil {
		return
	}
	notifier.Notify(ctx, userID, title, body, actionURL)
}
