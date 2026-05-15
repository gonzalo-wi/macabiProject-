package projecthttp

import (
	projectusecases "macabi-back/internal/project/application/usecases"
	projectdomain "macabi-back/internal/project/domain"
	"macabi-back/internal/shared/pagination"
)

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

type AddProjectMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type ProjectMemberResponse struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at,omitempty"`
}

func toProjectResponse(p *projectdomain.Project) ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Format(timeRFC3339),
	}
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

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

func (r CreateProjectRequest) toInput() projectusecases.CreateProjectInput {
	return projectusecases.CreateProjectInput{
		Name:          r.Name,
		Description:   r.Description,
		CoordinatorID: r.CoordinatorID,
	}
}

func (r UpdateProjectRequest) toInput(id string) projectusecases.UpdateProjectInput {
	return projectusecases.UpdateProjectInput{
		ID:          id,
		Name:        r.Name,
		Description: r.Description,
	}
}

func toMemberResponse(m *projectdomain.ProjectMember) ProjectMemberResponse {
	return ProjectMemberResponse{
		ID:        m.ID,
		ProjectID: m.ProjectID,
		UserID:    m.UserID,
		Role:      string(m.Role),
		CreatedAt: m.CreatedAt.Format(timeRFC3339),
	}
}
