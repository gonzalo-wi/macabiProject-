package deps

import (
	eventusecases "macabi-back/internal/event/application/usecases"
	eventhttp "macabi-back/internal/event/infrastructure/http"
	eventmail "macabi-back/internal/event/infrastructure/mail"
	eventpersistence "macabi-back/internal/event/infrastructure/persistence"
	"macabi-back/internal/shared/config"

	"gorm.io/gorm"
)

func buildEventDeps(db *gorm.DB, cfg *config.Config) (*eventhttp.Handler, *eventusecases.SendEventReminders) {
	eventRepo := eventpersistence.NewRepositoryPG(db)
	eventMailer := eventmail.NewBrevoEventMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom, cfg.FrontendPublicURL)
	sendEventRemindersUC := eventusecases.NewSendEventReminders(eventRepo, eventMailer)

	eventHandler := eventhttp.NewHandler(
		eventusecases.NewCreateEvent(eventRepo),
		eventusecases.NewListEvents(eventRepo),
		eventusecases.NewGetEventDetail(eventRepo),
		eventusecases.NewUpdateEvent(eventRepo),
		eventusecases.NewDeleteEvent(eventRepo),
		eventusecases.NewSetEventProjects(eventRepo),
		eventusecases.NewCreateModule(eventRepo),
		eventusecases.NewUpdateModule(eventRepo),
		eventusecases.NewDeleteModule(eventRepo),
		eventusecases.NewSetModuleProjects(eventRepo),
		eventusecases.NewCreateOptionGroup(eventRepo),
		eventusecases.NewUpdateOptionGroup(eventRepo),
		eventusecases.NewDeleteOptionGroup(eventRepo),
		eventusecases.NewCreateOption(eventRepo),
		eventusecases.NewUpdateOption(eventRepo),
		eventusecases.NewDeleteOption(eventRepo),
		eventusecases.NewSubmitResponse(eventRepo, eventRepo),
		eventusecases.NewGetMyResponse(eventRepo),
		eventusecases.NewListEventResponses(eventRepo, eventRepo),
		eventusecases.NewGetModuleResponseSummary(eventRepo),
	)

	return eventHandler, sendEventRemindersUC
}
