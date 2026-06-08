package deps

import (
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/database"
	userports "macabi-back/internal/user/application/ports"
	userusecases "macabi-back/internal/user/application/usecases"
	userhttp "macabi-back/internal/user/infrastructure/http"
	usermail "macabi-back/internal/user/infrastructure/mail"
	userpersistence "macabi-back/internal/user/infrastructure/persistence"
	usersecurity "macabi-back/internal/user/infrastructure/security"

	"gorm.io/gorm"
)

func buildUserDeps(db *gorm.DB, cfg *config.Config) (*userhttp.AuthHandler, *userhttp.UserHandler, userports.TokenProvider) {
	userRepo := userpersistence.NewUserRepositoryPG(db)
	inviteRepo := userpersistence.NewUserInvitationRepositoryPG(db)
	tokenRepo := userpersistence.NewPasswordResetTokenRepositoryPG(db)
	hasher := usersecurity.NewBcryptHasher()
	jwtProvider := usersecurity.NewJWTProvider(cfg.JWTSecret, cfg.JWTExpiration)
	transactor := database.NewGORMTransactor(db)
	invitationMailer := usermail.NewBrevoTransactionalMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom)
	passwordResetMailer := usermail.NewBrevoPasswordResetMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom)

	loginUC := userusecases.NewLogin(userRepo, hasher, jwtProvider)
	acceptInvitationUC := userusecases.NewAcceptInvitation(transactor, userRepo, inviteRepo, hasher)
	createInvitationUC := userusecases.NewCreateUserInvitation(
		userRepo,
		inviteRepo,
		invitationMailer,
		cfg.FrontendPublicURL,
		cfg.InvitationTTL,
	)
	listPendingInvitationsUC := userusecases.NewListPendingInvitations(inviteRepo, userRepo)
	resendInvitationUC := userusecases.NewResendUserInvitation(
		userRepo,
		inviteRepo,
		invitationMailer,
		cfg.FrontendPublicURL,
		cfg.InvitationTTL,
	)
	revokeInvitationUC := userusecases.NewRevokeUserInvitation(inviteRepo)
	requestPasswordUC := userusecases.NewRequestPasswordReset(
		userRepo,
		tokenRepo,
		passwordResetMailer,
		cfg.FrontendPublicURL,
		cfg.PasswordResetTTL,
	)
	resetPasswordUC := userusecases.NewResetPassword(transactor, userRepo, tokenRepo, hasher)
	getCurrentUserUC := userusecases.NewGetCurrentUser(userRepo)
	changeRoleUC := userusecases.NewChangeRole(userRepo)
	listUsersUC := userusecases.NewListUsers(userRepo)
	setUserStatusUC := userusecases.NewSetUserStatus(userRepo)
	updateUserUC := userusecases.NewUpdateUser(userRepo)
	changePasswordUC := userusecases.NewChangePassword(userRepo, hasher)

	authHandler := userhttp.NewAuthHandler(
		loginUC,
		acceptInvitationUC,
		requestPasswordUC,
		resetPasswordUC,
	)
	userHandler := userhttp.NewUserHandler(
		getCurrentUserUC,
		changeRoleUC,
		listUsersUC,
		setUserStatusUC,
		updateUserUC,
		changePasswordUC,
		createInvitationUC,
		listPendingInvitationsUC,
		resendInvitationUC,
		revokeInvitationUC,
	)

	return authHandler, userHandler, jwtProvider
}
