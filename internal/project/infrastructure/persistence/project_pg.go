package projectpersistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	projectdomain "macabi-back/internal/project/domain"
	projectports "macabi-back/internal/project/application/ports"
	"macabi-back/internal/shared/pagination"
)

type ProjectModel struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time
}

func (ProjectModel) TableName() string { return "projects" }

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

func applyProjectListFilter(q *gorm.DB, filter projectports.ProjectListFilter) *gorm.DB {
	if s := strings.TrimSpace(filter.Query); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("(LOWER(name) LIKE ? OR LOWER(COALESCE(description, '')) LIKE ?)", like, like)
	}
	return q
}

func (r *ProjectRepositoryPG) FindAll(ctx context.Context, filter projectports.ProjectListFilter, params pagination.Params) (pagination.Result[projectdomain.Project], error) {
	base := applyProjectListFilter(r.db.WithContext(ctx).Model(&ProjectModel{}), filter)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return pagination.Result[projectdomain.Project]{}, err
	}
	var models []ProjectModel
	if err := applyProjectListFilter(r.db.WithContext(ctx).Model(&ProjectModel{}), filter).
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
