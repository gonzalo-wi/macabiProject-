package projectpersistence

import (
	"context"
	"errors"
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
	AdminUserID string `gorm:"type:uuid;not null"`
	Active      bool   `gorm:"not null;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ProjectModel) TableName() string { return "projects" }

type ProjectRepositoryPG struct {
	db *gorm.DB
}

func NewProjectRepositoryPG(db *gorm.DB) *ProjectRepositoryPG {
	return &ProjectRepositoryPG{db: db}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&ProjectModel{})
}

func (r *ProjectRepositoryPG) Save(ctx context.Context, p *projectdomain.Project) error {
	m := toProjectModel(p)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
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
			"name":          p.Name,
			"description":   p.Description,
			"admin_user_id": p.AdminUserID,
			"active":        p.Active,
		}).Error
}

func (r *ProjectRepositoryPG) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&ProjectModel{}, "id = ?", id).Error
}

// Ensure interface compliance
var _ projectports.ProjectRepository = (*ProjectRepositoryPG)(nil)

func toProjectModel(p *projectdomain.Project) ProjectModel {
	return ProjectModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		AdminUserID: p.AdminUserID,
		Active:      p.Active,
	}
}

func toDomainProject(m ProjectModel) *projectdomain.Project {
	return &projectdomain.Project{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		AdminUserID: m.AdminUserID,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
