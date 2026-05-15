package stockhttp

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

		api.GET("/stock/resources", h.ListResources)
		api.GET("/stock/resources/:id", h.GetResource)
		api.POST("/stock/resources",
			userhttp.RequireRole(userdomain.RoleAdmin),
			h.CreateResource,
		)
		api.PUT("/stock/resources/:id",
			userhttp.RequireRole(userdomain.RoleAdmin),
			h.UpdateResource,
		)
		api.DELETE("/stock/resources/:id",
			userhttp.RequireRole(userdomain.RoleAdmin),
			h.DeleteResource,
		)

		api.POST("/stock/requests", h.CreateRequest)
		api.GET("/stock/requests/my", h.ListMyRequests)
		api.GET("/stock/requests/:id", h.GetRequestDetail)
		api.GET("/stock/requests",
			userhttp.RequireRole(userdomain.RoleAdmin),
			h.ListRequests,
		)
		api.PATCH("/stock/requests/:id/approve", h.ApproveRequest)
		api.PATCH("/stock/requests/:id/reject", h.RejectRequest)
		api.PATCH("/stock/requests/:id/deliver", h.DeliverRequest)
		api.PATCH("/stock/requests/:id/return", h.ReturnRequest)

		api.GET("/stock/notifications/unread-count", h.UnreadCount)
		api.GET("/stock/notifications", h.ListNotifications)
		api.PATCH("/stock/notifications/:id/read", h.MarkNotificationRead)
	}
}
