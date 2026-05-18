package stockports

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type StockRepository interface {
	// --- Resources ---
	SaveResource(ctx context.Context, r *stockdomain.Resource) error
	FindResourceByID(ctx context.Context, id string) (*stockdomain.Resource, error)
	FindAllResources(ctx context.Context, params pagination.Params) (pagination.Result[stockdomain.Resource], error)
	UpdateResource(ctx context.Context, r *stockdomain.Resource) error
	DeleteResource(ctx context.Context, id string) error
	HasActiveRequests(ctx context.Context, resourceID string) (bool, error)

	// --- Requests ---
	SaveRequest(ctx context.Context, r *stockdomain.ResourceRequest) error
	FindRequestByID(ctx context.Context, id string) (*stockdomain.ResourceRequest, error)
	FindRequestDetailByID(ctx context.Context, id string) (*stockdomain.RequestDetail, error)
	// ApproveRequest atomically sets status to RESERVADO and decrements available_stock.
	ApproveRequest(ctx context.Context, requestID string) error
	UpdateRequestStatus(ctx context.Context, id string, status stockdomain.RequestStatus) error
	// ReturnRequest atomically sets status to DEVUELTO and increments available_stock.
	ReturnRequest(ctx context.Context, requestID string) error
	// ListRequests lists requests filtered by optional project ID; if onlyRequestedBy is non-nil, restrict to submissions by that user.
	ListRequests(ctx context.Context, params pagination.Params, projectID string, onlyRequestedBy *string) (pagination.Result[stockdomain.RequestDetail], error)
	ListMyRequests(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.RequestDetail], error)

	// --- Notifications ---
	SaveNotification(ctx context.Context, n *stockdomain.StockNotification) error
	ListNotificationsByUser(ctx context.Context, userID string) ([]stockdomain.StockNotification, error)
	MarkNotificationRead(ctx context.Context, id, userID string) error
	UnreadCount(ctx context.Context, userID string) (int64, error)
}
