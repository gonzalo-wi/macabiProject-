package eventhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	eventusecases "macabi-back/internal/event/application/usecases"
	eventdto "macabi-back/internal/event/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

func (h *Handler) SubmitResponse(c *gin.Context) {
	var req eventdto.SubmitResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	userID := c.GetString(userhttp.AuthUserIDKey)
	if err := h.submitResponse.Execute(c.Request.Context(), eventusecases.SubmitResponseInput{
		EventID:   c.Param("id"),
		UserID:    userID,
		ProjectID: req.ProjectID,
		Answers:   req.Answers,
	}); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetMyResponse(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	resp, answers, err := h.getMyResponse.Execute(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToMyResponseResult(resp, answers))
}

func (h *Handler) ListEventResponsesForAdmin(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	result, err := h.listEventResponses.Execute(c.Request.Context(), c.Param("id"), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToAdminResponseListResponse(result))
}

func (h *Handler) GetModuleResponseSummary(c *gin.Context) {
	summary, err := h.getModuleResponseSummary.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToModuleResponseSummaryJSON(summary))
}
