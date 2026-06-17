package stockports

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type RequestRepository interface {
	SaveRequest(ctx context.Context, r *stockdomain.ResourceRequest) error
	FindRequestByID(ctx context.Context, id string) (*stockdomain.ResourceRequest, error)
	FindRequestDetailByID(ctx context.Context, id string) (*stockdomain.RequestDetail, error)
	ApproveRequest(ctx context.Context, requestID string) error
	UpdateRequestStatus(ctx context.Context, id string, status stockdomain.RequestStatus) error
	ReturnRequest(ctx context.Context, requestID string) error
	ListRequests(ctx context.Context, params pagination.Params, projectID string, onlyRequestedBy *string) (pagination.Result[stockdomain.RequestDetail], error)
	ListMyRequests(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.RequestDetail], error)
}
