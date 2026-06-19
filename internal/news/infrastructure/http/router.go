package newshttp

import (
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
	userhttp "macabi-back/internal/user/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, tokenPrv userports.TokenProvider) {
	api := r.Group("/api")
	api.Use(userhttp.AuthMiddleware(tokenPrv))
	{
		// Lectura (todos los autenticados: madrijim y coordinadores).
		api.GET("/news", h.ListPublished)
		api.GET("/news/latest", h.Latest)

		// Campana in-app.
		api.GET("/news/notifications/unread-count", h.UnreadNotificationCount)
		api.GET("/news/notifications", h.ListNotifications)
		api.PATCH("/news/notifications/read-all", h.MarkAllNotificationsRead)
		api.PATCH("/news/notifications/:id/read", h.MarkNotificationRead)

		// Administración.
		api.GET("/news/all", userhttp.RequireRole(userdomain.RoleAdmin), h.ListAll)
		api.POST("/news", userhttp.RequireRole(userdomain.RoleAdmin), h.Create)

		api.GET("/news/:id", h.Get)
		api.PATCH("/news/:id", userhttp.RequireRole(userdomain.RoleAdmin), h.Patch)
		api.DELETE("/news/:id", userhttp.RequireRole(userdomain.RoleAdmin), h.Delete)
		api.POST("/news/:id/image", userhttp.RequireRole(userdomain.RoleAdmin), h.ImageUpload)
		api.GET("/news/:id/image", h.ImageView)
		api.DELETE("/news/:id/image", userhttp.RequireRole(userdomain.RoleAdmin), h.RemoveImage)
	}
}
