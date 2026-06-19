package newspersistence

import (
	"context"
	"errors"
	"time"

	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"

	"gorm.io/gorm"
)

type NewsModel struct {
	ID               string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title            string `gorm:"not null"`
	Body             string `gorm:"type:text;not null"`
	ImageStoragePath *string `gorm:"type:text"`
	Status           string `gorm:"not null;default:'draft';index"`
	AuthorID         string `gorm:"type:uuid;not null"`
	PublishedAt      *time.Time `gorm:"index"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (NewsModel) TableName() string { return "news" }

type NewsRepositoryPG struct {
	db *gorm.DB
}

func NewNewsRepositoryPG(db *gorm.DB) *NewsRepositoryPG {
	return &NewsRepositoryPG{db: db}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&NewsModel{},
		&NewsNotificationModel{},
		&NewsProjectModel{},
	)
}

func (r *NewsRepositoryPG) Save(ctx context.Context, n *newsdomain.News) error {
	m := toNewsModel(n)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	n.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *NewsRepositoryPG) Update(ctx context.Context, n *newsdomain.News) error {
	// Map para poder fijar image_storage_path / published_at a NULL cuando corresponda.
	return r.db.WithContext(ctx).Model(&NewsModel{}).Where("id = ?", n.ID).Updates(map[string]interface{}{
		"title":              n.Title,
		"body":               n.Body,
		"image_storage_path": n.ImageStoragePath,
		"status":             string(n.Status),
		"published_at":       n.PublishedAt,
	}).Error
}

func (r *NewsRepositoryPG) DeleteByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&NewsModel{}).Error
}

func (r *NewsRepositoryPG) FindByID(ctx context.Context, id string) (*newsdomain.News, error) {
	var m NewsModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newsdomain.ErrNewsNotFound
		}
		return nil, err
	}
	n := toNewsDomain(&m)
	pids, err := r.projectIDsFor(ctx, id)
	if err != nil {
		return nil, err
	}
	n.ProjectIDs = pids
	return n, nil
}

// applyVisibility limita a noticias visibles para los proyectos del viewer:
// sin proyectos asociados (todos) OR asociada a alguno de los del viewer.
func applyVisibility(q *gorm.DB, viewerProjectIDs []string) *gorm.DB {
	if len(viewerProjectIDs) > 0 {
		return q.Where(
			"(NOT EXISTS (SELECT 1 FROM news_projects np WHERE np.news_id = news.id)"+
				" OR EXISTS (SELECT 1 FROM news_projects np WHERE np.news_id = news.id AND np.project_id IN ?))",
			viewerProjectIDs,
		)
	}
	return q.Where("NOT EXISTS (SELECT 1 FROM news_projects np WHERE np.news_id = news.id)")
}

func (r *NewsRepositoryPG) FindLatestPublished(ctx context.Context, viewerProjectIDs []string) (*newsdomain.News, error) {
	q := r.db.WithContext(ctx).Model(&NewsModel{}).Where("status = ?", string(newsdomain.StatusPublished))
	q = applyVisibility(q, viewerProjectIDs)
	var m NewsModel
	err := q.Order("published_at DESC").First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	n := toNewsDomain(&m)
	pids, err := r.projectIDsFor(ctx, n.ID)
	if err != nil {
		return nil, err
	}
	n.ProjectIDs = pids
	return n, nil
}

func (r *NewsRepositoryPG) ListPublished(ctx context.Context, viewerProjectIDs []string, params pagination.Params) (pagination.Result[newsdomain.News], error) {
	return r.list(ctx, params, true, viewerProjectIDs)
}

func (r *NewsRepositoryPG) ListAll(ctx context.Context, params pagination.Params) (pagination.Result[newsdomain.News], error) {
	return r.list(ctx, params, false, nil)
}

// list pagina noticias. onlyPublished aplica el filtro de estado + visibilidad por
// proyecto del viewer; con onlyPublished=false (admin) no se filtra nada.
func (r *NewsRepositoryPG) list(ctx context.Context, params pagination.Params, onlyPublished bool, viewerProjectIDs []string) (pagination.Result[newsdomain.News], error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&NewsModel{})
		if onlyPublished {
			q = q.Where("status = ?", string(newsdomain.StatusPublished))
			q = applyVisibility(q, viewerProjectIDs)
		}
		return q
	}

	var total int64
	if err := build().Count(&total).Error; err != nil {
		return pagination.Result[newsdomain.News]{}, err
	}

	query := build()
	if onlyPublished {
		query = query.Order("published_at DESC")
	} else {
		query = query.Order("created_at DESC")
	}

	var models []NewsModel
	if err := query.Offset(params.Offset()).Limit(params.PageSize).Find(&models).Error; err != nil {
		return pagination.Result[newsdomain.News]{}, err
	}

	ids := make([]string, len(models))
	for i := range models {
		ids[i] = models[i].ID
	}
	projByNews, err := r.projectIDsByNews(ctx, ids)
	if err != nil {
		return pagination.Result[newsdomain.News]{}, err
	}

	out := make([]newsdomain.News, len(models))
	for i := range models {
		out[i] = *toNewsDomain(&models[i])
		out[i].ProjectIDs = projByNews[models[i].ID]
	}
	return pagination.NewResult(out, total, params), nil
}

func toNewsModel(n *newsdomain.News) NewsModel {
	return NewsModel{
		ID:               n.ID,
		Title:            n.Title,
		Body:             n.Body,
		ImageStoragePath: n.ImageStoragePath,
		Status:           string(n.Status),
		AuthorID:         n.AuthorID,
		PublishedAt:      n.PublishedAt,
	}
}

func toNewsDomain(m *NewsModel) *newsdomain.News {
	return &newsdomain.News{
		ID:               m.ID,
		Title:            m.Title,
		Body:             m.Body,
		ImageStoragePath: m.ImageStoragePath,
		Status:           newsdomain.NewsStatus(m.Status),
		AuthorID:         m.AuthorID,
		PublishedAt:      m.PublishedAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
