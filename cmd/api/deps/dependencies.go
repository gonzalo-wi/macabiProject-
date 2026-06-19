package deps

import (
	eventusecases "macabi-back/internal/event/application/usecases"
	eventhttp "macabi-back/internal/event/infrastructure/http"
	expenseshttp "macabi-back/internal/expenses/infrastructure/http"
	newshttp "macabi-back/internal/news/infrastructure/http"
	projecthttp "macabi-back/internal/project/infrastructure/http"
	"macabi-back/internal/shared/config"
	stockhttp "macabi-back/internal/stock/infrastructure/http"
	userports "macabi-back/internal/user/application/ports"
	userhttp "macabi-back/internal/user/infrastructure/http"

	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler        *userhttp.AuthHandler
	UserHandler        *userhttp.UserHandler
	ProjectHandler     *projecthttp.ProjectHandler
	EventHandler       *eventhttp.Handler
	StockHandler       *stockhttp.Handler
	ExpensesHandler    *expenseshttp.Handler
	NewsHandler        *newshttp.Handler
	TokenPrv           userports.TokenProvider
	SendEventReminders *eventusecases.SendEventReminders
}

func Build(db *gorm.DB, cfg *config.Config) *Dependencies {
	push := buildPushNotifier(db, cfg)

	authHandler, userHandler, tokenPrv, memberDir := buildUserDeps(db, cfg)
	projectHandler := buildProjectDeps(db)
	eventHandler, sendReminders := buildEventDeps(db, cfg, push)
	stockHandler, stockRepo := buildStockDeps(db, cfg, push)
	expensesHandler := buildExpensesDeps(db, cfg, stockRepo, stockRepo, stockRepo, push)
	newsHandler := buildNewsDeps(db, cfg, memberDir, stockRepo, push)

	return &Dependencies{
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		ProjectHandler:     projectHandler,
		EventHandler:       eventHandler,
		StockHandler:       stockHandler,
		ExpensesHandler:    expensesHandler,
		NewsHandler:        newsHandler,
		TokenPrv:           tokenPrv,
		SendEventReminders: sendReminders,
	}
}
