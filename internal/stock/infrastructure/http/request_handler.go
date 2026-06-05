package stockhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	stockusecases "macabi-back/internal/stock/application/usecases"
	stockdto "macabi-back/internal/stock/infrastructure/http/dto"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) CreateRequest(c *gin.Context) {
	var req stockdto.CreateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	input, err := req.ToInput(c.GetString(userhttp.AuthUserIDKey), c.GetString(userhttp.AuthRoleKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	request, err := h.createRequestUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, stockdto.ToRequestResponse(request))
}

func (h *Handler) ListRequests(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	result, err := h.listRequestsUC.Execute(c.Request.Context(), stockusecases.ListRequestsInput{
		Params:    params,
		ProjectID: c.Query("project_id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, stockdto.ToRequestDetailListResponse(result))
}

func (h *Handler) ListMyRequests(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	result, err := h.listMyRequestsUC.Execute(c.Request.Context(), c.GetString(userhttp.AuthUserIDKey), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, stockdto.ToRequestDetailListResponse(result))
}

func (h *Handler) GetRequestDetail(c *gin.Context) {
	detail, err := h.getRequestDetailUC.Execute(c.Request.Context(), stockusecases.GetRequestDetailInput{
		RequestID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, stockdto.ToRequestDetailResponse(detail))
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	if err := h.approveRequestUC.Execute(c.Request.Context(), stockusecases.ApproveRequestInput{
		RequestID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) RejectRequest(c *gin.Context) {
	if err := h.rejectRequestUC.Execute(c.Request.Context(), stockusecases.RejectRequestInput{
		RequestID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CancelRequest(c *gin.Context) {
	if err := h.cancelRequestUC.Execute(c.Request.Context(), stockusecases.CancelRequestInput{
		RequestID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeliverRequest(c *gin.Context) {
	if err := h.deliverRequestUC.Execute(c.Request.Context(), stockusecases.DeliverRequestInput{
		RequestID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ReturnRequest(c *gin.Context) {
	if err := h.returnRequestUC.Execute(c.Request.Context(), stockusecases.ReturnRequestInput{
		RequestID: c.Param("id"),
		UserID:    c.GetString(userhttp.AuthUserIDKey),
		UserRole:  c.GetString(userhttp.AuthRoleKey),
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
