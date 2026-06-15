package expensespersistence

import (
	"context"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
)

func (r *ExpenseRepositoryPG) expenseListBase(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Table("project_expenses e").
		Joins("JOIN users u ON u.id = e.submitted_by_user_id").
		Joins("JOIN projects p ON p.id = e.project_id").
		Joins("LEFT JOIN expense_categories cat ON cat.id = e.category_id")
}

func applyExpenseListFilter(q *gorm.DB, filter expensesports.ExpenseListFilter) *gorm.DB {
	if filter.ProjectID != "" {
		q = q.Where("e.project_id = ?", filter.ProjectID)
	}
	if filter.OnlySubmittedBy != nil && *filter.OnlySubmittedBy != "" {
		q = q.Where("e.submitted_by_user_id = ?", *filter.OnlySubmittedBy)
	}
	if filter.Status != "" {
		q = q.Where("e.status = ?", filter.Status)
	}
	if filter.From != nil {
		q = q.Where("e.expense_date >= ?", filter.From.UTC())
	}
	if filter.To != nil {
		q = q.Where("e.expense_date <= ?", filter.To.UTC())
	}
	if s := strings.TrimSpace(filter.Query); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("(LOWER(e.description) LIKE ? OR LOWER(u.name) LIKE ? OR LOWER(p.name) LIKE ?)", like, like, like)
	}
	return q
}

func (r *ExpenseRepositoryPG) PeriodMetrics(ctx context.Context, filter expensesports.ExpenseListFilter) (*expensesports.ExpensePeriodMetrics, error) {
	type statusRow struct {
		Status string
		Cnt    int64
		Total  decimal.Decimal
	}
	var rows []statusRow
	if err := applyExpenseListFilter(r.expenseListBase(ctx), filter).
		Select("e.status, COUNT(*) as cnt, COALESCE(SUM(e.amount), 0) as total").
		Group("e.status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := &expensesports.ExpensePeriodMetrics{}
	approved := string(expensesdomain.StatusApproved)
	for _, row := range rows {
		out.TotalCount += row.Cnt
		switch row.Status {
		case string(expensesdomain.StatusPending):
			out.PendingCount = row.Cnt
			out.PendingTotal = row.Total
		case approved:
			out.ApprovedCount = row.Cnt
			out.ApprovedTotal = row.Total
		case string(expensesdomain.StatusRejected):
			out.RejectedCount = row.Cnt
			out.RejectedTotal = row.Total
		}
	}
	return out, nil
}
