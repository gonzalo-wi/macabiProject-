package stockusecases

import (
	"context"
	"strings"

	stockports "macabi-back/internal/stock/application/ports"
)

type UnregisterPushSubscription struct {
	subRepo stockports.PushSubscriptionRepository
}

func NewUnregisterPushSubscription(subRepo stockports.PushSubscriptionRepository) *UnregisterPushSubscription {
	return &UnregisterPushSubscription{subRepo: subRepo}
}

func (uc *UnregisterPushSubscription) Execute(ctx context.Context, endpoint string) error {
	if strings.TrimSpace(endpoint) == "" {
		return nil
	}
	return uc.subRepo.DeleteByEndpoint(ctx, endpoint)
}
