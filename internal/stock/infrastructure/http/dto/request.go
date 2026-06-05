package stockdto

import (
	"time"

	"macabi-back/internal/shared/pagination"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockdomain "macabi-back/internal/stock/domain"
)

type CreateRequestRequest struct {
	ProjectID      string  `json:"project_id" binding:"required"`
	ResourceID     string  `json:"resource_id" binding:"required"`
	Quantity       int     `json:"quantity" binding:"required,min=1"`
	WithdrawalDate string  `json:"withdrawal_date" binding:"required"`
	ReturnDate     *string `json:"return_date"`
	Notes          string  `json:"notes"`
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

func ToRequestResponse(r *stockdomain.ResourceRequest) RequestResponse {
	resp := RequestResponse{
		ID:             r.ID,
		ProjectID:      r.ProjectID,
		ResourceID:     r.ResourceID,
		RequestedByID:  r.RequestedByID,
		Quantity:       r.Quantity,
		WithdrawalDate: r.WithdrawalDate.Format(TimeFormat),
		Status:         string(r.Status),
		Notes:          r.Notes,
		CreatedAt:      r.CreatedAt.Format(TimeFormat),
	}
	if r.ReturnDate != nil {
		s := r.ReturnDate.Format(TimeFormat)
		resp.ReturnDate = &s
	}
	return resp
}

func ToRequestDetailResponse(d *stockdomain.RequestDetail) RequestDetailResponse {
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
		WithdrawalDate: d.Request.WithdrawalDate.Format(TimeFormat),
		Status:         string(d.Request.Status),
		Notes:          d.Request.Notes,
		CreatedAt:      d.Request.CreatedAt.Format(TimeFormat),
	}
	if d.Request.ReturnDate != nil {
		s := d.Request.ReturnDate.Format(TimeFormat)
		resp.ReturnDate = &s
	}
	return resp
}

func ToRequestDetailListResponse(result pagination.Result[stockdomain.RequestDetail]) pagination.Result[RequestDetailResponse] {
	items := make([]RequestDetailResponse, len(result.Data))
	for i, d := range result.Data {
		items[i] = ToRequestDetailResponse(&d)
	}
	return pagination.Result[RequestDetailResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

func (r CreateRequestRequest) ToInput(requestedByID, userRole string) (stockusecases.CreateRequestInput, error) {
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
