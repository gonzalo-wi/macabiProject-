package newshttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	newsdto "macabi-back/internal/news/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) ListNotifications(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	result, err := h.listNotifs.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, newsdto.ToNewsNotificationListResponse(result))
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	if err := h.markNotifRead.Execute(c.Request.Context(), c.Param("id"), c.GetString(userhttp.AuthUserIDKey)); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) MarkAllNotificationsRead(c *gin.Context) {
	if err := h.markAllNotifs.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey)); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UnreadNotificationCount(c *gin.Context) {
	count, err := h.unreadNotifs.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
