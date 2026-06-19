package newsdto

import (
	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"
)

type NewsNotificationResponse struct {
	ID        string  `json:"id"`
	NewsID    string  `json:"news_id"`
	Message   string  `json:"message"`
	ReadAt    *string `json:"read_at"`
	CreatedAt string  `json:"created_at"`
}

func ToNewsNotificationListResponse(result pagination.Result[newsdomain.NewsNotification]) pagination.Result[NewsNotificationResponse] {
	out := make([]NewsNotificationResponse, len(result.Data))
	for i, n := range result.Data {
		resp := NewsNotificationResponse{
			ID:        n.ID,
			NewsID:    n.NewsID,
			Message:   n.Message,
			CreatedAt: n.CreatedAt.Format(TimeFormat),
		}
		if n.ReadAt != nil {
			s := n.ReadAt.Format(TimeFormat)
			resp.ReadAt = &s
		}
		out[i] = resp
	}
	return pagination.Result[NewsNotificationResponse]{
		Data:       out,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}
