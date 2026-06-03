package main

import (
	eventusecases "macabi-back/internal/event/application/usecases"
	eventhttp "macabi-back/internal/event/infrastructure/http"
	eventpersistence "macabi-back/internal/event/infrastructure/persistence"
	expensesports "macabi-back/internal/expenses/application/ports"
	expensesusecases "macabi-back/internal/expenses/application/usecases"
	expenseshttp "macabi-back/internal/expenses/infrastructure/http"
	expensesmail "macabi-back/internal/expenses/infrastructure/mail"
	expensespersistence "macabi-back/internal/expenses/infrastructure/persistence"
	expensesstorage "macabi-back/internal/expenses/infrastructure/storage"
	projectusecases "macabi-back/internal/project/application/usecases"
	projecthttp "macabi-back/internal/project/infrastructure/http"
	projectpersistence "macabi-back/internal/project/infrastructure/persistence"
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/database"
	stockports "macabi-back/internal/stock/application/ports"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockhttp "macabi-back/internal/stock/infrastructure/http"
	stockmail "macabi-back/internal/stock/infrastructure/mail"
	stockpersistence "macabi-back/internal/stock/infrastructure/persistence"
	stockpush "macabi-back/internal/stock/infrastructure/push"
	userports "macabi-back/internal/user/application/ports"
	userusecases "macabi-back/internal/user/application/usecases"
	userhttp "macabi-back/internal/user/infrastructure/http"
	usermail "macabi-back/internal/user/infrastructure/mail"
	userpersistence "macabi-back/internal/user/infrastructure/persistence"
	usersecurity "macabi-back/internal/user/infrastructure/security"
	"strings"

	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler     *userhttp.AuthHandler
	UserHandler     *userhttp.UserHandler
	ProjectHandler  *projecthttp.ProjectHandler
	EventHandler    *eventhttp.Handler
	StockHandler    *stockhttp.Handler
	ExpensesHandler *expenseshttp.Handler
	TokenPrv        userports.TokenProvider
}

func BuildDependencies(db *gorm.DB, cfg *config.Config) *Dependencies {
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

	projectRepo := projectpersistence.NewProjectRepositoryPG(db)
	createProjectUC := projectusecases.NewCreateProject(projectRepo)
	listProjectsUC := projectusecases.NewListProjects(projectRepo)
	getProjectUC := projectusecases.NewGetProject(projectRepo)
	updateProjectUC := projectusecases.NewUpdateProject(projectRepo)
	deleteProjectUC := projectusecases.NewDeleteProject(projectRepo)
	addMemberUC := projectusecases.NewAddProjectMember(projectRepo)
	removeMemberUC := projectusecases.NewRemoveProjectMember(projectRepo)
	listMembersUC := projectusecases.NewListProjectMembers(projectRepo)
	projectHandler := projecthttp.NewProjectHandler(
		createProjectUC,
		listProjectsUC,
		getProjectUC,
		updateProjectUC,
		deleteProjectUC,
		addMemberUC,
		removeMemberUC,
		listMembersUC,
	)

	eventRepo := eventpersistence.NewRepositoryPG(db)
	createEventUC := eventusecases.NewCreateEvent(eventRepo)
	listEventsUC := eventusecases.NewListEvents(eventRepo)
	getEventDetailUC := eventusecases.NewGetEventDetail(eventRepo)
	updateEventUC := eventusecases.NewUpdateEvent(eventRepo)
	deleteEventUC := eventusecases.NewDeleteEvent(eventRepo)
	setEventProjectsUC := eventusecases.NewSetEventProjects(eventRepo)
	createModuleUC := eventusecases.NewCreateModule(eventRepo)
	updateModuleUC := eventusecases.NewUpdateModule(eventRepo)
	deleteModuleUC := eventusecases.NewDeleteModule(eventRepo)
	setModuleProjectsUC := eventusecases.NewSetModuleProjects(eventRepo)
	createOptionGroupUC := eventusecases.NewCreateOptionGroup(eventRepo)
	updateOptionGroupUC := eventusecases.NewUpdateOptionGroup(eventRepo)
	deleteOptionGroupUC := eventusecases.NewDeleteOptionGroup(eventRepo)
	createOptionUC := eventusecases.NewCreateOption(eventRepo)
	updateOptionUC := eventusecases.NewUpdateOption(eventRepo)
	deleteOptionUC := eventusecases.NewDeleteOption(eventRepo)
	submitResponseUC := eventusecases.NewSubmitResponse(eventRepo)
	getMyResponseUC := eventusecases.NewGetMyResponse(eventRepo)
	listEventResponsesUC := eventusecases.NewListEventResponses(eventRepo)
	getModuleResponseSummaryUC := eventusecases.NewGetModuleResponseSummary(eventRepo)
	eventHandler := eventhttp.NewHandler(
		createEventUC,
		listEventsUC,
		getEventDetailUC,
		updateEventUC,
		deleteEventUC,
		setEventProjectsUC,
		createModuleUC,
		updateModuleUC,
		deleteModuleUC,
		setModuleProjectsUC,
		createOptionGroupUC,
		updateOptionGroupUC,
		deleteOptionGroupUC,
		createOptionUC,
		updateOptionUC,
		deleteOptionUC,
		submitResponseUC,
		getMyResponseUC,
		listEventResponsesUC,
		getModuleResponseSummaryUC,
	)

	stockRepo := stockpersistence.NewStockRepositoryPG(db)
	if err := stockpersistence.RunMigrations(db); err != nil {
		panic("stock migrations failed: " + err.Error())
	}
	stockMailer := stockmail.NewBrevoStockMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom, cfg.FrontendPublicURL)
	pushSubRepo := stockpersistence.NewPushSubscriptionRepositoryPG(db)
	var pushNotifier stockports.UserPushNotifier
	if cfg.VAPIDPublicKey != "" && cfg.VAPIDPrivateKey != "" && cfg.VAPIDSubject != "" {
		pushNotifier = stockpush.NewWebPushNotifier(pushSubRepo, cfg.VAPIDPublicKey, cfg.VAPIDPrivateKey, cfg.VAPIDSubject)
	} else {
		pushNotifier = stockpush.NoOpPushNotifier{}
	}
	stockHandler := stockhttp.NewHandler(
		stockusecases.NewCreateResource(stockRepo),
		stockusecases.NewListResources(stockRepo),
		stockusecases.NewGetResource(stockRepo),
		stockusecases.NewUpdateResource(stockRepo),
		stockusecases.NewDeleteResource(stockRepo),
		stockusecases.NewCreateRequest(stockRepo, stockRepo, stockRepo, stockMailer, pushNotifier),
		stockusecases.NewApproveRequest(stockRepo, stockRepo, stockRepo, stockMailer, pushNotifier),
		stockusecases.NewRejectRequest(stockRepo, stockRepo, stockRepo, stockMailer, pushNotifier),
		stockusecases.NewCancelRequest(stockRepo),
		stockusecases.NewDeliverRequest(stockRepo, stockRepo),
		stockusecases.NewReturnRequest(stockRepo, stockRepo),
		stockusecases.NewListRequests(stockRepo, stockRepo),
		stockusecases.NewListMyRequests(stockRepo),
		stockusecases.NewGetRequestDetail(stockRepo, stockRepo),
		stockusecases.NewListNotifications(stockRepo),
		stockusecases.NewMarkNotificationRead(stockRepo),
		stockusecases.NewMarkAllNotificationsRead(stockRepo),
		stockusecases.NewUnreadNotificationCount(stockRepo),
		stockusecases.NewRegisterPushSubscription(pushSubRepo),
		stockusecases.NewUnregisterPushSubscription(pushSubRepo),
		cfg.VAPIDPublicKey,
	)

	expenseRepo := expensespersistence.NewExpenseRepositoryPG(db)
	if err := expensespersistence.RunMigrations(db); err != nil {
		panic("expenses migrations failed: " + err.Error())
	}

	var receiptSigner expensesports.ReceiptSigner
	if cfg.SupabaseURL != "" && cfg.SupabaseServiceRoleKey != "" && cfg.SupabaseExpenseReceiptBucket != "" {
		receiptSigner = &expensesstorage.SupabaseSigner{
			BaseURL: strings.TrimSuffix(cfg.SupabaseURL, "/"),
			APIKey:  cfg.SupabaseServiceRoleKey,
			Bucket:  cfg.SupabaseExpenseReceiptBucket,
		}
	}

	var expenseMailer expensesports.ExpenseMailer
	if cfg.BrevoAPIKey != "" && cfg.BrevoEmailFrom != "" {
		expenseMailer = expensesmail.NewBrevoExpenseMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom, cfg.FrontendPublicURL)
	} else {
		expenseMailer = expensesmail.NewNoOpExpenseMailer()
	}

	createExpenseUC := expensesusecases.NewCreateExpense(expenseRepo, stockRepo, stockRepo, stockRepo, expenseMailer)
	receiptUploadFileUC := expensesusecases.NewReceiptUploadFile(receiptSigner, expenseRepo, stockRepo)

	expensesHandler := expenseshttp.NewHandler(
		createExpenseUC,
		expensesusecases.NewCreateExpenseWithReceipt(createExpenseUC, receiptUploadFileUC, expenseRepo),
		expensesusecases.NewGetExpense(expenseRepo, stockRepo),
		expensesusecases.NewListProjectExpenses(expenseRepo, stockRepo),
		expensesusecases.NewListAllExpenses(expenseRepo),
		expensesusecases.NewListMyExpenses(expenseRepo),
		expensesusecases.NewUpdateExpense(expenseRepo, stockRepo),
		expensesusecases.NewApproveExpense(expenseRepo, stockRepo, stockRepo, expenseMailer),
		expensesusecases.NewRejectExpense(expenseRepo, stockRepo, stockRepo, expenseMailer),
		expensesusecases.NewDeleteExpense(expenseRepo, stockRepo, receiptSigner),
		expensesusecases.NewProjectExpenseSummaryUC(expenseRepo, stockRepo),
		expensesusecases.NewReceiptUploadURL(receiptSigner, expenseRepo, stockRepo),
		receiptUploadFileUC,
		expensesusecases.NewReceiptDownloadURL(receiptSigner, expenseRepo, stockRepo),
		expensesusecases.NewRemoveReceipt(expenseRepo, stockRepo, receiptSigner),
		expensesusecases.NewExpenseAnalytics(expenseRepo),
		expensesusecases.NewListExpenseNotifications(expenseRepo),
		expensesusecases.NewMarkExpenseNotificationRead(expenseRepo),
		expensesusecases.NewMarkAllExpenseNotificationsRead(expenseRepo),
		expensesusecases.NewUnreadExpenseNotificationCount(expenseRepo),
	)

	return &Dependencies{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		ProjectHandler:  projectHandler,
		EventHandler:    eventHandler,
		StockHandler:    stockHandler,
		ExpensesHandler: expensesHandler,
		TokenPrv:        jwtProvider,
	}
}
