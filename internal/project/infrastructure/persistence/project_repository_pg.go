package projectpersistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
	"macabi-back/internal/shared/pagination"
)

type ProjectModel struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time
}

func (ProjectModel) TableName() string { return "projects" }

type ProjectMemberModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID string `gorm:"type:uuid;not null"`
	UserID    string `gorm:"type:uuid;not null"`
	Role      string `gorm:"not null;default:'madrij'"`
	CreatedAt time.Time
}

func (ProjectMemberModel) TableName() string { return "project_members" }

type ProjectRepositoryPG struct {
	db *gorm.DB
}

func NewProjectRepositoryPG(db *gorm.DB) *ProjectRepositoryPG {
	return &ProjectRepositoryPG{db: db}
}

func RunMigrations(_ *gorm.DB) error {
	return nil
}

func (r *ProjectRepositoryPG) Save(ctx context.Context, p *projectdomain.Project) error {
	m := toProjectModel(p)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	return nil
}

func (r *ProjectRepositoryPG) FindByID(ctx context.Context, id string) (*projectdomain.Project, error) {
	var m ProjectModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, projectdomain.ErrProjectNotFound
		}
		return nil, err
	}
	return toDomainProject(m), nil
}

func (r *ProjectRepositoryPG) FindAll(ctx context.Context, params pagination.Params) (pagination.Result[projectdomain.Project], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ProjectModel{}).Count(&total).Error; err != nil {
		return pagination.Result[projectdomain.Project]{}, err
	}

	var models []ProjectModel
	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error; err != nil {
		return pagination.Result[projectdomain.Project]{}, err
	}

	items := make([]projectdomain.Project, len(models))
	for i, m := range models {
		items[i] = *toDomainProject(m)
	}
	return pagination.NewResult(items, total, params), nil
}

func (r *ProjectRepositoryPG) Update(ctx context.Context, p *projectdomain.Project) error {
	return r.db.WithContext(ctx).Model(&ProjectModel{}).
		Where("id = ?", p.ID).
		Updates(map[string]interface{}{
			"name":        p.Name,
			"description": p.Description,
		}).Error
}

func (r *ProjectRepositoryPG) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", id).Delete(&ProjectMemberModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&ProjectModel{}, "id = ?", id).Error
	})
}

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

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// Postgres unique_violation 23505
	return strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key")
}

var _ projectports.ProjectRepository = (*ProjectRepositoryPG)(nil)

func toProjectModel(p *projectdomain.Project) ProjectModel {
	return ProjectModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}

func toDomainProject(m ProjectModel) *projectdomain.Project {
	return &projectdomain.Project{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
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
