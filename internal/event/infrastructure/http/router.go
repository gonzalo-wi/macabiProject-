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
		api.GET("/event-instances", h.ListEvents)
		api.GET("/event-instances/:id", h.GetEventDetail)
		api.GET("/event-instances/:id/responses/me", h.GetMyResponse)
		api.POST("/event-instances/:id/responses", h.SubmitResponse)

		admin := api.Group("")
		admin.Use(userhttp.RequireRole(userdomain.RoleAdmin))
		{
			admin.POST("/event-instances", h.CreateEvent)
			admin.PATCH("/event-instances/:id", h.PatchEvent)
			admin.PUT("/event-instances/:id/projects", h.SetEventProjects)
			admin.POST("/event-modules", h.CreateModule)
			admin.PATCH("/event-modules/:id", h.PatchModule)
			admin.DELETE("/event-modules/:id", h.DeleteModule)
			admin.PUT("/event-modules/:id/projects", h.SetModuleProjects)
			admin.POST("/event-option-groups", h.CreateOptionGroup)
			admin.PATCH("/event-option-groups/:id", h.PatchOptionGroup)
			admin.POST("/event-options", h.CreateOption)
			admin.PATCH("/event-options/:id", h.PatchOption)
		}
	}
}
