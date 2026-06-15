package expensespersistence

import (
	"context"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
)

func (r *ExpenseRepositoryPG) Analytics(ctx context.Context, from, to *time.Time, granularity string) (*expensesports.ExpenseAnalyticsResult, error) {
	filter := expensesports.ExpenseListFilter{From: from, To: to}
	metrics, err := r.PeriodMetrics(ctx, filter)
	if err != nil {
		return nil, err
	}

	out := &expensesports.ExpenseAnalyticsResult{
		Granularity:   granularity,
		TotalApproved: metrics.ApprovedTotal,
		TotalCount:    metrics.TotalCount,
		PendingCount:  metrics.PendingCount,
		PendingTotal:  metrics.PendingTotal,
		ApprovedCount: metrics.ApprovedCount,
		RejectedCount: metrics.RejectedCount,
		RejectedTotal: metrics.RejectedTotal,
	}

	approved := string(expensesdomain.StatusApproved)
	applyRange := func(q *gorm.DB) *gorm.DB {
		return applyExpenseListFilter(q, filter)
	}

	type projRow struct {
		ProjectID   string          `gorm:"column:project_id"`
		ProjectName string          `gorm:"column:project_name"`
		Total       decimal.Decimal `gorm:"column:total"`
	}
	var prs []projRow
	if err := applyRange(r.expenseListBase(ctx)).
		Select("e.project_id as project_id, p.name as project_name, COALESCE(SUM(e.amount),0) as total").
		Where("e.status = ?", approved).
		Group("e.project_id, p.name").
		Order("total DESC").Scan(&prs).Error; err != nil {
		return nil, err
	}
	for _, p := range prs {
		out.ByProject = append(out.ByProject, expensesports.AnalyticsProjectTotal{
			ProjectID: p.ProjectID, ProjectName: p.ProjectName, Total: p.Total,
		})
	}

	truncUnit := "month"
	bucketLayout := "2006-01"
	if granularity == "day" {
		truncUnit = "day"
		bucketLayout = "2006-01-02"
	}
	type bucketRow struct {
		Bucket time.Time       `gorm:"column:bucket"`
		Total  decimal.Decimal `gorm:"column:total"`
	}
	var brs []bucketRow
	if err := applyRange(r.expenseListBase(ctx)).
		Select("date_trunc('"+truncUnit+"', e.expense_date) as bucket, COALESCE(SUM(e.amount),0) as total").
		Where("e.status = ?", approved).
		Group("bucket").Order("bucket ASC").Scan(&brs).Error; err != nil {
		return nil, err
	}
	for _, b := range brs {
		out.ByBucket = append(out.ByBucket, expensesports.AnalyticsBucketTotal{
			Bucket: b.Bucket.UTC().Format(bucketLayout), Total: b.Total,
		})
	}

	return out, nil
}

func (r *ExpenseRepositoryPG) SummaryByProject(ctx context.Context, projectID string, onlySubmittedBy *string, from, to *time.Time) (*expensesports.ProjectExpenseSummary, error) {
	filter := expensesports.ExpenseListFilter{
		ProjectID:       projectID,
		OnlySubmittedBy: onlySubmittedBy,
		From:            from,
		To:              to,
	}
	metrics, err := r.PeriodMetrics(ctx, filter)
	if err != nil {
		return nil, err
	}

	approved := string(expensesdomain.StatusApproved)
	var rows []ExpenseModel
	if err := applyExpenseListFilter(r.expenseListBase(ctx), filter).
		Select("e.*").
		Where("e.status = ?", approved).
		Order("e.expense_date ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	monthTotals := make(map[string]decimal.Decimal)
	for _, m := range rows {
		key := time.Date(m.ExpenseDate.UTC().Year(), m.ExpenseDate.UTC().Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
		prev := monthTotals[key]
		monthTotals[key] = prev.Add(m.Amount)
	}

	months := make([]string, 0, len(monthTotals))
	for k := range monthTotals {
		months = append(months, k)
	}
	sort.Strings(months)

	byMonth := make([]expensesports.MonthlyTotal, 0, len(months))
	for _, ym := range months {
		byMonth = append(byMonth, expensesports.MonthlyTotal{MonthYYYYMM: ym, Total: monthTotals[ym]})
	}

	return &expensesports.ProjectExpenseSummary{
		TotalApproved: metrics.ApprovedTotal,
		TotalCount:    metrics.TotalCount,
		ApprovedCount: metrics.ApprovedCount,
		PendingCount:  metrics.PendingCount,
		PendingTotal:  metrics.PendingTotal,
		RejectedCount: metrics.RejectedCount,
		RejectedTotal: metrics.RejectedTotal,
		ByMonth:       byMonth,
	}, nil
}

func (r *ExpenseRepositoryPG) FindProjectName(ctx context.Context, projectID string) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("projects").Select("name").Where("id = ?", projectID).Scan(&name).Error
	return name, err
}
