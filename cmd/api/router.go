package main

import (
	"net/http"

	apideps "macabi-back/cmd/api/deps"
	eventhttp "macabi-back/internal/event/infrastructure/http"
	expenseshttp "macabi-back/internal/expenses/infrastructure/http"
	projecthttp "macabi-back/internal/project/infrastructure/http"
	"macabi-back/internal/shared/middleware"
	stockhttp "macabi-back/internal/stock/infrastructure/http"
	userhttp "macabi-back/internal/user/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d *apideps.Dependencies) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	userhttp.RegisterRoutes(r, d.UserHandlers, d.TokenPrv)
	projecthttp.RegisterRoutes(r, d.ProjectHandler, d.TokenPrv)
	eventhttp.RegisterRoutes(r, d.EventHandler, d.TokenPrv)
	stockhttp.RegisterRoutes(r, d.StockHandler, d.TokenPrv)
	expenseshttp.RegisterRoutes(r, d.ExpensesHandler, d.TokenPrv)

	return r
}
