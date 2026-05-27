package expenseshttp

import (
	"time"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

const timeFormat = "2006-01-02T15:04:05Z07:00"

type CreateExpenseBody struct {
	ProjectID   string `json:"project_id" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	Description string `json:"description" binding:"required"`
	ExpenseDate string `json:"expense_date" binding:"required"`
	Currency    string `json:"currency"`
}

type PatchExpenseBody struct {
	Amount             *string `json:"amount"`
	Description        *string `json:"description"`
	ExpenseDate        *string `json:"expense_date"`
	Currency           *string `json:"currency"`
	ReceiptStoragePath *string `json:"receipt_storage_path"`
}

type RejectBody struct {
	RejectionReason string `json:"rejection_reason"`
}

type ReceiptUploadBody struct {
	ContentType string `json:"content_type" binding:"required"`
}

type expenseNotificationResponse struct {
	ID        string  `json:"id"`
	ExpenseID string  `json:"expense_id"`
	ProjectID string  `json:"project_id"`
	Message   string  `json:"message"`
	ReadAt    *string `json:"read_at"`
	CreatedAt string  `json:"created_at"`
}

type expenseResponse struct {
	ID                 string  `json:"id"`
	ProjectID          string  `json:"project_id"`
	ProjectName        string  `json:"project_name,omitempty"`
	SubmittedByUserID  string  `json:"submitted_by_user_id"`
	SubmitterName      string  `json:"submitter_name,omitempty"`
	Amount             string  `json:"amount"`
	Currency           string  `json:"currency"`
	Description        string  `json:"description"`
	ExpenseDate        string  `json:"expense_date"`
	Status             string  `json:"status"`
	ReceiptStoragePath *string `json:"receipt_storage_path"`
	ApprovedByUserID   *string `json:"approved_by_user_id"`
	ApprovedAt         *string `json:"approved_at"`
	RejectionReason    *string `json:"rejection_reason"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type summaryResponse struct {
	TotalApproved string                  `json:"total_approved"`
	ByMonth       []monthlyBucketResponse `json:"by_month"`
}

type monthlyBucketResponse struct {
	Month string `json:"month"`
	Total string `json:"total"`
}

func expenseToResp(e expensesdomain.Expense, submitter string) expenseResponse {
	var approved *string
	if e.ApprovedAt != nil {
		s := e.ApprovedAt.Format(timeFormat)
		approved = &s
	}
	return expenseResponse{
		ID:                 e.ID,
		ProjectID:          e.ProjectID,
		SubmittedByUserID:  e.SubmittedByUserID,
		SubmitterName:      submitter,
		Amount:             e.Amount.String(),
		Currency:           e.Currency,
		Description:        e.Description,
		ExpenseDate:        e.ExpenseDate.Format("2006-01-02"),
		Status:             string(e.Status),
		ReceiptStoragePath: e.ReceiptStoragePath,
		ApprovedByUserID:   e.ApprovedByUserID,
		ApprovedAt:         approved,
		RejectionReason:    e.RejectionReason,
		CreatedAt:          e.CreatedAt.Format(timeFormat),
		UpdatedAt:          e.UpdatedAt.Format(timeFormat),
	}
}

func listItemToExpenseResponse(row expensesdomain.ExpenseListItem) expenseResponse {
	out := expenseToResp(row.Expense, row.SubmitterName)
	out.ProjectName = row.ProjectName
	return out
}

func listResp(res pagination.Result[expensesdomain.ExpenseListItem]) pagination.Result[expenseResponse] {
	out := pagination.Result[expenseResponse]{
		Total:      res.Total,
		Page:       res.Page,
		PageSize:   res.PageSize,
		TotalPages: res.TotalPages,
	}
	out.Data = make([]expenseResponse, len(res.Data))
	for i, row := range res.Data {
		out.Data[i] = listItemToExpenseResponse(row)
	}
	return out
}

func summaryToResp(s *expensesports.ProjectExpenseSummary) summaryResponse {
	by := make([]monthlyBucketResponse, len(s.ByMonth))
	for i, b := range s.ByMonth {
		by[i] = monthlyBucketResponse{Month: b.MonthYYYYMM, Total: b.Total.String()}
	}
	return summaryResponse{TotalApproved: s.TotalApproved.String(), ByMonth: by}
}

func toExpenseNotificationListResponse(result pagination.Result[expensesdomain.ExpenseNotification]) pagination.Result[expenseNotificationResponse] {
	out := make([]expenseNotificationResponse, len(result.Data))
	for i, n := range result.Data {
		resp := expenseNotificationResponse{
			ID:        n.ID,
			ExpenseID: n.ExpenseID,
			ProjectID: n.ProjectID,
			Message:   n.Message,
			CreatedAt: n.CreatedAt.Format(timeFormat),
		}
		if n.ReadAt != nil {
			s := n.ReadAt.Format(timeFormat)
			resp.ReadAt = &s
		}
		out[i] = resp
	}
	return pagination.Result[expenseNotificationResponse]{
		Data:       out,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

func parseExpenseDate(layout, v string) (time.Time, error) {
	t, err := time.Parse(layout, v)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC().Truncate(24 * time.Hour), nil
}
