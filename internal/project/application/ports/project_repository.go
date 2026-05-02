package projectports

import (
	"context"

	projectdomain "macabi-back/internal/project/domain"
	"macabi-back/internal/shared/pagination"
)

type ProjectRepository interface {
	Save(ctx context.Context, p *projectdomain.Project) error
	FindByID(ctx context.Context, id string) (*projectdomain.Project, error)
	FindAll(ctx context.Context, params pagination.Params) (pagination.Result[projectdomain.Project], error)
	Update(ctx context.Context, p *projectdomain.Project) error
	Delete(ctx context.Context, id string) error
}
