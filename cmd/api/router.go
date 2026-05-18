package main

import (
	"net/http"

	eventhttp "macabi-back/internal/event/infrastructure/http"
	expenseshttp "macabi-back/internal/expenses/infrastructure/http"
	projecthttp "macabi-back/internal/project/infrastructure/http"
	"macabi-back/internal/shared/middleware"
	stockhttp "macabi-back/internal/stock/infrastructure/http"
	userhttp "macabi-back/internal/user/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(deps *Dependencies) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	userhttp.RegisterRoutes(r, deps.AuthHandler, deps.UserHandler, deps.TokenPrv)
	projecthttp.RegisterRoutes(r, deps.ProjectHandler, deps.TokenPrv)
	eventhttp.RegisterRoutes(r, deps.EventHandler, deps.TokenPrv)
	stockhttp.RegisterRoutes(r, deps.StockHandler, deps.TokenPrv)
	expenseshttp.RegisterRoutes(r, deps.ExpensesHandler, deps.TokenPrv)

	return r
}
