package expenseshttp

import (
	"mime"
	"strings"

	"github.com/gin-gonic/gin"

	expensesusecases "macabi-back/internal/expenses/application/usecases"
)

type Handler struct {
	create            *expensesusecases.CreateExpense
	createWithReceipt *expensesusecases.CreateExpenseWithReceipt
	getExp            *expensesusecases.GetExpense
	listProj          *expensesusecases.ListProjectExpenses
	listAll           *expensesusecases.ListAllExpenses
	listMine          *expensesusecases.ListMyExpenses
	upd               *expensesusecases.UpdateExpense
	approve           *expensesusecases.ApproveExpense
	reject            *expensesusecases.RejectExpense
	deleteUC          *expensesusecases.DeleteExpense
	summaryUC         *expensesusecases.ProjectExpenseSummaryUC
	recUpload         *expensesusecases.ReceiptUploadURL
	recUploadFile     *expensesusecases.ReceiptUploadFile
	recView           *expensesusecases.ReceiptDownloadURL
	removeReceipt     *expensesusecases.RemoveReceipt
	analytics         *expensesusecases.ExpenseAnalytics
	listNotifs        *expensesusecases.ListExpenseNotifications
	markNotifRead     *expensesusecases.MarkExpenseNotificationRead
	markAllNotifsRead *expensesusecases.MarkAllExpenseNotificationsRead
	unreadNotifs      *expensesusecases.UnreadExpenseNotificationCount
	createCategory    *expensesusecases.CreateExpenseCategory
	listCategories    *expensesusecases.ListExpenseCategories
	deleteCategory    *expensesusecases.DeleteExpenseCategory
}

func NewHandler(
	create *expensesusecases.CreateExpense,
	createWithReceipt *expensesusecases.CreateExpenseWithReceipt,
	getExp *expensesusecases.GetExpense,
	listProj *expensesusecases.ListProjectExpenses,
	listAll *expensesusecases.ListAllExpenses,
	listMine *expensesusecases.ListMyExpenses,
	upd *expensesusecases.UpdateExpense,
	approve *expensesusecases.ApproveExpense,
	reject *expensesusecases.RejectExpense,
	deleteUC *expensesusecases.DeleteExpense,
	summaryUC *expensesusecases.ProjectExpenseSummaryUC,
	recUpload *expensesusecases.ReceiptUploadURL,
	recUploadFile *expensesusecases.ReceiptUploadFile,
	recView *expensesusecases.ReceiptDownloadURL,
	removeReceipt *expensesusecases.RemoveReceipt,
	analytics *expensesusecases.ExpenseAnalytics,
	listNotifs *expensesusecases.ListExpenseNotifications,
	markNotifRead *expensesusecases.MarkExpenseNotificationRead,
	markAllNotifsRead *expensesusecases.MarkAllExpenseNotificationsRead,
	unreadNotifs *expensesusecases.UnreadExpenseNotificationCount,
	createCategory *expensesusecases.CreateExpenseCategory,
	listCategories *expensesusecases.ListExpenseCategories,
	deleteCategory *expensesusecases.DeleteExpenseCategory,
) *Handler {
	return &Handler{
		create:            create,
		createWithReceipt: createWithReceipt,
		getExp:            getExp,
		listProj:          listProj,
		listAll:           listAll,
		listMine:          listMine,
		upd:               upd,
		approve:           approve,
		reject:            reject,
		deleteUC:          deleteUC,
		summaryUC:         summaryUC,
		recUpload:         recUpload,
		recUploadFile:     recUploadFile,
		recView:           recView,
		removeReceipt:     removeReceipt,
		analytics:         analytics,
		listNotifs:        listNotifs,
		markNotifRead:     markNotifRead,
		markAllNotifsRead: markAllNotifsRead,
		unreadNotifs:      unreadNotifs,
		createCategory:    createCategory,
		listCategories:    listCategories,
		deleteCategory:    deleteCategory,
	}
}

func isMultipartRequest(c *gin.Context) bool {
	ct := strings.TrimSpace(c.GetHeader("Content-Type"))
	if ct == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return strings.Contains(strings.ToLower(ct), "multipart/form-data")
	}
	return mediaType == "multipart/form-data"
}
