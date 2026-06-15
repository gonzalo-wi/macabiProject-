package eventpersistence

import (
	"gorm.io/gorm"

	eventports "macabi-back/internal/event/application/ports"
)

type RepositoryPG struct {
	db *gorm.DB
}

func NewRepositoryPG(db *gorm.DB) *RepositoryPG {
	return &RepositoryPG{db: db}
}

var _ eventports.EventRepository = (*RepositoryPG)(nil)
var _ eventports.ModuleRepository = (*RepositoryPG)(nil)
var _ eventports.OptionRepository = (*RepositoryPG)(nil)
var _ eventports.ResponseRepository = (*RepositoryPG)(nil)
var _ eventports.ReminderRepository = (*RepositoryPG)(nil)
var _ eventports.EventNotificationRepository = (*RepositoryPG)(nil)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&EventNotificationModel{})
}
