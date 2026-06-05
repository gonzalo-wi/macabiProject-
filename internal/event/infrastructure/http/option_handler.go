package eventhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	eventusecases "macabi-back/internal/event/application/usecases"
	eventdto "macabi-back/internal/event/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
)

func (h *Handler) CreateOptionGroup(c *gin.Context) {
	var req eventdto.CreateOptionGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	ot, err := eventdto.ParseOptionGroupType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	g, err := h.createOptionGroup.Execute(c.Request.Context(), eventusecases.CreateOptionGroupInput{
		ModuleID:   req.ModuleID,
		Name:       req.Name,
		Type:       ot,
		SortOrder:  req.SortOrder,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, eventdto.ToOptionGroupResponse(g))
}

func (h *Handler) PatchOptionGroup(c *gin.Context) {
	var req eventdto.PatchOptionGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	ot, err := eventdto.ParseOptionGroupType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	g, err := h.updateOptionGroup.Execute(c.Request.Context(), eventusecases.UpdateOptionGroupInput{
		ID:         c.Param("id"),
		Name:       req.Name,
		Type:       ot,
		SortOrder:  req.SortOrder,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToOptionGroupResponse(g))
}

func (h *Handler) DeleteOptionGroup(c *gin.Context) {
	if err := h.deleteOptionGroup.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateOption(c *gin.Context) {
	var req eventdto.CreateOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	o, err := h.createOption.Execute(c.Request.Context(), eventusecases.CreateOptionInput{
		GroupID:     req.GroupID,
		Label:       req.Label,
		MaxCapacity: req.MaxCapacity,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, eventdto.ToOptionResponse(o))
}

func (h *Handler) PatchOption(c *gin.Context) {
	var req eventdto.PatchOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	o, err := h.updateOption.Execute(c.Request.Context(), eventusecases.UpdateOptionInput{
		ID:           c.Param("id"),
		Label:        req.Label,
		MaxCapacity:  req.MaxCapacity,
		SortOrder:    req.SortOrder,
		CurrentCount: req.CurrentCount,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToOptionResponse(o))
}

func (h *Handler) DeleteOption(c *gin.Context) {
	if err := h.deleteOption.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
