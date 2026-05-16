package stockusecases

import (
	"context"
	"strings"

	stockports "macabi-back/internal/stock/application/ports"
)

type RegisterPushSubscriptionInput struct {
	UserID   string
	Endpoint string
	P256DH   string
	Auth     string
}

type RegisterPushSubscription struct {
	subRepo stockports.PushSubscriptionRepository
}

func NewRegisterPushSubscription(subRepo stockports.PushSubscriptionRepository) *RegisterPushSubscription {
	return &RegisterPushSubscription{subRepo: subRepo}
}

func (uc *RegisterPushSubscription) Execute(ctx context.Context, in RegisterPushSubscriptionInput) error {
	if strings.TrimSpace(in.Endpoint) == "" || strings.TrimSpace(in.P256DH) == "" || strings.TrimSpace(in.Auth) == "" {
		return nil
	}
	return uc.subRepo.Upsert(ctx, in.UserID, in.Endpoint, in.P256DH, in.Auth)
}
