package main

import (
	"net/http"

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

	return r
}
