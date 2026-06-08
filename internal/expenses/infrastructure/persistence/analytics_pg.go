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
	applyRange := func(q *gorm.DB) *gorm.DB {
		if from != nil {
			q = q.Where("expense_date >= ?", from.UTC())
		}
		if to != nil {
			q = q.Where("expense_date <= ?", to.UTC())
		}
		return q
	}
	out := &expensesports.ExpenseAnalyticsResult{Granularity: granularity}
	approved := string(expensesdomain.StatusApproved)

	type statusCount struct {
		Status string
		Cnt    int64
	}
	var scs []statusCount
	if err := applyRange(r.db.WithContext(ctx).Model(&ExpenseModel{})).
		Select("status, COUNT(*) as cnt").Group("status").Scan(&scs).Error; err != nil {
		return nil, err
	}
	for _, s := range scs {
		out.TotalCount += s.Cnt
		switch s.Status {
		case string(expensesdomain.StatusPending):
			out.PendingCount = s.Cnt
		case approved:
			out.ApprovedCount = s.Cnt
		case string(expensesdomain.StatusRejected):
			out.RejectedCount = s.Cnt
		}
	}

	type projRow struct {
		ProjectID   string          `gorm:"column:project_id"`
		ProjectName string          `gorm:"column:project_name"`
		Total       decimal.Decimal `gorm:"column:total"`
	}
	var prs []projRow
	if err := applyRange(r.db.WithContext(ctx).Table("project_expenses e")).
		Select("e.project_id as project_id, p.name as project_name, COALESCE(SUM(e.amount),0) as total").
		Joins("JOIN projects p ON p.id = e.project_id").
		Where("e.status = ?", approved).
		Group("e.project_id, p.name").
		Order("total DESC").Scan(&prs).Error; err != nil {
		return nil, err
	}
	for _, p := range prs {
		out.TotalApproved = out.TotalApproved.Add(p.Total)
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
	if err := applyRange(r.db.WithContext(ctx).Table("project_expenses e")).
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
	q := r.db.WithContext(ctx).Where("project_id = ? AND status = ?", projectID, string(expensesdomain.StatusApproved))
	if onlySubmittedBy != nil && *onlySubmittedBy != "" {
		q = q.Where("submitted_by_user_id = ?", *onlySubmittedBy)
	}
	if from != nil {
		q = q.Where("expense_date >= ?", from.UTC())
	}
	if to != nil {
		q = q.Where("expense_date <= ?", to.UTC())
	}

	var models []ExpenseModel
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}

	monthTotals := make(map[string]decimal.Decimal)
	var grand decimal.Decimal
	for _, m := range models {
		grand = grand.Add(m.Amount)
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
		TotalApproved: grand,
		ByMonth:       byMonth,
	}, nil
}

func (r *ExpenseRepositoryPG) FindProjectName(ctx context.Context, projectID string) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("projects").Select("name").Where("id = ?", projectID).Scan(&name).Error
	return name, err
}
