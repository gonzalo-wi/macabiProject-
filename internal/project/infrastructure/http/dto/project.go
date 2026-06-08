package projectdto

import (
	projectusecases "macabi-back/internal/project/application/usecases"
	projectdomain "macabi-back/internal/project/domain"
	"macabi-back/internal/shared/pagination"
)

const TimeRFC3339 = "2006-01-02T15:04:05Z07:00"

type CreateProjectRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	CoordinatorID string `json:"coordinator_id" binding:"required"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type ProjectResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at,omitempty"`
}

func (r CreateProjectRequest) ToInput() projectusecases.CreateProjectInput {
	return projectusecases.CreateProjectInput{
		Name:          r.Name,
		Description:   r.Description,
		CoordinatorID: r.CoordinatorID,
	}
}

func (r UpdateProjectRequest) ToInput(id string) projectusecases.UpdateProjectInput {
	return projectusecases.UpdateProjectInput{
		ID:          id,
		Name:        r.Name,
		Description: r.Description,
	}
}

func ToProjectResponse(p *projectdomain.Project) ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Format(TimeRFC3339),
	}
}

func ToProjectListResponse(result pagination.Result[projectdomain.Project]) pagination.Result[ProjectResponse] {
	items := make([]ProjectResponse, len(result.Data))
	for i, p := range result.Data {
		items[i] = ToProjectResponse(&p)
	}
	return pagination.Result[ProjectResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}
