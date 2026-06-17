package stockpersistence

import (
	"gorm.io/gorm"

	stockports "macabi-back/internal/stock/application/ports"
	projectports "macabi-back/internal/project/application/ports"
	userports "macabi-back/internal/user/application/ports"
)

type StockRepositoryPG struct {
	db *gorm.DB
}

func NewStockRepositoryPG(db *gorm.DB) *StockRepositoryPG {
	return &StockRepositoryPG{db: db}
}

var _ stockports.ResourceRepository = (*StockRepositoryPG)(nil)
var _ stockports.RequestRepository = (*StockRepositoryPG)(nil)
var _ stockports.NotificationRepository = (*StockRepositoryPG)(nil)
var _ projectports.ProjectMemberReader = (*StockRepositoryPG)(nil)
var _ userports.UserEmailReader = (*StockRepositoryPG)(nil)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&ResourceModel{},
		&RequestModel{},
		&NotificationModel{},
	)
}
