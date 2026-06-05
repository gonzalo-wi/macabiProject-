package stockhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	stockdto "macabi-back/internal/stock/infrastructure/http/dto"
)

func (h *Handler) CreateResource(c *gin.Context) {
	var req stockdto.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	r, err := h.createResourceUC.Execute(c.Request.Context(), req.ToInput())
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, stockdto.ToResourceResponse(r))
}

func (h *Handler) ListResources(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	result, err := h.listResourcesUC.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, stockdto.ToResourceListResponse(result))
}

func (h *Handler) GetResource(c *gin.Context) {
	r, err := h.getResourceUC.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, stockdto.ToResourceResponse(r))
}

func (h *Handler) UpdateResource(c *gin.Context) {
	var req stockdto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	r, err := h.updateResourceUC.Execute(c.Request.Context(), req.ToInput(c.Param("id")))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, stockdto.ToResourceResponse(r))
}

func (h *Handler) DeleteResource(c *gin.Context) {
	if err := h.deleteResourceUC.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
