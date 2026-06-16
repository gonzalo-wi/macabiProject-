package stockhttp

import (
	stockusecases "macabi-back/internal/stock/application/usecases"
)

type Handler struct {
	createResourceUC           *stockusecases.CreateResource
	listResourcesUC            *stockusecases.ListResources
	getResourceUC              *stockusecases.GetResource
	updateResourceUC           *stockusecases.UpdateResource
	deleteResourceUC           *stockusecases.DeleteResource
	createRequestUC            *stockusecases.CreateRequest
	approveRequestUC           *stockusecases.ApproveRequest
	rejectRequestUC            *stockusecases.RejectRequest
	cancelRequestUC            *stockusecases.CancelRequest
	deliverRequestUC           *stockusecases.DeliverRequest
	returnRequestUC            *stockusecases.ReturnRequest
	listRequestsUC             *stockusecases.ListRequests
	listMyRequestsUC           *stockusecases.ListMyRequests
	getRequestDetailUC         *stockusecases.GetRequestDetail
	listNotificationsUC        *stockusecases.ListNotifications
	markNotificationReadUC     *stockusecases.MarkNotificationRead
	markAllNotificationsReadUC *stockusecases.MarkAllNotificationsRead
	unreadCountUC              *stockusecases.UnreadNotificationCount
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
	markAllNotificationsReadUC *stockusecases.MarkAllNotificationsRead,
	unreadCountUC *stockusecases.UnreadNotificationCount,
) *Handler {
	return &Handler{
		createResourceUC:           createResourceUC,
		listResourcesUC:            listResourcesUC,
		getResourceUC:              getResourceUC,
		updateResourceUC:           updateResourceUC,
		deleteResourceUC:           deleteResourceUC,
		createRequestUC:            createRequestUC,
		approveRequestUC:           approveRequestUC,
		rejectRequestUC:            rejectRequestUC,
		cancelRequestUC:            cancelRequestUC,
		deliverRequestUC:           deliverRequestUC,
		returnRequestUC:            returnRequestUC,
		listRequestsUC:             listRequestsUC,
		listMyRequestsUC:           listMyRequestsUC,
		getRequestDetailUC:         getRequestDetailUC,
		listNotificationsUC:        listNotificationsUC,
		markNotificationReadUC:     markNotificationReadUC,
		markAllNotificationsReadUC: markAllNotificationsReadUC,
		unreadCountUC:              unreadCountUC,
	}
}
