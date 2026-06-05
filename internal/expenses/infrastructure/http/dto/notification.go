package expensesdto

import (
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
)

type ExpenseNotificationResponse struct {
	ID        string  `json:"id"`
	ExpenseID string  `json:"expense_id"`
	ProjectID string  `json:"project_id"`
	Message   string  `json:"message"`
	ReadAt    *string `json:"read_at"`
	CreatedAt string  `json:"created_at"`
}

func ToExpenseNotificationListResponse(result pagination.Result[expensesdomain.ExpenseNotification]) pagination.Result[ExpenseNotificationResponse] {
	out := make([]ExpenseNotificationResponse, len(result.Data))
	for i, n := range result.Data {
		resp := ExpenseNotificationResponse{
			ID:        n.ID,
			ExpenseID: n.ExpenseID,
			ProjectID: n.ProjectID,
			Message:   n.Message,
			CreatedAt: n.CreatedAt.Format(TimeFormat),
		}
		if n.ReadAt != nil {
			s := n.ReadAt.Format(TimeFormat)
			resp.ReadAt = &s
		}
		out[i] = resp
	}
	return pagination.Result[ExpenseNotificationResponse]{
		Data:       out,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}
