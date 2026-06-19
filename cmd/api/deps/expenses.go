package deps

import (
	"strings"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesservices "macabi-back/internal/expenses/application/services"
	expensesusecases "macabi-back/internal/expenses/application/usecases"
	expenseshttp "macabi-back/internal/expenses/infrastructure/http"
	expensesmail "macabi-back/internal/expenses/infrastructure/mail"
	expensespersistence "macabi-back/internal/expenses/infrastructure/persistence"
	projectports "macabi-back/internal/project/application/ports"
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/notifications"
	sharedstorage "macabi-back/internal/shared/storage"
	userports "macabi-back/internal/user/application/ports"

	"gorm.io/gorm"
)

func buildExpensesDeps(
	db *gorm.DB,
	cfg *config.Config,
	membership projectports.ProjectMembership,
	coordinatorReader projectports.ProjectCoordinatorReader,
	emailReader userports.UserEmailReader,
	push notifications.PushNotifier,
) *expenseshttp.Handler {
	expenseRepo := expensespersistence.NewExpenseRepositoryPG(db)
	if err := expensespersistence.RunMigrations(db); err != nil {
		panic("expenses migrations failed: " + err.Error())
	}

	var receiptSigner expensesports.ReceiptSigner
	if cfg.SupabaseURL != "" && cfg.SupabaseServiceRoleKey != "" && cfg.SupabaseExpenseReceiptBucket != "" {
		receiptSigner = &sharedstorage.SupabaseSigner{
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

	notifier := expensesservices.NewExpenseNotificationService(expenseRepo, coordinatorReader, emailReader, expenseMailer, expenseRepo, push)

	createExpenseUC := expensesusecases.NewCreateExpense(expenseRepo, membership, notifier)
	receiptUploadFileUC := expensesusecases.NewReceiptUploadFile(receiptSigner, expenseRepo, membership)

	return expenseshttp.NewHandler(
		createExpenseUC,
		expensesusecases.NewCreateExpenseWithReceipt(createExpenseUC, receiptUploadFileUC, expenseRepo),
		expensesusecases.NewGetExpense(expenseRepo, membership),
		expensesusecases.NewListProjectExpenses(expenseRepo, membership),
		expensesusecases.NewListAllExpenses(expenseRepo),
		expensesusecases.NewListMyExpenses(expenseRepo),
		expensesusecases.NewUpdateExpense(expenseRepo, membership),
		expensesusecases.NewApproveExpense(expenseRepo, membership, notifier),
		expensesusecases.NewRejectExpense(expenseRepo, membership, notifier),
		expensesusecases.NewDeleteExpense(expenseRepo, expenseRepo, membership, receiptSigner),
		expensesusecases.NewProjectExpenseSummaryUC(expenseRepo, membership),
		expensesusecases.NewReceiptUploadURL(receiptSigner, expenseRepo, membership),
		receiptUploadFileUC,
		expensesusecases.NewReceiptDownloadURL(receiptSigner, expenseRepo, membership),
		expensesusecases.NewRemoveReceipt(expenseRepo, membership, receiptSigner),
		expensesusecases.NewExpenseAnalytics(expenseRepo),
		expensesusecases.NewListExpenseNotifications(expenseRepo),
		expensesusecases.NewMarkExpenseNotificationRead(expenseRepo),
		expensesusecases.NewMarkAllExpenseNotificationsRead(expenseRepo),
		expensesusecases.NewUnreadExpenseNotificationCount(expenseRepo),
		expensesusecases.NewCreateExpenseCategory(expenseRepo),
		expensesusecases.NewListExpenseCategories(expenseRepo),
		expensesusecases.NewDeleteExpenseCategory(expenseRepo),
		expensesusecases.NewGetProjectBudget(expenseRepo, membership),
		expensesusecases.NewSetProjectBudget(expenseRepo),
	)
}
