package attendancehttp

import (
	"github.com/gin-gonic/gin"

	userports "macabi-back/internal/user/application/ports"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func RegisterRoutes(r *gin.Engine, handler *AttendanceHandler, tokenPrv userports.TokenProvider) {
	api := r.Group("/api")
	api.Use(userhttp.AuthMiddleware(tokenPrv))
	{
		api.POST("/projects/:id/attendance", handler.Confirm)
		api.GET("/projects/:id/attendance", handler.GetCount)
	}
}