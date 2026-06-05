package stockdto

type RegisterPushSubscriptionRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
	P256DH   string `json:"p256dh"   binding:"required"`
	Auth     string `json:"auth"     binding:"required"`
}

type UnregisterPushSubscriptionRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
}
