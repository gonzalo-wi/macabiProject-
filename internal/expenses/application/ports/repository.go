package expensesports

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

// ExpenseRepository covers transactional operations and reads needed by write use cases.
type ExpenseRepository interface {
	Save(ctx context.Context, e *expensesdomain.Expense) error
	Update(ctx context.Context, e *expensesdomain.Expense) error
	DeleteByID(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*expensesdomain.Expense, error)
	FindProjectName(ctx context.Context, projectID string) (string, error)
}

// ExpenseQueryRepository covers read-only queries and analytics.
type ExpenseQueryRepository interface {
	FindDetailByID(ctx context.Context, id string) (*expensesdomain.ExpenseDetailItem, error)
	ListAll(ctx context.Context, filter ExpenseListFilter, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	ListByProject(ctx context.Context, projectID string, onlySubmittedBy *string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	ListMine(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error)
	SummaryByProject(ctx context.Context, projectID string, onlySubmittedBy *string, from, to *time.Time) (*ProjectExpenseSummary, error)
	Analytics(ctx context.Context, from, to *time.Time, granularity string) (*ExpenseAnalyticsResult, error)

	GetBudget(ctx context.Context, projectID string) (*decimal.Decimal, error)
	SetBudget(ctx context.Context, projectID string, amount *decimal.Decimal) error
}

// ProjectBudgetStatus = presupuesto mensual del proyecto + lo aprobado en el mes actual.
type ProjectBudgetStatus struct {
	MonthlyAmount        *decimal.Decimal // nil = sin presupuesto definido
	CurrentMonthApproved decimal.Decimal
	Month                string // "YYYY-MM"
}

type ExpenseAnalyticsResult struct {
	TotalApproved decimal.Decimal
	TotalCount    int64
	PendingCount  int64
	PendingTotal  decimal.Decimal
	ApprovedCount int64
	RejectedCount int64
	RejectedTotal decimal.Decimal
	ByProject     []AnalyticsProjectTotal
	ByBucket      []AnalyticsBucketTotal
	Granularity   string
}

type ExpensePeriodMetrics struct {
	TotalCount    int64
	ApprovedCount int64
	ApprovedTotal decimal.Decimal
	PendingCount  int64
	PendingTotal  decimal.Decimal
	RejectedCount int64
	RejectedTotal decimal.Decimal
}

type AnalyticsProjectTotal struct {
	ProjectID   string
	ProjectName string
	Total       decimal.Decimal
}

type AnalyticsBucketTotal struct {
	Bucket string
	Total  decimal.Decimal
}

type ExpenseListFilter struct {
	ProjectID       string
	OnlySubmittedBy *string
	Status          string
	From            *time.Time
	To              *time.Time
	Query           string
}

type ProjectExpenseSummary struct {
	TotalApproved decimal.Decimal
	TotalCount    int64
	ApprovedCount int64
	PendingCount  int64
	PendingTotal  decimal.Decimal
	RejectedCount int64
	RejectedTotal decimal.Decimal
	ByMonth       []MonthlyTotal
}

type MonthlyTotal struct {
	MonthYYYYMM string `json:"month"`
	Total       decimal.Decimal
}
