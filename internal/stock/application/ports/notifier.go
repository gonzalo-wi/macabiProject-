package stockports

import "context"

// UserNotifier campana in-app + push + email para eventos de pedidos de stock.
type UserNotifier interface {
	NotifyCoordinatorsNewRequest(ctx context.Context, requestID, projectID, resourceName string, quantity int)
	NotifyRequesterApproved(ctx context.Context, requestID, requesterID, resourceName string, quantity int)
	NotifyRequesterRejected(ctx context.Context, requestID, requesterID, resourceName string, quantity int)
}
