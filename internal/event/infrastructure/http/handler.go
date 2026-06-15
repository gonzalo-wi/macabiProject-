package eventhttp

import (
	eventusecases "macabi-back/internal/event/application/usecases"
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
	listNotifs               *eventusecases.ListEventNotifications
	markNotifRead            *eventusecases.MarkEventNotificationRead
	markAllNotifsRead        *eventusecases.MarkAllEventNotificationsRead
	unreadNotifs             *eventusecases.UnreadEventNotificationCount
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
	listNotifs *eventusecases.ListEventNotifications,
	markNotifRead *eventusecases.MarkEventNotificationRead,
	markAllNotifsRead *eventusecases.MarkAllEventNotificationsRead,
	unreadNotifs *eventusecases.UnreadEventNotificationCount,
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
		listNotifs:               listNotifs,
		markNotifRead:            markNotifRead,
		markAllNotifsRead:        markAllNotifsRead,
		unreadNotifs:             unreadNotifs,
	}
}
