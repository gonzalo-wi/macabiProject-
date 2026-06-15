package expensespersistence

import (
	"context"

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
	base := func() *gorm.DB {
		return r.expenseListBase(ctx)
	}

	var total int64
	if err := applyExpenseListFilter(base(), filter).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	var rows []expenseListRow
	if err := applyExpenseListFilter(base().Select("e.*, u.name as submitter_name, p.name as project_name, COALESCE(cat.name,'') as category_name"), filter).
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
	filter := expensesports.ExpenseListFilter{
		ProjectID:       projectID,
		OnlySubmittedBy: onlySubmittedBy,
	}
	base := func() *gorm.DB {
		return r.expenseListBase(ctx)
	}

	var total int64
	if err := applyExpenseListFilter(base(), filter).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	var rows []expenseListRow
	if err := applyExpenseListFilter(base().Select("e.*, u.name as submitter_name, p.name as project_name, COALESCE(cat.name,'') as category_name"), filter).
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

func (r *ExpenseRepositoryPG) ListMine(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseListItem], error) {
	filter := expensesports.ExpenseListFilter{
		OnlySubmittedBy: &userID,
	}
	base := func() *gorm.DB {
		return r.expenseListBase(ctx)
	}

	var total int64
	if err := applyExpenseListFilter(base(), filter).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, err
	}

	var rows []expenseListRow
	if err := applyExpenseListFilter(base().Select("e.*, u.name as submitter_name, p.name as project_name, COALESCE(cat.name,'') as category_name"), filter).
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
