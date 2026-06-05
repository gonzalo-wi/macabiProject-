package expensesdto

import (
	"time"

	expensesdomain "macabi-back/internal/expenses/domain"
)

type CreateExpenseCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type ExpenseCategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ToExpenseCategoryResponse(c expensesdomain.ExpenseCategory) ExpenseCategoryResponse {
	return ExpenseCategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}

func ToExpenseCategoryListResponse(cats []expensesdomain.ExpenseCategory) []ExpenseCategoryResponse {
	out := make([]ExpenseCategoryResponse, len(cats))
	for i, c := range cats {
		out[i] = ToExpenseCategoryResponse(c)
	}
	return out
}
