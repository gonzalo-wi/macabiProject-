package expensesports

import (
	"context"

	expensesdomain "macabi-back/internal/expenses/domain"
)

type CategoryRepository interface {
	SaveCategory(ctx context.Context, c *expensesdomain.ExpenseCategory) error
	ListCategories(ctx context.Context) ([]expensesdomain.ExpenseCategory, error)
	FindCategoryByID(ctx context.Context, id string) (*expensesdomain.ExpenseCategory, error)
	DeleteCategory(ctx context.Context, id string) error
}
