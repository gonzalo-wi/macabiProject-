package stockhttp

import (
	"time"

	"macabi-back/internal/shared/pagination"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockdomain "macabi-back/internal/stock/domain"
)

const timeFormat = "2006-01-02T15:04:05Z07:00"

type CreateResourceRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	TotalStock int    `json:"total_stock" binding:"required,min=0"`
}

type UpdateResourceRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	TotalStock int    `json:"total_stock" binding:"min=0"`
}

type CreateRequestRequest struct {
	ProjectID      string  `json:"project_id" binding:"required"`
	ResourceID     string  `json:"resource_id" binding:"required"`
	Quantity       int     `json:"quantity" binding:"required,min=1"`
	WithdrawalDate string  `json:"withdrawal_date" binding:"required"`
	ReturnDate     *string `json:"return_date"`
	Notes          string  `json:"notes"`
}

type ResourceResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	TotalStock     int    `json:"total_stock"`
	AvailableStock int    `json:"available_stock"`
	CreatedAt      string `json:"created_at"`
}

type RequestResponse struct {
	ID             string  `json:"id"`
	ProjectID      string  `json:"project_id"`
	ResourceID     string  `json:"resource_id"`
	RequestedByID  string  `json:"requested_by_id"`
	Quantity       int     `json:"quantity"`
	WithdrawalDate string  `json:"withdrawal_date"`
	ReturnDate     *string `json:"return_date"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes"`
	CreatedAt      string  `json:"created_at"`
}

type RequestDetailResponse struct {
	ID             string  `json:"id"`
	ProjectID      string  `json:"project_id"`
	ProjectName    string  `json:"project_name"`
	ResourceID     string  `json:"resource_id"`
	ResourceName   string  `json:"resource_name"`
	ResourceType   string  `json:"resource_type"`
	RequestedByID  string  `json:"requested_by_id"`
	RequesterName  string  `json:"requester_name"`
	Quantity       int     `json:"quantity"`
	WithdrawalDate string  `json:"withdrawal_date"`
	ReturnDate     *string `json:"return_date"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes"`
	CreatedAt      string  `json:"created_at"`
}

type NotificationResponse struct {
	ID        string  `json:"id"`
	RequestID string  `json:"request_id"`
	Message   string  `json:"message"`
	ReadAt    *string `json:"read_at"`
	CreatedAt string  `json:"created_at"`
}

func toResourceResponse(r *stockdomain.Resource) ResourceResponse {
	return ResourceResponse{
		ID:             r.ID,
		Name:           r.Name,
		Type:           string(r.Type),
		TotalStock:     r.TotalStock,
		AvailableStock: r.AvailableStock,
		CreatedAt:      r.CreatedAt.Format(timeFormat),
	}
}

func toResourceListResponse(result pagination.Result[stockdomain.Resource]) pagination.Result[ResourceResponse] {
	items := make([]ResourceResponse, len(result.Data))
	for i, r := range result.Data {
		items[i] = toResourceResponse(&r)
	}
	return pagination.Result[ResourceResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

func toRequestResponse(r *stockdomain.ResourceRequest) RequestResponse {
	resp := RequestResponse{
		ID:             r.ID,
		ProjectID:      r.ProjectID,
		ResourceID:     r.ResourceID,
		RequestedByID:  r.RequestedByID,
		Quantity:       r.Quantity,
		WithdrawalDate: r.WithdrawalDate.Format(timeFormat),
		Status:         string(r.Status),
		Notes:          r.Notes,
		CreatedAt:      r.CreatedAt.Format(timeFormat),
	}
	if r.ReturnDate != nil {
		s := r.ReturnDate.Format(timeFormat)
		resp.ReturnDate = &s
	}
	return resp
}

func toRequestDetailResponse(d *stockdomain.RequestDetail) RequestDetailResponse {
	resp := RequestDetailResponse{
		ID:             d.Request.ID,
		ProjectID:      d.Request.ProjectID,
		ProjectName:    d.ProjectName,
		ResourceID:     d.Request.ResourceID,
		ResourceName:   d.ResourceName,
		ResourceType:   string(d.ResourceType),
		RequestedByID:  d.Request.RequestedByID,
		RequesterName:  d.RequesterName,
		Quantity:       d.Request.Quantity,
		WithdrawalDate: d.Request.WithdrawalDate.Format(timeFormat),
		Status:         string(d.Request.Status),
		Notes:          d.Request.Notes,
		CreatedAt:      d.Request.CreatedAt.Format(timeFormat),
	}
	if d.Request.ReturnDate != nil {
		s := d.Request.ReturnDate.Format(timeFormat)
		resp.ReturnDate = &s
	}
	return resp
}

func toRequestDetailListResponse(result pagination.Result[stockdomain.RequestDetail]) pagination.Result[RequestDetailResponse] {
	items := make([]RequestDetailResponse, len(result.Data))
	for i, d := range result.Data {
		items[i] = toRequestDetailResponse(&d)
	}
	return pagination.Result[RequestDetailResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

func toNotificationListResponse(result pagination.Result[stockdomain.StockNotification]) pagination.Result[NotificationResponse] {
	out := make([]NotificationResponse, len(result.Data))
	for i, n := range result.Data {
		resp := NotificationResponse{
			ID:        n.ID,
			RequestID: n.RequestID,
			Message:   n.Message,
			CreatedAt: n.CreatedAt.Format(timeFormat),
		}
		if n.ReadAt != nil {
			s := n.ReadAt.Format(timeFormat)
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

func (r CreateResourceRequest) toInput() stockusecases.CreateResourceInput {
	return stockusecases.CreateResourceInput{
		Name:       r.Name,
		Type:       r.Type,
		TotalStock: r.TotalStock,
	}
}

func (r UpdateResourceRequest) toInput(id string) stockusecases.UpdateResourceInput {
	return stockusecases.UpdateResourceInput{
		ID:         id,
		Name:       r.Name,
		Type:       r.Type,
		TotalStock: r.TotalStock,
	}
}

func (r CreateRequestRequest) toInput(requestedByID, userRole string) (stockusecases.CreateRequestInput, error) {
	withdrawalDate, err := time.Parse(time.RFC3339, r.WithdrawalDate)
	if err != nil {
		return stockusecases.CreateRequestInput{}, err
	}
	var returnDate *time.Time
	if r.ReturnDate != nil {
		rd, err := time.Parse(time.RFC3339, *r.ReturnDate)
		if err != nil {
			return stockusecases.CreateRequestInput{}, err
		}
		returnDate = &rd
	}
	return stockusecases.CreateRequestInput{
		ProjectID:      r.ProjectID,
		ResourceID:     r.ResourceID,
		RequestedByID:  requestedByID,
		UserRole:       userRole,
		Quantity:       r.Quantity,
		WithdrawalDate: withdrawalDate,
		ReturnDate:     returnDate,
		Notes:          r.Notes,
	}, nil
}
