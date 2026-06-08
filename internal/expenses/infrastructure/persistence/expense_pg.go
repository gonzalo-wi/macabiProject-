package expensespersistence

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	expensesdomain "macabi-back/internal/expenses/domain"
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
	CategoryID         *string         `gorm:"type:uuid;index"`
	ReceiptStoragePath *string         `gorm:"type:text"`
	ApprovedByUserID   *string         `gorm:"type:uuid"`
	ApprovedAt         *time.Time
	RejectionReason    *string `gorm:"type:text"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (ExpenseModel) TableName() string { return "project_expenses" }

type expenseDetailRow struct {
	ExpenseModel
	SubmitterName string `gorm:"column:submitter_name"`
	ProjectName   string `gorm:"column:project_name"`
	ApproverName  string `gorm:"column:approver_name"`
	CategoryName  string `gorm:"column:category_name"`
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

func (r *ExpenseRepositoryPG) FindDetailByID(ctx context.Context, id string) (*expensesdomain.ExpenseDetailItem, error) {
	var row expenseDetailRow
	err := r.db.WithContext(ctx).
		Table("project_expenses e").
		Select("e.*, u.name as submitter_name, p.name as project_name, a.name as approver_name, COALESCE(cat.name,'') as category_name").
		Joins("JOIN users u ON u.id = e.submitted_by_user_id").
		Joins("JOIN projects p ON p.id = e.project_id").
		Joins("LEFT JOIN users a ON a.id = e.approved_by_user_id").
		Joins("LEFT JOIN expense_categories cat ON cat.id = e.category_id").
		Where("e.id = ?", id).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == "" {
		return nil, expensesdomain.ErrExpenseNotFound
	}
	var e expensesdomain.Expense
	toDomainExpense(&row.ExpenseModel, &e)
	return &expensesdomain.ExpenseDetailItem{
		Expense:       e,
		SubmitterName: row.SubmitterName,
		ProjectName:   row.ProjectName,
		ApproverName:  row.ApproverName,
		CategoryName:  row.CategoryName,
	}, nil
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
		CategoryID:         e.CategoryID,
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
	e.CategoryID = m.CategoryID
	e.ReceiptStoragePath = m.ReceiptStoragePath
	e.ApprovedByUserID = m.ApprovedByUserID
	e.ApprovedAt = m.ApprovedAt
	e.RejectionReason = m.RejectionReason
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
}
