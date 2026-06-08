package projectdto

import projectdomain "macabi-back/internal/project/domain"

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

func ToMemberResponse(m *projectdomain.ProjectMember) ProjectMemberResponse {
	return ProjectMemberResponse{
		ID:        m.ID,
		ProjectID: m.ProjectID,
		UserID:    m.UserID,
		Role:      string(m.Role),
		CreatedAt: m.CreatedAt.Format(TimeRFC3339),
	}
}
