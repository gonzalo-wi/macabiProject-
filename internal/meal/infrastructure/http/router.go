package mealhttp

import (
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"
	userhttp "macabi-back/internal/user/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, mealHandler *MealHandler, bookingHandler *BookingHandler, templateHandler *MealTemplateHandler, tokenPrv userports.TokenProvider) {
	api := r.Group("/api")
	api.Use(userhttp.AuthMiddleware(tokenPrv))
	{
		// Templates (recetas reutilizables)
		api.GET("/meal-templates", templateHandler.List)
		api.POST("/meal-templates",
			userhttp.RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			templateHandler.Create,
		)
		api.PUT("/meal-templates/:id",
			userhttp.RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			templateHandler.Update,
		)
		api.DELETE("/meal-templates/:id",
			userhttp.RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			templateHandler.Delete,
		)

		// Meals (asignación de template a fecha + cantidad)
		api.GET("/meals", mealHandler.ListByDate)
		api.POST("/meals",
			userhttp.RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			mealHandler.Create,
		)
		api.DELETE("/meals/:id",
			userhttp.RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			mealHandler.Delete,
		)

		// Bookings
		api.POST("/bookings", bookingHandler.Book)
		api.GET("/bookings/mine", bookingHandler.ListMine)
		api.DELETE("/bookings/:id", bookingHandler.Cancel)

		api.GET("/admin/bookings/daily-summary",
			userhttp.RequireRole(userdomain.RoleSuperAdmin, userdomain.RoleAdmin),
			bookingHandler.DailySummary,
		)
	}
}
