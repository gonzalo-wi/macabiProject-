package expensespersistence

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type ExpenseModel struct {
	ID                 string          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID          string          `gorm:"type:uuid;not null;index:idx_project_expense"`
	SubmittedByUserID  string          `gorm:"type:uuid;not null;index"`
	Amount             decimal.Decimal `gorm:"type:numeric(16,4);not null"`
	Currency           string          `gorm:"type:varchar(8);not null;default:'ARS'"`
	Description        string          `gorm:"type:text;not null"`
	ExpenseDate        time.Time       `gorm:"type:date;not null;index"`
	Status             string          `gorm:"type:varchar(16);not null;default:'PENDIENTE';index"`
	ReceiptStoragePath *string         `gorm:"type:text"`
	ApprovedByUserID   *string         `gorm:"type:uuid"`
	ApprovedAt         *time.Time
	RejectionReason    *string `gorm:"type:text"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (ExpenseModel) TableName() string { return "project_expenses" }

type ExpenseNotificationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null;index"`
	ExpenseID string `gorm:"type:uuid;not null"`
	ProjectID string `gorm:"type:uuid;not null"`
	Message   string `gorm:"not null"`
	ReadAt    *time.Time
	CreatedAt time.Time
}

func (ExpenseNotificationModel) TableName() string { return "expense_notifications" }

type expenseListRow struct {
	ExpenseModel
	SubmitterName string `gorm:"column:submitter_name"`
	ProjectName   string `gorm:"column:project_name"`
}

type ExpenseRepositoryPG struct {
	db *gorm.DB
}

var _ expensesports.ExpenseRepository = (*ExpenseRepositoryPG)(nil)

func NewExpenseRepositoryPG(db *gorm.DB) *ExpenseRepositoryPG {
	return &ExpenseRepositoryPG{db: db}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&ExpenseModel{}, &ExpenseNotificationModel{})
}

func (r *ExpenseRepositoryPG) Save(ctx context.Context, e *expensesdomain.Expense) error {
	m := fromDomainExpense(e)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	toDomainExpense(&m, e)
	return nil
}

func (r *ExpenseRepositoryPG) Update(ctx context.Context, e *expensesdomain.Expense) error {
	m := fromDomainExpense(e)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *ExpenseRepositoryPG) DeleteByID(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&ExpenseModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return expensesdomain.ErrExpenseNotFound
	}
	return nil
}

func (r *ExpenseRepositoryPG) FindByID(ctx context.Context, id string) (*expensesdomain.Expense, error) {
	var m ExpenseModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, expensesdomain.ErrExpenseNotFound
		}
		return nil, err
	}
	var out expensesdomain.Expense
	toDomainExpense(&m, &out)
	return &out, nil
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
		Select("e.*, u.name as submitter_name, p.name as project_name").
		Joins("JOIN users u ON u.id = e.submitted_by_user_id").
		Joins("JOIN projects p ON p.id = e.project_id").
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
		Select("e.*, u.name as submitter_name, p.name as project_name").
		Joins("JOIN users u ON u.id = e.submitted_by_user_id").
		Joins("JOIN projects p ON p.id = e.project_id").
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

func listRowToItem(row expenseListRow) expensesdomain.ExpenseListItem {
	var e expensesdomain.Expense
	toDomainExpense(&row.ExpenseModel, &e)
	return expensesdomain.ExpenseListItem{Expense: e, SubmitterName: row.SubmitterName, ProjectName: row.ProjectName}
}

func fromDomainExpense(e *expensesdomain.Expense) ExpenseModel {
	m := ExpenseModel{
		ID:                 e.ID,
		ProjectID:          e.ProjectID,
		SubmittedByUserID:  e.SubmittedByUserID,
		Amount:             e.Amount,
		Currency:           e.Currency,
		Description:        e.Description,
		ExpenseDate:        e.ExpenseDate,
		Status:             string(e.Status),
		ReceiptStoragePath: e.ReceiptStoragePath,
		ApprovedByUserID:   e.ApprovedByUserID,
		ApprovedAt:         e.ApprovedAt,
		RejectionReason:    e.RejectionReason,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
	if m.Currency == "" {
		m.Currency = expensesdomain.DefaultCurrency
	}
	return m
}

func toDomainExpense(m *ExpenseModel, e *expensesdomain.Expense) {
	e.ID = m.ID
	e.ProjectID = m.ProjectID
	e.SubmittedByUserID = m.SubmittedByUserID
	e.Amount = m.Amount
	e.Currency = m.Currency
	e.Description = m.Description
	e.ExpenseDate = m.ExpenseDate.UTC().Truncate(24 * time.Hour)
	e.Status = expensesdomain.ExpenseStatus(m.Status)
	e.ReceiptStoragePath = m.ReceiptStoragePath
	e.ApprovedByUserID = m.ApprovedByUserID
	e.ApprovedAt = m.ApprovedAt
	e.RejectionReason = m.RejectionReason
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
}

func (r *ExpenseRepositoryPG) SaveNotification(ctx context.Context, n *expensesdomain.ExpenseNotification) error {
	m := ExpenseNotificationModel{
		UserID:    n.UserID,
		ExpenseID: n.ExpenseID,
		ProjectID: n.ProjectID,
		Message:   n.Message,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	return nil
}

func (r *ExpenseRepositoryPG) ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[expensesdomain.ExpenseNotification], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ExpenseNotificationModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseNotification]{}, err
	}
	var models []ExpenseNotificationModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error; err != nil {
		return pagination.Result[expensesdomain.ExpenseNotification]{}, err
	}
	out := make([]expensesdomain.ExpenseNotification, len(models))
	for i, m := range models {
		out[i] = toDomainExpenseNotification(m)
	}
	return pagination.NewResult(out, total, params), nil
}

func (r *ExpenseRepositoryPG) MarkNotificationRead(ctx context.Context, id, userID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&ExpenseNotificationModel{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return expensesdomain.ErrNotificationNotFound
	}
	return nil
}

func (r *ExpenseRepositoryPG) UnreadNotificationCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ExpenseNotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *ExpenseRepositoryPG) DeleteNotificationsByExpenseID(ctx context.Context, expenseID string) error {
	return r.db.WithContext(ctx).Where("expense_id = ?", expenseID).Delete(&ExpenseNotificationModel{}).Error
}

func (r *ExpenseRepositoryPG) FindProjectName(ctx context.Context, projectID string) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("projects").Select("name").Where("id = ?", projectID).Scan(&name).Error
	return name, err
}

func toDomainExpenseNotification(m ExpenseNotificationModel) expensesdomain.ExpenseNotification {
	return expensesdomain.ExpenseNotification{
		ID:        m.ID,
		UserID:    m.UserID,
		ExpenseID: m.ExpenseID,
		ProjectID: m.ProjectID,
		Message:   m.Message,
		ReadAt:    m.ReadAt,
		CreatedAt: m.CreatedAt,
	}
}
