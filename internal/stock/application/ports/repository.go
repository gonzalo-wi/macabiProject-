package stockports

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type ResourceListFilter struct {
	Query string
}

type RequestListFilter struct {
	Query  string
	Status string
}

type StockRepository interface {
	SaveResource(ctx context.Context, r *stockdomain.Resource) error
	FindResourceByID(ctx context.Context, id string) (*stockdomain.Resource, error)
	FindAllResources(ctx context.Context, filter ResourceListFilter, params pagination.Params) (pagination.Result[stockdomain.Resource], error)
	UpdateResource(ctx context.Context, r *stockdomain.Resource) error
	DeleteResource(ctx context.Context, id string) error
	HasActiveRequests(ctx context.Context, resourceID string) (bool, error)

	SaveRequest(ctx context.Context, r *stockdomain.ResourceRequest) error
	FindRequestByID(ctx context.Context, id string) (*stockdomain.ResourceRequest, error)
	FindRequestDetailByID(ctx context.Context, id string) (*stockdomain.RequestDetail, error)
	ApproveRequest(ctx context.Context, requestID string) error
	UpdateRequestStatus(ctx context.Context, id string, status stockdomain.RequestStatus) error
	ReturnRequest(ctx context.Context, requestID string) error
	ListRequests(ctx context.Context, filter RequestListFilter, params pagination.Params, projectID string, onlyRequestedBy *string) (pagination.Result[stockdomain.RequestDetail], error)
	ListMyRequests(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.RequestDetail], error)

	SaveNotification(ctx context.Context, n *stockdomain.StockNotification) error
	ListNotificationsByUser(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.StockNotification], error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	MarkAllNotificationsRead(ctx context.Context, userID string) (int64, error)
	UnreadCount(ctx context.Context, userID string) (int64, error)
}
