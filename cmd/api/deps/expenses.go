package deps

import (
	"strings"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesusecases "macabi-back/internal/expenses/application/usecases"
	expenseshttp "macabi-back/internal/expenses/infrastructure/http"
	expensesmail "macabi-back/internal/expenses/infrastructure/mail"
	expensespersistence "macabi-back/internal/expenses/infrastructure/persistence"
	expensesstorage "macabi-back/internal/expenses/infrastructure/storage"
	"macabi-back/internal/shared/config"
	stockpersistence "macabi-back/internal/stock/infrastructure/persistence"

	"gorm.io/gorm"
)

func buildExpensesDeps(db *gorm.DB, cfg *config.Config, stockRepo *stockpersistence.StockRepositoryPG) *expenseshttp.Handler {
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

	createExpenseUC := expensesusecases.NewCreateExpense(expenseRepo, expenseRepo, stockRepo, stockRepo, stockRepo, expenseMailer)
	receiptUploadFileUC := expensesusecases.NewReceiptUploadFile(receiptSigner, expenseRepo, stockRepo)

	return expenseshttp.NewHandler(
		createExpenseUC,
		expensesusecases.NewCreateExpenseWithReceipt(createExpenseUC, receiptUploadFileUC, expenseRepo),
		expensesusecases.NewGetExpense(expenseRepo, stockRepo),
		expensesusecases.NewListProjectExpenses(expenseRepo, stockRepo),
		expensesusecases.NewListAllExpenses(expenseRepo),
		expensesusecases.NewListMyExpenses(expenseRepo),
		expensesusecases.NewUpdateExpense(expenseRepo, stockRepo),
		expensesusecases.NewApproveExpense(expenseRepo, expenseRepo, stockRepo, stockRepo, expenseMailer),
		expensesusecases.NewRejectExpense(expenseRepo, expenseRepo, stockRepo, stockRepo, expenseMailer),
		expensesusecases.NewDeleteExpense(expenseRepo, expenseRepo, stockRepo, receiptSigner),
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
		expensesusecases.NewCreateExpenseCategory(expenseRepo),
		expensesusecases.NewListExpenseCategories(expenseRepo),
		expensesusecases.NewDeleteExpenseCategory(expenseRepo),
	)
}
