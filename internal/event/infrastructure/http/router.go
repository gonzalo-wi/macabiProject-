package eventhttp

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
		api.GET("/event-instances/notifications/unread-count", h.UnreadNotificationCount)
		api.GET("/event-instances/notifications", h.ListNotifications)
		api.PATCH("/event-instances/notifications/read-all", h.MarkAllNotificationsRead)
		api.PATCH("/event-instances/notifications/:id/read", h.MarkNotificationRead)
		api.GET("/event-instances", h.ListEvents)
		api.GET("/event-instances/:id", h.GetEventDetail)
		api.GET("/event-instances/:id/responses/me", h.GetMyResponse)
		api.POST("/event-instances/:id/responses", h.SubmitResponse)

		admin := api.Group("")
		admin.Use(userhttp.RequireRole(userdomain.RoleAdmin))
		{
			admin.GET("/event-instances/:id/participant-responses", h.ListEventResponsesForAdmin)
			admin.GET("/event-instances/:id/responses", h.ListEventResponsesForAdmin)
			admin.GET("/event-modules/:id/response-summary", h.GetModuleResponseSummary)
			admin.POST("/event-instances", h.CreateEvent)
			admin.PATCH("/event-instances/:id", h.PatchEvent)
			admin.DELETE("/event-instances/:id", h.DeleteEvent)
			admin.PUT("/event-instances/:id/projects", h.SetEventProjects)
			admin.POST("/event-modules", h.CreateModule)
			admin.PATCH("/event-modules/:id", h.PatchModule)
			admin.DELETE("/event-modules/:id", h.DeleteModule)
			admin.PUT("/event-modules/:id/projects", h.SetModuleProjects)
			admin.POST("/event-option-groups", h.CreateOptionGroup)
			admin.PATCH("/event-option-groups/:id", h.PatchOptionGroup)
			admin.DELETE("/event-option-groups/:id", h.DeleteOptionGroup)
			admin.POST("/event-options", h.CreateOption)
			admin.PATCH("/event-options/:id", h.PatchOption)
			admin.DELETE("/event-options/:id", h.DeleteOption)
		}
	}
}
