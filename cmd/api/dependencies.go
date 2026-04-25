package main

import (
	mealusecases "macabi-back/internal/meal/application/usecases"
	mealhttp "macabi-back/internal/meal/infrastructure/http"
	mealpersistence "macabi-back/internal/meal/infrastructure/persistence"
	"macabi-back/internal/shared/config"
	userports "macabi-back/internal/user/application/ports"
	userusecases "macabi-back/internal/user/application/usecases"
	userhttp "macabi-back/internal/user/infrastructure/http"
	userpersistence "macabi-back/internal/user/infrastructure/persistence"
	usersecurity "macabi-back/internal/user/infrastructure/security"

	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler    *userhttp.AuthHandler
	UserHandler    *userhttp.UserHandler
	MealHandler    *mealhttp.MealHandler
	BookingHandler *mealhttp.BookingHandler
	TokenPrv       userports.TokenProvider
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

	// User handlers
	authHandler := userhttp.NewAuthHandler(registerUC, loginUC)
	userHandler := userhttp.NewUserHandler(getCurrentUserUC, changeRoleUC, listUsersUC)

	// Meal infrastructure
	mealRepo := mealpersistence.NewMealRepositoryPG(db)
	bookingRepo := mealpersistence.NewBookingRepositoryPG(db)

	// Meal use cases
	createMealUC := mealusecases.NewCreateMeal(mealRepo)
	listAvailableMealsUC := mealusecases.NewListAvailableMeals(mealRepo)
	bookMealUC := mealusecases.NewBookMeal(mealRepo, bookingRepo)
	cancelBookingUC := mealusecases.NewCancelBooking(bookingRepo, mealRepo)
	listMyBookingsUC := mealusecases.NewListMyBookings(bookingRepo)
	getDailySummaryUC := mealusecases.NewGetDailySummary(bookingRepo)

	// Meal handlers
	mealHandler := mealhttp.NewMealHandler(createMealUC, listAvailableMealsUC)
	bookingHandler := mealhttp.NewBookingHandler(bookMealUC, cancelBookingUC, listMyBookingsUC, getDailySummaryUC)

	return &Dependencies{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		MealHandler:    mealHandler,
		BookingHandler: bookingHandler,
		TokenPrv:       jwtProvider,
	}
}
