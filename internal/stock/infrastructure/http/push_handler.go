package stockhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharederrors "macabi-back/internal/shared/errors"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockdto "macabi-back/internal/stock/infrastructure/http/dto"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) GetVAPIDPublicKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"public_key": h.vapidPublicKey})
}

func (h *Handler) RegisterPushSubscription(c *gin.Context) {
	var req stockdto.RegisterPushSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.registerPushSubUC.Execute(c.Request.Context(), stockusecases.RegisterPushSubscriptionInput{
		UserID:   c.GetString(userhttp.AuthUserIDKey),
		Endpoint: req.Endpoint,
		P256DH:   req.P256DH,
		Auth:     req.Auth,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UnregisterPushSubscription(c *gin.Context) {
	var req stockdto.UnregisterPushSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.unregisterPushSubUC.Execute(c.Request.Context(), req.Endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
