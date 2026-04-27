package main

import (
	mealusecases "macabi-back/internal/meal/application/usecases"
	mealhttp "macabi-back/internal/meal/infrastructure/http"
	mealpersistence "macabi-back/internal/meal/infrastructure/persistence"
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/database"
	userports "macabi-back/internal/user/application/ports"
	userusecases "macabi-back/internal/user/application/usecases"
	userhttp "macabi-back/internal/user/infrastructure/http"
	userpersistence "macabi-back/internal/user/infrastructure/persistence"
	usersecurity "macabi-back/internal/user/infrastructure/security"

	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler         *userhttp.AuthHandler
	UserHandler         *userhttp.UserHandler
	MealHandler         *mealhttp.MealHandler
	BookingHandler      *mealhttp.BookingHandler
	MealTemplateHandler *mealhttp.MealTemplateHandler
	TokenPrv            userports.TokenProvider
}

func BuildDependencies(db *gorm.DB, cfg *config.Config) *Dependencies {
	// User infrastructure
	userRepo := userpersistence.NewUserRepositoryPG(db)
	hasher := usersecurity.NewBcryptHasher()
	jwtProvider := usersecurity.NewJWTProvider(cfg.JWTSecret, cfg.JWTExpiration)

	// User use cases
	registerUC := userusecases.NewRegisterUser(userRepo, hasher)
	loginUC := userusecases.NewLogin(userRepo, hasher, jwtProvider)
	getCurrentUserUC := userusecases.NewGetCurrentUser(userRepo)
	changeRoleUC := userusecases.NewChangeRole(userRepo)
	listUsersUC := userusecases.NewListUsers(userRepo)
	setUserStatusUC := userusecases.NewSetUserStatus(userRepo)
	updateUserUC := userusecases.NewUpdateUser(userRepo)
	changePasswordUC := userusecases.NewChangePassword(userRepo, hasher)

	// User handlers
	authHandler := userhttp.NewAuthHandler(registerUC, loginUC)
	userHandler := userhttp.NewUserHandler(getCurrentUserUC, changeRoleUC, listUsersUC, setUserStatusUC, updateUserUC, changePasswordUC)

	// Meal infrastructure
	mealRepo := mealpersistence.NewMealRepositoryPG(db)
	templateRepo := mealpersistence.NewMealTemplateRepositoryPG(db)
	bookingRepo := mealpersistence.NewBookingRepositoryPG(db)
	transactor := database.NewGORMTransactor(db)

	// Meal use cases
	createMealTemplateUC := mealusecases.NewCreateMealTemplate(templateRepo)
	listMealTemplatesUC := mealusecases.NewListMealTemplates(templateRepo)
	updateMealTemplateUC := mealusecases.NewUpdateMealTemplate(templateRepo)
	deleteMealTemplateUC := mealusecases.NewDeleteMealTemplate(templateRepo)
	createMealUC := mealusecases.NewCreateMeal(mealRepo, templateRepo)
	listAvailableMealsUC := mealusecases.NewListAvailableMeals(mealRepo)
	deleteMealUC := mealusecases.NewDeleteMeal(mealRepo)
	bookMealUC := mealusecases.NewBookMeal(mealRepo, bookingRepo, transactor)
	cancelBookingUC := mealusecases.NewCancelBooking(bookingRepo, mealRepo, transactor)
	listMyBookingsUC := mealusecases.NewListMyBookings(bookingRepo)
	getDailySummaryUC := mealusecases.NewGetDailySummary(bookingRepo)

	// Meal handlers
	mealHandler := mealhttp.NewMealHandler(createMealUC, listAvailableMealsUC, deleteMealUC)
	bookingHandler := mealhttp.NewBookingHandler(bookMealUC, cancelBookingUC, listMyBookingsUC, getDailySummaryUC)
	templateHandler := mealhttp.NewMealTemplateHandler(createMealTemplateUC, listMealTemplatesUC, updateMealTemplateUC, deleteMealTemplateUC)

	return &Dependencies{
		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		MealHandler:         mealHandler,
		BookingHandler:      bookingHandler,
		MealTemplateHandler: templateHandler,
		TokenPrv:            jwtProvider,
	}
}
