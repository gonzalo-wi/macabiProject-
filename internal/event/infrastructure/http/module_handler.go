package eventhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	eventusecases "macabi-back/internal/event/application/usecases"
	eventdto "macabi-back/internal/event/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
)

func (h *Handler) CreateModule(c *gin.Context) {
	var req eventdto.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	mt, err := eventdto.ParseModuleType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	m, err := h.createModule.Execute(c.Request.Context(), eventusecases.CreateModuleInput{
		EventInstanceID: req.EventInstanceID,
		Title:           req.Title,
		Type:            mt,
		SortOrder:       req.SortOrder,
		IsRequired:      req.IsRequired,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, eventdto.ToModuleResponse(m))
}

func (h *Handler) PatchModule(c *gin.Context) {
	var req eventdto.PatchModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	mt, err := eventdto.ParseModuleType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	m, err := h.updateModule.Execute(c.Request.Context(), eventusecases.UpdateModuleInput{
		ID:         c.Param("id"),
		Title:      req.Title,
		Type:       mt,
		SortOrder:  req.SortOrder,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToModuleResponse(m))
}

func (h *Handler) DeleteModule(c *gin.Context) {
	if err := h.deleteModule.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) SetModuleProjects(c *gin.Context) {
	var req eventdto.SetModuleProjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.setModuleProjects.Execute(c.Request.Context(), c.Param("id"), req.ProjectIDs); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
