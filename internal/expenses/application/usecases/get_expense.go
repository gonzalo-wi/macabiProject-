package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	userdomain "macabi-back/internal/user/domain"
)

type GetExpense struct {
	repo     expensesports.ExpenseRepository
	projects expensesports.ProjectMembership
}

func NewGetExpense(repo expensesports.ExpenseRepository, projects expensesports.ProjectMembership) *GetExpense {
	return &GetExpense{repo: repo, projects: projects}
}

type GetExpenseInput struct {
	ID       string
	UserID   string
	UserRole string
}

func (uc *GetExpense) Execute(ctx context.Context, input GetExpenseInput) (*expensesdomain.Expense, error) {
	exp, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if userdomain.Role(input.UserRole) == userdomain.RoleAdmin {
		return exp, nil
	}
	if exp.SubmittedByUserID == input.UserID {
		return exp, nil
	}
	coord, err := uc.projects.IsProjectCoordinator(ctx, exp.ProjectID, input.UserID)
	if err != nil {
		return nil, err
	}
	if coord {
		return exp, nil
	}

	return nil, expensesdomain.ErrForbidden
}
