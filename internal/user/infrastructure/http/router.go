package userhttp

import (
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authHandler *AuthHandler, userHandler *UserHandler, tokenPrv userports.TokenProvider) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterDisabled)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ConfirmPasswordReset)
		auth.POST("/accept-invitation", authHandler.AcceptInvitation)
	}

	api := r.Group("/api")
	api.Use(AuthMiddleware(tokenPrv))
	{
		api.GET("/me", userHandler.Me)
		api.PATCH("/me/password", userHandler.ChangePassword)
		api.GET("/users/invitations",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.ListPendingInvitations,
		)
		api.POST("/users/invitations/:id/resend",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.ResendInvitation,
		)
		api.DELETE("/users/invitations/:id",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.RevokeInvitation,
		)
		api.POST("/users/invitations",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.CreateInvitation,
		)
		api.GET("/users",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.ListUsers,
		)
		api.PATCH("/users/:id/role",
			RequireRole(userdomain.RoleSuperAdmin),
			userHandler.ChangeRole,
		)
		api.PATCH("/users/:id/status",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.SetStatus,
		)
		api.PUT("/users/:id",
			RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			userHandler.UpdateUser,
		)
	}
}
