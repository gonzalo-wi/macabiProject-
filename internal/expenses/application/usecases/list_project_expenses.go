package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	"macabi-back/internal/shared/pagination"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type ListProjectExpenses struct {
	repo     expensesports.ExpenseQueryRepository
	projects projectports.ProjectMembership
}

func NewListProjectExpenses(repo expensesports.ExpenseQueryRepository, projects projectports.ProjectMembership) *ListProjectExpenses {
	return &ListProjectExpenses{repo: repo, projects: projects}
}

type ListProjectExpensesInput struct {
	ProjectID string
	UserID    string
	UserRole  string
	Params    pagination.Params
}

func (uc *ListProjectExpenses) Execute(ctx context.Context, input ListProjectExpensesInput) (pagination.Result[expensesdomain.ExpenseListItem], error) {
	if input.ProjectID == "" {
		return pagination.Result[expensesdomain.ExpenseListItem]{}, expensesdomain.ErrMissingRequiredField
	}

	isAdmin := userdomain.Role(input.UserRole) == userdomain.RoleAdmin
	if !isAdmin {
		ok, err := uc.projects.IsProjectMember(ctx, input.ProjectID, input.UserID)
		if err != nil {
			return pagination.Result[expensesdomain.ExpenseListItem]{}, err
		}
		if !ok {
			return pagination.Result[expensesdomain.ExpenseListItem]{}, expensesdomain.ErrForbidden
		}
	}

	coord := false
	if !isAdmin {
		var err error
		coord, err = uc.projects.IsProjectCoordinator(ctx, input.ProjectID, input.UserID)
		if err != nil {
			return pagination.Result[expensesdomain.ExpenseListItem]{}, err
		}
	}

	var filter *string
	if !isAdmin && !coord {
		u := input.UserID
		filter = &u
	}
	return uc.repo.ListByProject(ctx, input.ProjectID, filter, input.Params)
}
