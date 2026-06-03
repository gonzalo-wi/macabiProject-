package expenseshttp

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
		api.POST("/expenses", h.Create)
		api.GET("/expenses", userhttp.RequireRole(userdomain.RoleAdmin), h.ListAll)
		api.GET("/expenses/my", h.ListMine)
		api.GET("/expenses/notifications/unread-count", h.UnreadNotificationCount)
		api.GET("/expenses/notifications", h.ListNotifications)
		api.PATCH("/expenses/notifications/read-all", h.MarkAllNotificationsRead)
		api.PATCH("/expenses/notifications/:id/read", h.MarkNotificationRead)
		// Rutas estáticas de /expenses ANTES de la wildcard /:id (precedencia de Gin por orden).
		api.GET("/expenses/analytics", userhttp.RequireRole(userdomain.RoleAdmin), h.Analytics)
		api.GET("/expenses/:id", h.Get)
		api.PATCH("/expenses/:id", h.Patch)
		api.PATCH("/expenses/:id/approve", h.Approve)
		api.PATCH("/expenses/:id/reject", h.Reject)
		api.DELETE("/expenses/:id", h.Delete)
		api.POST("/expenses/:id/receipt-upload", h.ReceiptUpload)
		api.POST("/expenses/:id/receipt", h.ReceiptUploadMultipart)
		api.GET("/expenses/:id/receipt", h.ReceiptView)
		api.DELETE("/expenses/:id/receipt", h.RemoveReceipt)

		api.GET("/projects/:id/expenses/summary", h.SummaryByProject)
		api.GET("/projects/:id/expenses", h.ListByProject)
	}
}
