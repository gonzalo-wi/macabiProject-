package deps

import (
	eventusecases "macabi-back/internal/event/application/usecases"
	eventhttp "macabi-back/internal/event/infrastructure/http"
	expenseshttp "macabi-back/internal/expenses/infrastructure/http"
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
	TokenPrv           userports.TokenProvider
	SendEventReminders *eventusecases.SendEventReminders
}

func Build(db *gorm.DB, cfg *config.Config) *Dependencies {
	authHandler, userHandler, tokenPrv := buildUserDeps(db, cfg)
	projectHandler := buildProjectDeps(db)
	eventHandler, sendReminders := buildEventDeps(db, cfg)
	stockHandler, stockRepo := buildStockDeps(db, cfg)
	// stockRepo satisfies ProjectMembership, ProjectCoordinatorReader and UserEmailReader
	// because StockRepositoryPG implements all three interfaces.
	expensesHandler := buildExpensesDeps(db, cfg, stockRepo, stockRepo, stockRepo)

	return &Dependencies{
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		ProjectHandler:     projectHandler,
		EventHandler:       eventHandler,
		StockHandler:       stockHandler,
		ExpensesHandler:    expensesHandler,
		TokenPrv:           tokenPrv,
		SendEventReminders: sendReminders,
	}
}
