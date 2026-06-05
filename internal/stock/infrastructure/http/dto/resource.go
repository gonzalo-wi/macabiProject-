package stockdto

import (
	"macabi-back/internal/shared/pagination"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockdomain "macabi-back/internal/stock/domain"
)

const TimeFormat = "2006-01-02T15:04:05Z07:00"

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

type ResourceResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	TotalStock     int    `json:"total_stock"`
	AvailableStock int    `json:"available_stock"`
	CreatedAt      string `json:"created_at"`
}

func ToResourceResponse(r *stockdomain.Resource) ResourceResponse {
	return ResourceResponse{
		ID:             r.ID,
		Name:           r.Name,
		Type:           string(r.Type),
		TotalStock:     r.TotalStock,
		AvailableStock: r.AvailableStock,
		CreatedAt:      r.CreatedAt.Format(TimeFormat),
	}
}

func ToResourceListResponse(result pagination.Result[stockdomain.Resource]) pagination.Result[ResourceResponse] {
	items := make([]ResourceResponse, len(result.Data))
	for i, r := range result.Data {
		items[i] = ToResourceResponse(&r)
	}
	return pagination.Result[ResourceResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

func (r CreateResourceRequest) ToInput() stockusecases.CreateResourceInput {
	return stockusecases.CreateResourceInput{
		Name:       r.Name,
		Type:       r.Type,
		TotalStock: r.TotalStock,
	}
}

func (r UpdateResourceRequest) ToInput(id string) stockusecases.UpdateResourceInput {
	return stockusecases.UpdateResourceInput{
		ID:         id,
		Name:       r.Name,
		Type:       r.Type,
		TotalStock: r.TotalStock,
	}
}
