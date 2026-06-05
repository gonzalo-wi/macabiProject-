package stockports

import "context"

type StockMailer interface {
	NotifyCoordinatorsNewRequest(ctx context.Context, coordinatorEmails []string, resourceName string, quantity int, requestID string) error
	NotifyRequesterApproved(ctx context.Context, requesterEmail, resourceName string, quantity int) error
	NotifyRequesterRejected(ctx context.Context, requesterEmail, resourceName string, quantity int) error
}
