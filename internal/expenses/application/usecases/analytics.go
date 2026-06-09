package expensesusecases

import (
	"context"
	"time"

	expensesports "macabi-back/internal/expenses/application/ports"
)

type ExpenseAnalytics struct {
	repo expensesports.ExpenseQueryRepository
}

func NewExpenseAnalytics(repo expensesports.ExpenseQueryRepository) *ExpenseAnalytics {
	return &ExpenseAnalytics{repo: repo}
}

type ExpenseAnalyticsInput struct {
	From *time.Time
	To   *time.Time
}

func (uc *ExpenseAnalytics) Execute(ctx context.Context, in ExpenseAnalyticsInput) (*expensesports.ExpenseAnalyticsResult, error) {
	granularity := "month"
	if in.From != nil && in.To != nil && in.To.Sub(*in.From) <= 62*24*time.Hour {
		granularity = "day"
	}
	return uc.repo.Analytics(ctx, in.From, in.To, granularity)
}
