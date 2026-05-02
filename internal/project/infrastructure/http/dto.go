package projecthttp

import (
	projectusecases "macabi-back/internal/project/application/usecases"
	projectdomain "macabi-back/internal/project/domain"
	"macabi-back/internal/shared/pagination"
)

// --- Requests ---

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	AdminUserID string `json:"admin_user_id" binding:"required"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	AdminUserID string `json:"admin_user_id" binding:"required"`
}

// --- Responses ---

type ProjectResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AdminUserID string `json:"admin_user_id"`
	Active      bool   `json:"active"`
}

func toProjectResponse(p *projectdomain.Project) ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		AdminUserID: p.AdminUserID,
		Active:      p.Active,
	}
}

func toProjectListResponse(result pagination.Result[projectdomain.Project]) pagination.Result[ProjectResponse] {
	items := make([]ProjectResponse, len(result.Data))
	for i, p := range result.Data {
		items[i] = toProjectResponse(&p)
	}
	return pagination.Result[ProjectResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

// --- Input mappers ---

func (r CreateProjectRequest) toInput() projectusecases.CreateProjectInput {
	return projectusecases.CreateProjectInput{
		Name:        r.Name,
		Description: r.Description,
		AdminUserID: r.AdminUserID,
	}
}

func (r UpdateProjectRequest) toInput(id string) projectusecases.UpdateProjectInput {
	return projectusecases.UpdateProjectInput{
		ID:          id,
		Name:        r.Name,
		Description: r.Description,
		AdminUserID: r.AdminUserID,
	}
}
