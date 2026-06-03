package expensesports

import (
	"context"
	"io"
	"time"

	"github.com/shopspring/decimal"

	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

// ProjectMembership matches the subset of reads needed for authorization (implemented by StockRepositoryPG).
type ProjectMembership interface {
	IsProjectCoordinator(ctx context.Context, projectID, userID string) (bool, error)
	IsProjectMember(ctx context.Context, projectID, userID string) (bool, error)
}

// ExpenseRepository persists expenses.
type ExpenseRepository interface {
	Save(ctx context.Context, e *expensesdomain.Expense) error
	Update(ctx context.Context, e *expensesdomain.Expense) error
	FindByID(ctx context.Context, id string) (*expensesdomain.Expense, error)
	FindDetailByID(ctx context.Context, id string) (*expensesdomain.ExpenseDetailItem, error)
	ListAll(ctx context.Context, filter ExpenseListFilter, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	DeleteByID(ctx context.Context, id string) error
	ListByProject(ctx context.Context, projectID string, onlySubmittedBy *string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	ListMine(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	SummaryByProject(ctx context.Context, projectID string, onlySubmittedBy *string, from, to *time.Time) (*ProjectExpenseSummary, error)

	SaveNotification(ctx context.Context, n *expensesdomain.ExpenseNotification) error
	ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseNotification], error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error)
	UnreadNotificationCount(ctx context.Context, userID string) (int64, error)
	DeleteNotificationsByExpenseID(ctx context.Context, expenseID string) error
	FindProjectName(ctx context.Context, projectID string) (string, error)
	Analytics(ctx context.Context, from, to *time.Time, granularity string) (*ExpenseAnalyticsResult, error)
}

// ExpenseAnalyticsResult agrega gastos para los gráficos del panel admin (date-filtered, sin paginar).
type ExpenseAnalyticsResult struct {
	TotalApproved decimal.Decimal
	TotalCount    int64
	PendingCount  int64
	ApprovedCount int64
	RejectedCount int64
	ByProject     []AnalyticsProjectTotal // aprobados por proyecto (torta)
	ByBucket      []AnalyticsBucketTotal   // aprobados por día/mes (barras)
	Granularity   string                   // "day" | "month"
}

type AnalyticsProjectTotal struct {
	ProjectID   string
	ProjectName string
	Total       decimal.Decimal
}

type AnalyticsBucketTotal struct {
	Bucket string // "YYYY-MM-DD" (day) o "YYYY-MM" (month)
	Total  decimal.Decimal
}

// ExpenseListFilter holds optional filters for the global (admin) expenses list.
type ExpenseListFilter struct {
	ProjectID string
	Status    string
	From      *time.Time
	To        *time.Time
	Query     string // matches description / submitter name / project name (ILIKE)
}

// ProjectExpenseSummary is a lightweight aggregate for charts.
type ProjectExpenseSummary struct {
	TotalApproved decimal.Decimal
	ByMonth       []MonthlyTotal // month start UTC YYYY-MM-01 → total approved in that calendar month in project TZ as DB date truncation
}

// MonthlyTotal is one bucket for charting.
type MonthlyTotal struct {
	MonthYYYYMM string `json:"month"` // YYYY-MM
	Total       decimal.Decimal
}

// ReceiptSigner stores and serves private expense receipts (Supabase Storage).
type ReceiptSigner interface {
	CreateSignedUploadURL(ctx context.Context, objectKey, contentType string) (uploadURL string, err error)
	CreateSignedDownloadURL(ctx context.Context, objectKey string, expiresSec int) (downloadURL string, err error)
	UploadObject(ctx context.Context, objectKey, contentType string, body io.Reader) error
	DeleteObject(ctx context.Context, objectKey string) error
}
