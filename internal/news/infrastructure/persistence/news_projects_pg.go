package newspersistence

import (
	"context"

	"gorm.io/gorm"
)

type NewsProjectModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	NewsID    string `gorm:"type:uuid;not null;index"`
	ProjectID string `gorm:"type:uuid;not null"`
}

func (NewsProjectModel) TableName() string { return "news_projects" }

// SetNewsProjects reemplaza la asociación de proyectos de una noticia.
func (r *NewsRepositoryPG) SetNewsProjects(ctx context.Context, newsID string, projectIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("news_id = ?", newsID).Delete(&NewsProjectModel{}).Error; err != nil {
			return err
		}
		for _, pid := range projectIDs {
			row := NewsProjectModel{NewsID: newsID, ProjectID: pid}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// projectIDsByNews carga, en una sola query, los project_ids de un conjunto de noticias.
func (r *NewsRepositoryPG) projectIDsByNews(ctx context.Context, newsIDs []string) (map[string][]string, error) {
	out := make(map[string][]string)
	if len(newsIDs) == 0 {
		return out, nil
	}
	var rows []NewsProjectModel
	if err := r.db.WithContext(ctx).
		Where("news_id IN ?", newsIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.NewsID] = append(out[row.NewsID], row.ProjectID)
	}
	return out, nil
}

// projectIDsFor carga los project_ids de una sola noticia.
func (r *NewsRepositoryPG) projectIDsFor(ctx context.Context, newsID string) ([]string, error) {
	var pids []string
	err := r.db.WithContext(ctx).
		Model(&NewsProjectModel{}).
		Where("news_id = ?", newsID).
		Pluck("project_id", &pids).Error
	if err != nil {
		return nil, err
	}
	return pids, nil
}
