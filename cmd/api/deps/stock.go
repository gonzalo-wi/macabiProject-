package deps

import (
	"macabi-back/internal/shared/config"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockhttp "macabi-back/internal/stock/infrastructure/http"
	stockmail "macabi-back/internal/stock/infrastructure/mail"
	stockpersistence "macabi-back/internal/stock/infrastructure/persistence"

	"gorm.io/gorm"
)

func buildStockDeps(db *gorm.DB, cfg *config.Config) (*stockhttp.Handler, *stockpersistence.StockRepositoryPG) {
	stockRepo := stockpersistence.NewStockRepositoryPG(db)
	if err := stockpersistence.RunMigrations(db); err != nil {
		panic("stock migrations failed: " + err.Error())
	}

	stockMailer := stockmail.NewBrevoStockMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom, cfg.FrontendPublicURL)

	stockHandler := stockhttp.NewHandler(
		stockusecases.NewCreateResource(stockRepo),
		stockusecases.NewListResources(stockRepo),
		stockusecases.NewGetResource(stockRepo),
		stockusecases.NewUpdateResource(stockRepo),
		stockusecases.NewDeleteResource(stockRepo),
		stockusecases.NewCreateRequest(stockRepo, stockRepo, stockRepo, stockMailer),
		stockusecases.NewApproveRequest(stockRepo, stockRepo, stockRepo, stockMailer),
		stockusecases.NewRejectRequest(stockRepo, stockRepo, stockRepo, stockMailer),
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
	)

	return stockHandler, stockRepo
}
