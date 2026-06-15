package expensespersistence

import (
	"gorm.io/gorm"

	expensesports "macabi-back/internal/expenses/application/ports"
)

type ExpenseRepositoryPG struct {
	db *gorm.DB
}

var _ expensesports.ExpenseRepository = (*ExpenseRepositoryPG)(nil)
var _ expensesports.ExpenseQueryRepository = (*ExpenseRepositoryPG)(nil)
var _ expensesports.ExpenseNotificationRepository = (*ExpenseRepositoryPG)(nil)
var _ expensesports.CategoryRepository = (*ExpenseRepositoryPG)(nil)

func NewExpenseRepositoryPG(db *gorm.DB) *ExpenseRepositoryPG {
	return &ExpenseRepositoryPG{db: db}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&ExpenseCategoryModel{},
		&ExpenseModel{},
		&ExpenseNotificationModel{},
		&ProjectExpenseBudgetModel{},
	)
}
