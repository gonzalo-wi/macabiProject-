package userhttp

import (
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authHandler *AuthHandler, userHandler *UserHandler, tokenPrv userports.TokenProvider) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(AuthMiddleware(tokenPrv))
	{
		api.GET("/me", userHandler.Me)
		api.GET("/users",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.ListUsers,
		)
		api.PATCH("/users/:id/role",
			RequireRole(userdomain.RoleSuperAdmin),
			userHandler.ChangeRole,
		)
	}
}
