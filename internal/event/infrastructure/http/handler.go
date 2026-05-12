package eventhttp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	eventusecases "macabi-back/internal/event/application/usecases"
	eventdomain "macabi-back/internal/event/domain"
	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

type Handler struct {
	svc *eventusecases.Service
}

func NewHandler(svc *eventusecases.Service) *Handler {
	return &Handler{svc: svc}
}

func parseEventType(s string) (eventdomain.EventType, error) {
	switch s {
	case "", string(eventdomain.EventTypeActivity):
		return eventdomain.EventTypeActivity, nil
	case string(eventdomain.EventTypeCustom):
		return eventdomain.EventTypeCustom, nil
	default:
		return "", eventdomain.ErrInvalidEventType
	}
}

func parseEventStatus(s string) (eventdomain.EventStatus, error) {
	switch s {
	case "", string(eventdomain.EventStatusDraft):
		return eventdomain.EventStatusDraft, nil
	case string(eventdomain.EventStatusOpen):
		return eventdomain.EventStatusOpen, nil
	case string(eventdomain.EventStatusClosed):
		return eventdomain.EventStatusClosed, nil
	case string(eventdomain.EventStatusCancelled):
		return eventdomain.EventStatusCancelled, nil
	default:
		return "", eventdomain.ErrInvalidEventStatus
	}
}

func parseModuleType(s string) (eventdomain.ModuleType, error) {
	switch eventdomain.ModuleType(s) {
	case eventdomain.ModuleAttendance, eventdomain.ModuleMeal, eventdomain.ModuleTransport,
		eventdomain.ModuleMaterials, eventdomain.ModuleCustom:
		return eventdomain.ModuleType(s), nil
	default:
		return "", eventdomain.ErrInvalidModuleType
	}
}

func parseOptionGroupType(s string) (eventdomain.OptionGroupType, error) {
	switch eventdomain.OptionGroupType(s) {
	case eventdomain.OptionGroupSingleChoice, eventdomain.OptionGroupMultipleChoice,
		eventdomain.OptionGroupText, eventdomain.OptionGroupNumber:
		return eventdomain.OptionGroupType(s), nil
	default:
		return "", eventdomain.ErrInvalidOptionGroupType
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
	e, err := h.svc.CreateEvent(c.Request.Context(), req.Title, et, req.StartsAt, req.ResponseDeadlineAt, st)
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
	res, err := h.svc.ListEvents(c.Request.Context(), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, eventListResponse(res))
}

func (h *Handler) GetEventDetail(c *gin.Context) {
	d, err := h.svc.GetEventDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, detailToJSON(d))
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
	e, err := h.svc.UpdateEvent(c.Request.Context(), c.Param("id"), req.Title, et, req.StartsAt, req.ResponseDeadlineAt, st)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toEventInstanceJSON(e))
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	if err := h.svc.DeleteEvent(c.Request.Context(), c.Param("id")); err != nil {
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
	if err := h.svc.SetEventProjects(c.Request.Context(), c.Param("id"), req.ProjectIDs); err != nil {
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
	m, err := h.svc.CreateModule(c.Request.Context(), req.EventInstanceID, req.Title, mt, req.SortOrder, req.IsRequired)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, moduleToJSON(m))
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
	m, err := h.svc.UpdateModule(c.Request.Context(), c.Param("id"), req.Title, mt, req.SortOrder, req.IsRequired)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, moduleToJSON(m))
}

func (h *Handler) DeleteModule(c *gin.Context) {
	if err := h.svc.DeleteModule(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteOptionGroup(c *gin.Context) {
	if err := h.svc.DeleteOptionGroup(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteOption(c *gin.Context) {
	if err := h.svc.DeleteOption(c.Request.Context(), c.Param("id")); err != nil {
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
	if err := h.svc.SetModuleProjects(c.Request.Context(), c.Param("id"), req.ProjectIDs); err != nil {
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
	g, err := h.svc.CreateOptionGroup(c.Request.Context(), req.ModuleID, req.Name, ot, req.SortOrder, req.IsRequired)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, groupToJSON(g))
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
	g, err := h.svc.UpdateOptionGroup(c.Request.Context(), c.Param("id"), req.Name, ot, req.SortOrder, req.IsRequired)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, groupToJSON(g))
}

func (h *Handler) CreateOption(c *gin.Context) {
	var req createOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	o, err := h.svc.CreateOption(c.Request.Context(), req.GroupID, req.Label, req.MaxCapacity, req.SortOrder)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, optionToJSON(o))
}

func (h *Handler) PatchOption(c *gin.Context) {
	var req patchOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	o, err := h.svc.UpdateOption(c.Request.Context(), c.Param("id"), req.Label, req.MaxCapacity, req.SortOrder, req.CurrentCount)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, optionToJSON(o))
}

func (h *Handler) SubmitResponse(c *gin.Context) {
	var req submitResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	userID := c.GetString(userhttp.AuthUserIDKey)
	if err := h.svc.SubmitResponse(c.Request.Context(), c.Param("id"), userID, req.ProjectID, req.Answers); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetMyResponse(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	resp, answers, err := h.svc.GetMyResponse(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if resp == nil {
		c.JSON(http.StatusOK, gin.H{"response": nil, "answers": []any{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": responseToJSON(resp),
		"answers":  answersToJSON(answers),
	})
}

func (h *Handler) ListEventResponsesForAdmin(c *gin.Context) {
	list, err := h.svc.ListResponsesForAdmin(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	data := make([]gin.H, len(list))
	for i := range list {
		row := list[i]
		rmeta := responseToJSON(&row.Response)
		rmeta["user_name"] = row.UserName
		rmeta["user_email"] = row.UserEmail
		data[i] = gin.H{
			"response": rmeta,
			"answers":  answersToJSON(row.Answers),
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func moduleToJSON(m *eventdomain.EventModule) gin.H {
	return gin.H{
		"id":                m.ID,
		"event_instance_id": m.EventInstanceID,
		"title":             m.Title,
		"type":              string(m.Type),
		"sort_order":        m.SortOrder,
		"is_required":       m.IsRequired,
		"created_at":        m.CreatedAt,
	}
}

func groupToJSON(g *eventdomain.EventOptionGroup) gin.H {
	return gin.H{
		"id":          g.ID,
		"module_id":   g.ModuleID,
		"name":        g.Name,
		"type":        string(g.Type),
		"sort_order":  g.SortOrder,
		"is_required": g.IsRequired,
	}
}

func optionToJSON(o *eventdomain.EventOption) gin.H {
	return gin.H{
		"id":            o.ID,
		"group_id":      o.GroupID,
		"label":         o.Label,
		"max_capacity":  o.MaxCapacity,
		"current_count": o.CurrentCount,
		"sort_order":    o.SortOrder,
	}
}

func responseToJSON(r *eventdomain.EventResponse) gin.H {
	h := gin.H{
		"id":                r.ID,
		"event_instance_id": r.EventInstanceID,
		"user_id":           r.UserID,
		"created_at":        r.CreatedAt,
	}
	if r.ProjectID != nil {
		h["project_id"] = *r.ProjectID
	}
	return h
}

func answersToJSON(aa []eventdomain.EventResponseAnswer) []gin.H {
	out := make([]gin.H, len(aa))
	for i, a := range aa {
		h := gin.H{
			"id":          a.ID,
			"response_id": a.ResponseID,
		}
		if a.GroupID != nil {
			h["group_id"] = *a.GroupID
		}
		if a.OptionID != nil {
			h["option_id"] = *a.OptionID
		}
		if a.TextValue != nil {
			h["text_value"] = *a.TextValue
		}
		out[i] = h
	}
	return out
}

func detailToJSON(d *eventdomain.EventDetail) gin.H {
	mods := make([]gin.H, len(d.Modules))
	for i, md := range d.Modules {
		groups := make([]gin.H, len(md.OptionGroups))
		for j, gd := range md.OptionGroups {
			opts := make([]gin.H, len(gd.Options))
			for k, o := range gd.Options {
				opts[k] = optionToJSON(&o)
			}
			groups[j] = gin.H{
				"group":   groupToJSON(&gd.Group),
				"options": opts,
			}
		}
		mods[i] = gin.H{
			"module":        moduleToJSON(&md.Module),
			"project_ids":   md.ProjectIDs,
			"option_groups": groups,
		}
	}
	return gin.H{
		"instance":    toEventInstanceJSON(&d.Instance),
		"project_ids": d.ProjectIDs,
		"modules":     mods,
	}
}
