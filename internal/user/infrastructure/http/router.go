package userhttp

import (
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handlers, tokenPrv userports.TokenProvider) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Auth.RegisterDisabled)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/forgot-password", h.Auth.ForgotPassword)
		auth.POST("/reset-password", h.Auth.ConfirmPasswordReset)
		auth.POST("/accept-invitation", h.Auth.AcceptInvitation)
	}

	api := r.Group("/api")
	api.Use(AuthMiddleware(tokenPrv))
	{
		api.GET("/me", h.User.Me)
		api.PATCH("/me/password", h.User.ChangePassword)
		api.GET("/users/invitations",
			RequireRole(userdomain.RoleAdmin),
			h.Invitation.ListPending,
		)
		api.POST("/users/invitations/:id/resend",
			RequireRole(userdomain.RoleAdmin),
			h.Invitation.Resend,
		)
		api.DELETE("/users/invitations/:id",
			RequireRole(userdomain.RoleAdmin),
			h.Invitation.Revoke,
		)
		api.POST("/users/invitations",
			RequireRole(userdomain.RoleAdmin),
			h.Invitation.Create,
		)
		api.GET("/users",
			RequireRole(userdomain.RoleAdmin),
			h.User.ListUsers,
		)
		api.PATCH("/users/:id/role",
			RequireRole(userdomain.RoleAdmin),
			h.User.ChangeRole,
		)
		api.PATCH("/users/:id/status",
			RequireRole(userdomain.RoleAdmin),
			h.User.SetStatus,
		)
		api.PUT("/users/:id",
			RequireRole(userdomain.RoleAdmin),
			h.User.UpdateUser,
		)
	}
}
