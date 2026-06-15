package deps

import (
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/notifications"
	stockpersistence "macabi-back/internal/stock/infrastructure/persistence"
	stockpush "macabi-back/internal/stock/infrastructure/push"

	"gorm.io/gorm"
)

func buildPushNotifier(db *gorm.DB, cfg *config.Config) notifications.PushNotifier {
	if cfg.VAPIDPublicKey != "" && cfg.VAPIDPrivateKey != "" && cfg.VAPIDSubject != "" {
		pushSubRepo := stockpersistence.NewPushSubscriptionRepositoryPG(db)
		return stockpush.NewWebPushNotifier(pushSubRepo, cfg.VAPIDPublicKey, cfg.VAPIDPrivateKey, cfg.VAPIDSubject)
	}
	return notifications.NoOpPush{}
}
