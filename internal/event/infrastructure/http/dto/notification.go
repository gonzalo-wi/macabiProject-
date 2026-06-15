package eventdto

import (
	"macabi-back/internal/shared/pagination"
	eventdomain "macabi-back/internal/event/domain"
)

const TimeFormat = "2006-01-02T15:04:05Z07:00"

type EventNotificationResponse struct {
	ID              string  `json:"id"`
	EventInstanceID string  `json:"event_instance_id"`
	Message         string  `json:"message"`
	ReadAt          *string `json:"read_at"`
	CreatedAt       string  `json:"created_at"`
}

func ToEventNotificationListResponse(result pagination.Result[eventdomain.EventNotification]) pagination.Result[EventNotificationResponse] {
	out := make([]EventNotificationResponse, len(result.Data))
	for i, n := range result.Data {
		resp := EventNotificationResponse{
			ID:              n.ID,
			EventInstanceID: n.EventInstanceID,
			Message:         n.Message,
			CreatedAt:       n.CreatedAt.Format(TimeFormat),
		}
		if n.ReadAt != nil {
			s := n.ReadAt.Format(TimeFormat)
			resp.ReadAt = &s
		}
		out[i] = resp
	}
	return pagination.Result[EventNotificationResponse]{
		Data:       out,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}
