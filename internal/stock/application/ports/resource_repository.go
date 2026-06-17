package stockports

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type ResourceRepository interface {
	SaveResource(ctx context.Context, r *stockdomain.Resource) error
	FindResourceByID(ctx context.Context, id string) (*stockdomain.Resource, error)
	FindAllResources(ctx context.Context, params pagination.Params) (pagination.Result[stockdomain.Resource], error)
	UpdateResource(ctx context.Context, r *stockdomain.Resource) error
	DeleteResource(ctx context.Context, id string) error
	HasActiveRequests(ctx context.Context, resourceID string) (bool, error)
}
