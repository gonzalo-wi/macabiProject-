package deps

import (
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/notifications"
	stockservices "macabi-back/internal/stock/application/services"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockhttp "macabi-back/internal/stock/infrastructure/http"
	stockmail "macabi-back/internal/stock/infrastructure/mail"
	stockpersistence "macabi-back/internal/stock/infrastructure/persistence"

	"gorm.io/gorm"
)

func buildStockDeps(db *gorm.DB, cfg *config.Config, push notifications.PushNotifier) (*stockhttp.Handler, *stockpersistence.StockRepositoryPG) {
	stockRepo := stockpersistence.NewStockRepositoryPG(db)
	if err := stockpersistence.RunMigrations(db); err != nil {
		panic("stock migrations failed: " + err.Error())
	}

	stockMailer := stockmail.NewBrevoStockMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom, cfg.FrontendPublicURL)
	pushSubRepo := stockpersistence.NewPushSubscriptionRepositoryPG(db)
	stockNotifier := stockservices.NewStockNotificationService(stockRepo, stockRepo, stockRepo, stockMailer, push)

	stockHandler := stockhttp.NewHandler(
		stockusecases.NewCreateResource(stockRepo),
		stockusecases.NewListResources(stockRepo),
		stockusecases.NewGetResource(stockRepo),
		stockusecases.NewUpdateResource(stockRepo),
		stockusecases.NewDeleteResource(stockRepo),
		stockusecases.NewCreateRequest(stockRepo, stockRepo, stockNotifier),
		stockusecases.NewApproveRequest(stockRepo, stockRepo, stockNotifier),
		stockusecases.NewRejectRequest(stockRepo, stockRepo, stockNotifier),
		stockusecases.NewCancelRequest(stockRepo),
		stockusecases.NewDeliverRequest(stockRepo, stockRepo),
		stockusecases.NewReturnRequest(stockRepo, stockRepo),
		stockusecases.NewListRequests(stockRepo, stockRepo),
		stockusecases.NewListMyRequests(stockRepo),
		stockusecases.NewGetRequestDetail(stockRepo, stockRepo),
		stockusecases.NewListNotifications(stockRepo),
		stockusecases.NewMarkNotificationRead(stockRepo),
		stockusecases.NewMarkAllNotificationsRead(stockRepo),
		stockusecases.NewUnreadNotificationCount(stockRepo),
		stockusecases.NewRegisterPushSubscription(pushSubRepo),
		stockusecases.NewUnregisterPushSubscription(pushSubRepo),
		cfg.VAPIDPublicKey,
	)

	return stockHandler, stockRepo
}
