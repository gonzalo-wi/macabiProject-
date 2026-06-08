package expensespersistence

import (
	"context"
	"strings"

	"gorm.io/gorm"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type expenseListRow struct {
	ExpenseModel
	SubmitterName string `gorm:"column:submitter_name"`
	ProjectName   string `gorm:"column:project_name"`
	CategoryName  string `gorm:"column:category_name"`
}

func (r *ExpenseRepositoryPG) ListAll(ctx context.Context, filter expensesports.ExpenseListFilter, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error) {
	apply := func(q *gorm.DB) *gorm.DB {
		if filter.ProjectID != "" {
			q = q.Where("e.project_id = ?", filter.ProjectID)
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

	base := func() *gorm.DB {
		return r.db.WithContext(ctx).
			Table("project_expenses e").
			Joins("JOIN users u ON u.id = e.submitted_by_user_id").
			Joins("JOIN projects p ON p.id = e.project_id").
			Joins("LEFT JOIN expense_categories cat ON cat.id = e.category_id")
	}

	var total int64
	if err := apply(base()).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	var rows []expenseListRow
	if err := apply(base().Select("e.*, u.name as submitter_name, p.name as project_name, COALESCE(cat.name,'') as category_name")).
		Order("e.expense_date DESC, e.created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	items := make([]expensesdomain.ExpenseListItem, len(rows))
	for i, row := range rows {
		items[i] = listRowToItem(row)
	}
	return pagination.NewResult(items, total, params), nil
}

func (r *ExpenseRepositoryPG) ListByProject(ctx context.Context, projectID string, onlySubmittedBy *string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error) {
	base := r.db.WithContext(ctx).Model(&ExpenseModel{}).Where("project_id = ?", projectID)
	if onlySubmittedBy != nil && *onlySubmittedBy != "" {
		base = base.Where("submitted_by_user_id = ?", *onlySubmittedBy)
	}
	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	q := r.db.WithContext(ctx).
		Table("project_expenses e").
		Select("e.*, u.name as submitter_name, p.name as project_name, COALESCE(cat.name,'') as category_name").
		Joins("JOIN users u ON u.id = e.submitted_by_user_id").
		Joins("JOIN projects p ON p.id = e.project_id").
		Joins("LEFT JOIN expense_categories cat ON cat.id = e.category_id").
		Where("e.project_id = ?", projectID)
	if onlySubmittedBy != nil && *onlySubmittedBy != "" {
		q = q.Where("e.submitted_by_user_id = ?", *onlySubmittedBy)
	}

	var rows []expenseListRow
	if err := q.Order("e.expense_date DESC, e.created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	items := make([]expensesdomain.ExpenseListItem, len(rows))
	for i, row := range rows {
		items[i] = listRowToItem(row)
	}
	return pagination.NewResult(items, total, params), nil
}

func (r *ExpenseRepositoryPG) ListMine(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error) {
	base := r.db.WithContext(ctx).Model(&ExpenseModel{}).Where("submitted_by_user_id = ?", userID)
	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	q := r.db.WithContext(ctx).
		Table("project_expenses e").
		Select("e.*, u.name as submitter_name, p.name as project_name, COALESCE(cat.name,'') as category_name").
		Joins("JOIN users u ON u.id = e.submitted_by_user_id").
		Joins("JOIN projects p ON p.id = e.project_id").
		Joins("LEFT JOIN expense_categories cat ON cat.id = e.category_id").
		Where("e.submitted_by_user_id = ?", userID)

	var rows []expenseListRow
	if err := q.Order("e.expense_date DESC, e.created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	items := make([]expensesdomain.ExpenseListItem, len(rows))
	for i, row := range rows {
		items[i] = listRowToItem(row)
	}
	return pagination.NewResult(items, total, params), nil
}

func listRowToItem(row expenseListRow) expensesdomain.ExpenseListItem {
	var e expensesdomain.Expense
	toDomainExpense(&row.ExpenseModel, &e)
	return expensesdomain.ExpenseListItem{
		Expense:       e,
		SubmitterName: row.SubmitterName,
		ProjectName:   row.ProjectName,
		CategoryName:  row.CategoryName,
	}
}
