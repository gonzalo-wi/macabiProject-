package main

import (
	"net/http"

	mealhttp "macabi-back/internal/meal/infrastructure/http"
	projecthttp "macabi-back/internal/project/infrastructure/http"
	"macabi-back/internal/shared/middleware"
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
	mealhttp.RegisterRoutes(r, deps.MealHandler, deps.BookingHandler, deps.MealTemplateHandler, deps.TokenPrv)
	projecthttp.RegisterRoutes(r, deps.ProjectHandler, deps.TokenPrv)

	return r
}
