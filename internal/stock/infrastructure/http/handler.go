package stockhttp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	sharederrors "macabi-back/internal/shared/errors"
	"macabi-back/internal/shared/pagination"
	stockusecases "macabi-back/internal/stock/application/usecases"
	userhttp "macabi-back/internal/user/infrastructure/http"
)

type Handler struct {
	createResourceUC       *stockusecases.CreateResource
	listResourcesUC        *stockusecases.ListResources
	getResourceUC          *stockusecases.GetResource
	updateResourceUC       *stockusecases.UpdateResource
	deleteResourceUC       *stockusecases.DeleteResource
	createRequestUC        *stockusecases.CreateRequest
	approveRequestUC       *stockusecases.ApproveRequest
	rejectRequestUC        *stockusecases.RejectRequest
	cancelRequestUC        *stockusecases.CancelRequest
	deliverRequestUC       *stockusecases.DeliverRequest
	returnRequestUC        *stockusecases.ReturnRequest
	listRequestsUC         *stockusecases.ListRequests
	listMyRequestsUC       *stockusecases.ListMyRequests
	getRequestDetailUC     *stockusecases.GetRequestDetail
	listNotificationsUC    *stockusecases.ListNotifications
	markNotificationReadUC *stockusecases.MarkNotificationRead
	unreadCountUC          *stockusecases.UnreadNotificationCount
	registerPushSubUC      *stockusecases.RegisterPushSubscription
	unregisterPushSubUC    *stockusecases.UnregisterPushSubscription
	vapidPublicKey         string
}

func NewHandler(
	createResourceUC *stockusecases.CreateResource,
	listResourcesUC *stockusecases.ListResources,
	getResourceUC *stockusecases.GetResource,
	updateResourceUC *stockusecases.UpdateResource,
	deleteResourceUC *stockusecases.DeleteResource,
	createRequestUC *stockusecases.CreateRequest,
	approveRequestUC *stockusecases.ApproveRequest,
	rejectRequestUC *stockusecases.RejectRequest,
	cancelRequestUC *stockusecases.CancelRequest,
	deliverRequestUC *stockusecases.DeliverRequest,
	returnRequestUC *stockusecases.ReturnRequest,
	listRequestsUC *stockusecases.ListRequests,
	listMyRequestsUC *stockusecases.ListMyRequests,
	getRequestDetailUC *stockusecases.GetRequestDetail,
	listNotificationsUC *stockusecases.ListNotifications,
	markNotificationReadUC *stockusecases.MarkNotificationRead,
	unreadCountUC *stockusecases.UnreadNotificationCount,
	registerPushSubUC *stockusecases.RegisterPushSubscription,
	unregisterPushSubUC *stockusecases.UnregisterPushSubscription,
	vapidPublicKey string,
) *Handler {
	return &Handler{
		createResourceUC:       createResourceUC,
		listResourcesUC:        listResourcesUC,
		getResourceUC:          getResourceUC,
		updateResourceUC:       updateResourceUC,
		deleteResourceUC:       deleteResourceUC,
		createRequestUC:        createRequestUC,
		approveRequestUC:       approveRequestUC,
		rejectRequestUC:        rejectRequestUC,
		cancelRequestUC:        cancelRequestUC,
		deliverRequestUC:       deliverRequestUC,
		returnRequestUC:        returnRequestUC,
		listRequestsUC:         listRequestsUC,
		listMyRequestsUC:       listMyRequestsUC,
		getRequestDetailUC:     getRequestDetailUC,
		listNotificationsUC:    listNotificationsUC,
		markNotificationReadUC: markNotificationReadUC,
		unreadCountUC:          unreadCountUC,
		registerPushSubUC:      registerPushSubUC,
		unregisterPushSubUC:    unregisterPushSubUC,
		vapidPublicKey:         vapidPublicKey,
	}
}

func (h *Handler) CreateResource(c *gin.Context) {
	var req CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	r, err := h.createResourceUC.Execute(c.Request.Context(), req.toInput())
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, toResourceResponse(r))
}

func (h *Handler) ListResources(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	params := pagination.NewParams(page, pageSize)
	result, err := h.listResourcesUC.Execute(c.Request.Context(), params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toResourceListResponse(result))
}

func (h *Handler) GetResource(c *gin.Context) {
	r, err := h.getResourceUC.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toResourceResponse(r))
}

func (h *Handler) UpdateResource(c *gin.Context) {
	var req UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	r, err := h.updateResourceUC.Execute(c.Request.Context(), req.toInput(c.Param("id")))
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toResourceResponse(r))
}

func (h *Handler) DeleteResource(c *gin.Context) {
	if err := h.deleteResourceUC.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateRequest(c *gin.Context) {
	var req CreateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	userID := c.GetString(userhttp.AuthUserIDKey)
	userRole := c.GetString(userhttp.AuthRoleKey)
	input, err := req.toInput(userID, userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	request, err := h.createRequestUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, toRequestResponse(request))
}

func (h *Handler) ListRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	params := pagination.NewParams(page, pageSize)
	projectID := c.Query("project_id")
	userID := c.GetString(userhttp.AuthUserIDKey)
	userRole := c.GetString(userhttp.AuthRoleKey)
	result, err := h.listRequestsUC.Execute(c.Request.Context(), stockusecases.ListRequestsInput{
		Params:    params,
		ProjectID: projectID,
		UserID:    userID,
		UserRole:  userRole,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toRequestDetailListResponse(result))
}

func (h *Handler) ListMyRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	params := pagination.NewParams(page, pageSize)
	userID := c.GetString(userhttp.AuthUserIDKey)
	result, err := h.listMyRequestsUC.Execute(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toRequestDetailListResponse(result))
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
	c.JSON(http.StatusOK, toRequestDetailResponse(detail))
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	userRole := c.GetString(userhttp.AuthRoleKey)
	err := h.approveRequestUC.Execute(c.Request.Context(), stockusecases.ApproveRequestInput{
		RequestID: c.Param("id"),
		UserID:    userID,
		UserRole:  userRole,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) RejectRequest(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	userRole := c.GetString(userhttp.AuthRoleKey)
	err := h.rejectRequestUC.Execute(c.Request.Context(), stockusecases.RejectRequestInput{
		RequestID: c.Param("id"),
		UserID:    userID,
		UserRole:  userRole,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CancelRequest(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	userRole := c.GetString(userhttp.AuthRoleKey)
	err := h.cancelRequestUC.Execute(c.Request.Context(), stockusecases.CancelRequestInput{
		RequestID: c.Param("id"),
		UserID:    userID,
		UserRole:  userRole,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeliverRequest(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	userRole := c.GetString(userhttp.AuthRoleKey)
	err := h.deliverRequestUC.Execute(c.Request.Context(), stockusecases.DeliverRequestInput{
		RequestID: c.Param("id"),
		UserID:    userID,
		UserRole:  userRole,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ReturnRequest(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	userRole := c.GetString(userhttp.AuthRoleKey)
	err := h.returnRequestUC.Execute(c.Request.Context(), stockusecases.ReturnRequestInput{
		RequestID: c.Param("id"),
		UserID:    userID,
		UserRole:  userRole,
	})
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListNotifications(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	notifs, err := h.listNotificationsUC.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, toNotificationListResponse(notifs))
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	if err := h.markNotificationReadUC.Execute(c.Request.Context(), c.Param("id"), userID); err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UnreadCount(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	count, err := h.unreadCountUC.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(httpStatus(err), sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// --- Web Push handlers ---

func (h *Handler) GetVAPIDPublicKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"public_key": h.vapidPublicKey})
}

type registerPushSubscriptionRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
	P256DH   string `json:"p256dh"   binding:"required"`
	Auth     string `json:"auth"     binding:"required"`
}

func (h *Handler) RegisterPushSubscription(c *gin.Context) {
	userID := c.GetString(userhttp.AuthUserIDKey)
	var req registerPushSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	if err := h.registerPushSubUC.Execute(c.Request.Context(), stockusecases.RegisterPushSubscriptionInput{
		UserID:   userID,
		Endpoint: req.Endpoint,
		P256DH:   req.P256DH,
		Auth:     req.Auth,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, sharederrors.NewErrorResponse(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

type unregisterPushSubscriptionRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
}

func (h *Handler) UnregisterPushSubscription(c *gin.Context) {
	var req unregisterPushSubscriptionRequest
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
