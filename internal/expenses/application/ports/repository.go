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
	DeleteByID(ctx context.Context, id string) error
	ListByProject(ctx context.Context, projectID string, onlySubmittedBy *string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	ListMine(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	SummaryByProject(ctx context.Context, projectID string, onlySubmittedBy *string, from, to *time.Time) (*ProjectExpenseSummary, error)

	SaveNotification(ctx context.Context, n *expensesdomain.ExpenseNotification) error
	ListNotificationsByUser(ctx context.Context, userID string) ([]expensesdomain.ExpenseNotification, error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	UnreadNotificationCount(ctx context.Context, userID string) (int64, error)
	DeleteNotificationsByExpenseID(ctx context.Context, expenseID string) error
	FindProjectName(ctx context.Context, projectID string) (string, error)
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
