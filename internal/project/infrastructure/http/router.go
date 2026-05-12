package projecthttp

import (
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
	userhttp "macabi-back/internal/user/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, projectHandler *ProjectHandler, tokenPrv userports.TokenProvider) {
	api := r.Group("/api")
	api.Use(userhttp.AuthMiddleware(tokenPrv))
	{
		api.GET("/projects", projectHandler.List)
		api.GET("/projects/:id", projectHandler.Get)
		api.POST("/projects",
			userhttp.RequireRole(userdomain.RoleAdmin),
			projectHandler.Create,
		)
		api.PUT("/projects/:id",
			userhttp.RequireRole(userdomain.RoleAdmin),
			projectHandler.Update,
		)
		api.DELETE("/projects/:id",
			userhttp.RequireRole(userdomain.RoleAdmin),
			projectHandler.Delete,
		)
		api.GET("/projects/:id/members", projectHandler.ListMembers)
		api.POST("/projects/:id/members",
			userhttp.RequireRole(userdomain.RoleAdmin),
			projectHandler.AddMember,
		)
		api.DELETE("/projects/:id/members/:userId",
			userhttp.RequireRole(userdomain.RoleAdmin),
			projectHandler.RemoveMember,
		)
	}
}
