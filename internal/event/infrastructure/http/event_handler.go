package eventhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	eventusecases "macabi-back/internal/event/application/usecases"
	eventdto "macabi-back/internal/event/infrastructure/http/dto"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
)

func (h *Handler) CreateEvent(c *gin.Context) {
	var req eventdto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	et, err := eventdto.ParseEventType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	st, err := eventdto.ParseEventStatus(req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	e, err := h.createEvent.Execute(c.Request.Context(), eventusecases.CreateEventInput{
		Title:    req.Title,
		Type:     et,
		StartsAt: req.StartsAt,
		Deadline: req.ResponseDeadlineAt,
		Status:   st,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, eventdto.ToEventInstanceJSON(e))
}

func (h *Handler) ListEvents(c *gin.Context) {
	params := pagination.ParseParams(c.Query("page"), c.Query("page_size"))
	res, err := h.listEvents.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.EventListResponse(res))
}

func (h *Handler) GetEventDetail(c *gin.Context) {
	d, err := h.getEventDetail.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToEventDetailResponse(d))
}

func (h *Handler) PatchEvent(c *gin.Context) {
	var req eventdto.PatchEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	et, err := eventdto.ParseEventType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	st, err := eventdto.ParseEventStatus(req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	e, err := h.updateEvent.Execute(c.Request.Context(), eventusecases.UpdateEventInput{
		ID:       c.Param("id"),
		Title:    req.Title,
		Type:     et,
		StartsAt: req.StartsAt,
		Deadline: req.ResponseDeadlineAt,
		Status:   st,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventdto.ToEventInstanceJSON(e))
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	if err := h.deleteEvent.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) SetEventProjects(c *gin.Context) {
	var req eventdto.SetProjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.setEventProjects.Execute(c.Request.Context(), c.Param("id"), req.ProjectIDs); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
