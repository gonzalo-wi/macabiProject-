package eventhttp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	eventusecases "macabi-back/internal/event/application/usecases"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

type Handler struct {
	createEvent              *eventusecases.CreateEvent
	listEvents               *eventusecases.ListEvents
	getEventDetail           *eventusecases.GetEventDetail
	updateEvent              *eventusecases.UpdateEvent
	deleteEvent              *eventusecases.DeleteEvent
	setEventProjects         *eventusecases.SetEventProjects
	createModule             *eventusecases.CreateModule
	updateModule             *eventusecases.UpdateModule
	deleteModule             *eventusecases.DeleteModule
	setModuleProjects        *eventusecases.SetModuleProjects
	createOptionGroup        *eventusecases.CreateOptionGroup
	updateOptionGroup        *eventusecases.UpdateOptionGroup
	deleteOptionGroup        *eventusecases.DeleteOptionGroup
	createOption             *eventusecases.CreateOption
	updateOption             *eventusecases.UpdateOption
	deleteOption             *eventusecases.DeleteOption
	submitResponse           *eventusecases.SubmitResponse
	getMyResponse            *eventusecases.GetMyResponse
	listEventResponses       *eventusecases.ListEventResponses
	getModuleResponseSummary *eventusecases.GetModuleResponseSummary
}

func NewHandler(
	createEvent *eventusecases.CreateEvent,
	listEvents *eventusecases.ListEvents,
	getEventDetail *eventusecases.GetEventDetail,
	updateEvent *eventusecases.UpdateEvent,
	deleteEvent *eventusecases.DeleteEvent,
	setEventProjects *eventusecases.SetEventProjects,
	createModule *eventusecases.CreateModule,
	updateModule *eventusecases.UpdateModule,
	deleteModule *eventusecases.DeleteModule,
	setModuleProjects *eventusecases.SetModuleProjects,
	createOptionGroup *eventusecases.CreateOptionGroup,
	updateOptionGroup *eventusecases.UpdateOptionGroup,
	deleteOptionGroup *eventusecases.DeleteOptionGroup,
	createOption *eventusecases.CreateOption,
	updateOption *eventusecases.UpdateOption,
	deleteOption *eventusecases.DeleteOption,
	submitResponse *eventusecases.SubmitResponse,
	getMyResponse *eventusecases.GetMyResponse,
	listEventResponses *eventusecases.ListEventResponses,
	getModuleResponseSummary *eventusecases.GetModuleResponseSummary,
) *Handler {
	return &Handler{
		createEvent:              createEvent,
		listEvents:               listEvents,
		getEventDetail:           getEventDetail,
		updateEvent:              updateEvent,
		deleteEvent:              deleteEvent,
		setEventProjects:         setEventProjects,
		createModule:             createModule,
		updateModule:             updateModule,
		deleteModule:             deleteModule,
		setModuleProjects:        setModuleProjects,
		createOptionGroup:        createOptionGroup,
		updateOptionGroup:        updateOptionGroup,
		deleteOptionGroup:        deleteOptionGroup,
		createOption:             createOption,
		updateOption:             updateOption,
		deleteOption:             deleteOption,
		submitResponse:           submitResponse,
		getMyResponse:            getMyResponse,
		listEventResponses:       listEventResponses,
		getModuleResponseSummary: getModuleResponseSummary,
	}
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var req createEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	et, err := parseEventType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	st, err := parseEventStatus(req.Status)
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
	c.JSON(http.StatusCreated, toEventInstanceJSON(e))
}

func (h *Handler) ListEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	params := pagination.NewParams(page, pageSize)
	res, err := h.listEvents.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventListResponse(res))
}

func (h *Handler) GetEventDetail(c *gin.Context) {
	d, err := h.getEventDetail.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toEventDetailResponse(d))
}

func (h *Handler) PatchEvent(c *gin.Context) {
	var req patchEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	et, err := parseEventType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	st, err := parseEventStatus(req.Status)
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
	c.JSON(http.StatusOK, toEventInstanceJSON(e))
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	if err := h.deleteEvent.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) SetEventProjects(c *gin.Context) {
	var req setProjectsRequest
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

func (h *Handler) CreateModule(c *gin.Context) {
	var req createModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	mt, err := parseModuleType(req.Type)
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
	c.JSON(http.StatusCreated, toModuleResponse(m))
}

func (h *Handler) PatchModule(c *gin.Context) {
	var req patchModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	mt, err := parseModuleType(req.Type)
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
	c.JSON(http.StatusOK, toModuleResponse(m))
}

func (h *Handler) DeleteModule(c *gin.Context) {
	if err := h.deleteModule.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteOptionGroup(c *gin.Context) {
	if err := h.deleteOptionGroup.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteOption(c *gin.Context) {
	if err := h.deleteOption.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) SetModuleProjects(c *gin.Context) {
	var req setModuleProjectsRequest
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

func (h *Handler) CreateOptionGroup(c *gin.Context) {
	var req createOptionGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	ot, err := parseOptionGroupType(req.Type)
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
	c.JSON(http.StatusCreated, toOptionGroupResponse(g))
}

func (h *Handler) PatchOptionGroup(c *gin.Context) {
	var req patchOptionGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	ot, err := parseOptionGroupType(req.Type)
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
	c.JSON(http.StatusOK, toOptionGroupResponse(g))
}

func (h *Handler) CreateOption(c *gin.Context) {
	var req createOptionRequest
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
	c.JSON(http.StatusCreated, toOptionResponse(o))
}

func (h *Handler) PatchOption(c *gin.Context) {
	var req patchOptionRequest
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
	c.JSON(http.StatusOK, toOptionResponse(o))
}

func (h *Handler) SubmitResponse(c *gin.Context) {
	var req submitResponseRequest
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
	c.JSON(http.StatusOK, toMyResponseResult(resp, answers))
}

func (h *Handler) ListEventResponsesForAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	params := pagination.NewParams(page, pageSize)
	result, err := h.listEventResponses.Execute(c.Request.Context(), c.Param("id"), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toAdminResponseListResponse(result))
}

func (h *Handler) GetModuleResponseSummary(c *gin.Context) {
	summary, err := h.getModuleResponseSummary.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toModuleResponseSummaryJSON(summary))
}
