package main

import (
	"macabi-back/internal/shared/config"
	userports "macabi-back/internal/user/application/ports"
	userusecases "macabi-back/internal/user/application/usecases"
	userhttp "macabi-back/internal/user/infrastructure/http"
	userpersistence "macabi-back/internal/user/infrastructure/persistence"
	usersecurity "macabi-back/internal/user/infrastructure/security"

	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler *userhttp.AuthHandler
	UserHandler *userhttp.UserHandler
	TokenPrv    userports.TokenProvider
}

func BuildDependencies(db *gorm.DB, cfg *config.Config) *Dependencies {
	// Infrastructure
	userRepo := userpersistence.NewUserRepositoryPG(db)
	hasher := usersecurity.NewBcryptHasher()
	jwtProvider := usersecurity.NewJWTProvider(cfg.JWTSecret, cfg.JWTExpiration)

	// Use cases
	registerUC := userusecases.NewRegisterUser(userRepo, hasher)
	loginUC := userusecases.NewLogin(userRepo, hasher, jwtProvider)
	getCurrentUserUC := userusecases.NewGetCurrentUser(userRepo)
	changeRoleUC := userusecases.NewChangeRole(userRepo)
	listUsersUC := userusecases.NewListUsers(userRepo)

	// Handlers
	authHandler := userhttp.NewAuthHandler(registerUC, loginUC)
	userHandler := userhttp.NewUserHandler(getCurrentUserUC, changeRoleUC, listUsersUC)

	return &Dependencies{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		TokenPrv:    jwtProvider,
	}
}
