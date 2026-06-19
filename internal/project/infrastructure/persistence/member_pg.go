package projectpersistence

import (
	"context"
	"time"

	projectdomain "macabi-back/internal/project/domain"
)

type ProjectMemberModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID string `gorm:"type:uuid;not null;uniqueIndex:idx_project_member"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex:idx_project_member"`
	Role      string `gorm:"not null;default:'madrij'"`
	CreatedAt time.Time
}

func (ProjectMemberModel) TableName() string { return "project_members" }

func (r *ProjectRepositoryPG) ListMembers(ctx context.Context, projectID string) ([]projectdomain.ProjectMember, error) {
	var rows []ProjectMemberModel
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]projectdomain.ProjectMember, len(rows))
	for i, row := range rows {
		out[i] = *toDomainMember(row)
	}
	return out, nil
}

func (r *ProjectRepositoryPG) AddMember(ctx context.Context, m *projectdomain.ProjectMember) error {
	row := ProjectMemberModel{
		ProjectID: m.ProjectID,
		UserID:    m.UserID,
		Role:      string(m.Role),
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		if isUniqueViolation(err) {
			return projectdomain.ErrDuplicateMember
		}
		return err
	}
	m.ID = row.ID
	m.CreatedAt = row.CreatedAt
	return nil
}

func (r *ProjectRepositoryPG) RemoveMember(ctx context.Context, projectID, userID string) error {
	res := r.db.WithContext(ctx).Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&ProjectMemberModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return projectdomain.ErrMemberNotFound
	}
	return nil
}

func (r *ProjectRepositoryPG) RemoveAllMembersByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&ProjectMemberModel{}).Error
}

func toDomainMember(m ProjectMemberModel) *projectdomain.ProjectMember {
	return &projectdomain.ProjectMember{
		ID:        m.ID,
		ProjectID: m.ProjectID,
		UserID:    m.UserID,
		Role:      projectdomain.ProjectMemberRole(m.Role),
		CreatedAt: m.CreatedAt,
	}
}
