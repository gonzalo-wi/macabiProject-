package expensesdto

import (
	"time"

	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

const TimeFormat = "2006-01-02T15:04:05Z07:00"

type CreateExpenseBody struct {
	ProjectID   string  `json:"project_id" binding:"required"`
	Amount      string  `json:"amount" binding:"required"`
	Description string  `json:"description" binding:"required"`
	ExpenseDate string  `json:"expense_date" binding:"required"`
	Currency    string  `json:"currency"`
	CategoryID  *string `json:"category_id"`
}

type PatchExpenseBody struct {
	Amount             *string `json:"amount"`
	Description        *string `json:"description"`
	ExpenseDate        *string `json:"expense_date"`
	Currency           *string `json:"currency"`
	ReceiptStoragePath *string `json:"receipt_storage_path"`
	CategoryID         *string `json:"category_id"`
}

type RejectBody struct {
	RejectionReason string `json:"rejection_reason"`
}

type ReceiptUploadBody struct {
	ContentType string `json:"content_type" binding:"required"`
}

type ExpenseResponse struct {
	ID                 string  `json:"id"`
	ProjectID          string  `json:"project_id"`
	ProjectName        string  `json:"project_name,omitempty"`
	SubmittedByUserID  string  `json:"submitted_by_user_id"`
	SubmitterName      string  `json:"submitter_name,omitempty"`
	ApprovedByName     string  `json:"approved_by_name,omitempty"`
	Amount             string  `json:"amount"`
	Currency           string  `json:"currency"`
	Description        string  `json:"description"`
	ExpenseDate        string  `json:"expense_date"`
	Status             string  `json:"status"`
	CategoryID         *string `json:"category_id"`
	CategoryName       string  `json:"category_name,omitempty"`
	ReceiptStoragePath *string `json:"receipt_storage_path"`
	ApprovedByUserID   *string `json:"approved_by_user_id"`
	ApprovedAt         *string `json:"approved_at"`
	RejectionReason    *string `json:"rejection_reason"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

func ExpenseToResp(e expensesdomain.Expense, submitter string) ExpenseResponse {
	var approved *string
	if e.ApprovedAt != nil {
		s := e.ApprovedAt.Format(TimeFormat)
		approved = &s
	}
	return ExpenseResponse{
		ID:                 e.ID,
		ProjectID:          e.ProjectID,
		SubmittedByUserID:  e.SubmittedByUserID,
		SubmitterName:      submitter,
		Amount:             e.Amount.String(),
		Currency:           e.Currency,
		Description:        e.Description,
		ExpenseDate:        e.ExpenseDate.Format("2006-01-02"),
		Status:             string(e.Status),
		CategoryID:         e.CategoryID,
		ReceiptStoragePath: e.ReceiptStoragePath,
		ApprovedByUserID:   e.ApprovedByUserID,
		ApprovedAt:         approved,
		RejectionReason:    e.RejectionReason,
		CreatedAt:          e.CreatedAt.Format(TimeFormat),
		UpdatedAt:          e.UpdatedAt.Format(TimeFormat),
	}
}

func ListItemToExpenseResponse(row expensesdomain.ExpenseListItem) ExpenseResponse {
	out := ExpenseToResp(row.Expense, row.SubmitterName)
	out.ProjectName = row.ProjectName
	out.CategoryName = row.CategoryName
	return out
}

func DetailItemToExpenseResponse(d expensesdomain.ExpenseDetailItem) ExpenseResponse {
	out := ExpenseToResp(d.Expense, d.SubmitterName)
	out.ProjectName = d.ProjectName
	out.ApprovedByName = d.ApproverName
	out.CategoryName = d.CategoryName
	return out
}

func ListResp(res pagination.Result[expensesdomain.ExpenseListItem]) pagination.Result[ExpenseResponse] {
	out := pagination.Result[ExpenseResponse]{
		Total:      res.Total,
		Page:       res.Page,
		PageSize:   res.PageSize,
		TotalPages: res.TotalPages,
	}
	out.Data = make([]ExpenseResponse, len(res.Data))
	for i, row := range res.Data {
		out.Data[i] = ListItemToExpenseResponse(row)
	}
	return out
}

func ParseExpenseDate(layout, v string) (time.Time, error) {
	t, err := time.Parse(layout, v)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC().Truncate(24 * time.Hour), nil
}
