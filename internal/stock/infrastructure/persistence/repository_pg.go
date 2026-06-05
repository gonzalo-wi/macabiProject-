package stockpersistence

import (
	"gorm.io/gorm"

	stockports "macabi-back/internal/stock/application/ports"
)

type StockRepositoryPG struct {
	db *gorm.DB
}

func NewStockRepositoryPG(db *gorm.DB) *StockRepositoryPG {
	return &StockRepositoryPG{db: db}
}

var _ stockports.StockRepository = (*StockRepositoryPG)(nil)
var _ stockports.ProjectMemberReader = (*StockRepositoryPG)(nil)
var _ stockports.UserEmailReader = (*StockRepositoryPG)(nil)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&ResourceModel{},
		&RequestModel{},
		&NotificationModel{},
		&PushSubscriptionModel{},
	)
}
