package expensesusecases

import (
	"context"
	"time"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	userdomain "macabi-back/internal/user/domain"
)

type ProjectExpenseSummaryUC struct {
	repo     expensesports.ExpenseRepository
	projects expensesports.ProjectMembership
}

func NewProjectExpenseSummaryUC(repo expensesports.ExpenseRepository, projects expensesports.ProjectMembership) *ProjectExpenseSummaryUC {
	return &ProjectExpenseSummaryUC{repo: repo, projects: projects}
}

type SummaryInput struct {
	ProjectID string
	UserID    string
	UserRole  string
	From      *time.Time
	To        *time.Time
}

func (uc *ProjectExpenseSummaryUC) Execute(ctx context.Context, in SummaryInput) (*expensesports.ProjectExpenseSummary, error) {
	if in.ProjectID == "" {
		return nil, expensesdomain.ErrMissingRequiredField
	}

	isAdmin := userdomain.Role(in.UserRole) == userdomain.RoleAdmin
	if !isAdmin {
		ok, err := uc.projects.IsProjectMember(ctx, in.ProjectID, in.UserID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, expensesdomain.ErrForbidden
		}
	}

	coord := false
	if !isAdmin {
		var err error
		coord, err = uc.projects.IsProjectCoordinator(ctx, in.ProjectID, in.UserID)
		if err != nil {
			return nil, err
		}
	}

	var filter *string
	if !isAdmin && !coord {
		u := in.UserID
		filter = &u
	}
	return uc.repo.SummaryByProject(ctx, in.ProjectID, filter, in.From, in.To)
}
