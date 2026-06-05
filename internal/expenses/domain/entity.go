package expensesdomain

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type ExpenseStatus string

const (
	StatusPending  ExpenseStatus = "PENDIENTE"
	StatusApproved ExpenseStatus = "APROBADO"
	StatusRejected ExpenseStatus = "RECHAZADO"
)

const DefaultCurrency = "ARS"

type Expense struct {
	ID                 string
	ProjectID          string
	SubmittedByUserID  string
	Amount             decimal.Decimal
	Currency           string
	Description        string
	ExpenseDate        time.Time
	Status             ExpenseStatus
	CategoryID         *string
	ReceiptStoragePath *string
	ApprovedByUserID   *string
	ApprovedAt         *time.Time
	RejectionReason    *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ExpenseCategory struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type ExpenseNotification struct {
	ID        string
	UserID    string
	ExpenseID string
	ProjectID string
	Message   string
	ReadAt    *time.Time
	CreatedAt time.Time
}

type ExpenseListItem struct {
	Expense
	SubmitterName string
	ProjectName   string
	CategoryName  string
}

type ExpenseDetailItem struct {
	Expense
	SubmitterName string
	ProjectName   string
	ApproverName  string
	CategoryName  string
}

func ReceiptKeyPrefix(projectID, expenseID string) string {
	return fmt.Sprintf("projects/%s/expenses/%s", projectID, expenseID)
}

func ReceiptFileName(storagePath string) string {
	storagePath = strings.TrimSpace(storagePath)
	if storagePath == "" {
		return "comprobante"
	}
	if i := strings.LastIndex(storagePath, "/"); i >= 0 && i < len(storagePath)-1 {
		return storagePath[i+1:]
	}
	return storagePath
}

func ValidateReceiptObjectKey(projectID, expenseID, objectKey string) bool {
	prefix := ReceiptKeyPrefix(projectID, expenseID) + "/"
	if !strings.HasPrefix(objectKey, prefix) {
		return false
	}
	rest := strings.TrimPrefix(objectKey, prefix)
	if rest == "" || strings.Contains(rest, "/") || strings.Contains(rest, "..") {
		return false
	}
	return true
}

func NewExpense(projectID, submittedBy string, amount decimal.Decimal, description string, expenseDate time.Time) (*Expense, error) {
	projectID = strings.TrimSpace(projectID)
	submittedBy = strings.TrimSpace(submittedBy)
	description = strings.TrimSpace(description)
	if projectID == "" || submittedBy == "" || description == "" {
		return nil, ErrMissingRequiredField
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}
	if expenseDate.IsZero() {
		return nil, ErrMissingRequiredField
	}
	return &Expense{
		ProjectID:         projectID,
		SubmittedByUserID: submittedBy,
		Amount:            amount,
		Currency:          DefaultCurrency,
		Description:       description,
		ExpenseDate:       expenseDate.UTC().Truncate(24 * time.Hour),
		Status:            StatusPending,
	}, nil
}
