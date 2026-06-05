package stockdto

import (
	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type NotificationResponse struct {
	ID        string  `json:"id"`
	RequestID string  `json:"request_id"`
	Message   string  `json:"message"`
	ReadAt    *string `json:"read_at"`
	CreatedAt string  `json:"created_at"`
}

func ToNotificationListResponse(result pagination.Result[stockdomain.StockNotification]) pagination.Result[NotificationResponse] {
	out := make([]NotificationResponse, len(result.Data))
	for i, n := range result.Data {
		resp := NotificationResponse{
			ID:        n.ID,
			RequestID: n.RequestID,
			Message:   n.Message,
			CreatedAt: n.CreatedAt.Format(TimeFormat),
		}
		if n.ReadAt != nil {
			s := n.ReadAt.Format(TimeFormat)
			resp.ReadAt = &s
		}
		out[i] = resp
	}
	return pagination.Result[NotificationResponse]{
		Data:       out,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}
