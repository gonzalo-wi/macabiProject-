package newsports

import (
	"context"

	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"
)

type NewsRepository interface {
	Save(ctx context.Context, n *newsdomain.News) error
	Update(ctx context.Context, n *newsdomain.News) error
	DeleteByID(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*newsdomain.News, error)
	// SetNewsProjects reemplaza la asociación de proyectos de una noticia (vacío = todos).
	SetNewsProjects(ctx context.Context, newsID string, projectIDs []string) error
	// FindLatestPublished devuelve la última noticia publicada VISIBLE para los proyectos
	// del viewer (vacío = solo las dirigidas a todos), o nil si no hay ninguna.
	FindLatestPublished(ctx context.Context, viewerProjectIDs []string) (*newsdomain.News, error)
	// ListPublished pagina noticias publicadas visibles para los proyectos del viewer.
	ListPublished(ctx context.Context, viewerProjectIDs []string, params pagination.Params) (pagination.Result[newsdomain.News], error)
	// ListAll pagina todas (incluye borradores) para administración, sin filtro.
	ListAll(ctx context.Context, params pagination.Params) (pagination.Result[newsdomain.News], error)
}
