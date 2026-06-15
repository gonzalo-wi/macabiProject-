package expensesusecases

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

type GetProjectBudget struct {
	repo     expensesports.ExpenseQueryRepository
	projects projectports.ProjectMembership
}

func NewGetProjectBudget(repo expensesports.ExpenseQueryRepository, projects projectports.ProjectMembership) *GetProjectBudget {
	return &GetProjectBudget{repo: repo, projects: projects}
}

func (uc *GetProjectBudget) Execute(ctx context.Context, projectID, userID, userRole string) (*expensesports.ProjectBudgetStatus, error) {
	if projectID == "" {
		return nil, expensesdomain.ErrMissingRequiredField
	}
	if userdomain.Role(userRole) != userdomain.RoleAdmin {
		ok, err := uc.projects.IsProjectMember(ctx, projectID, userID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, expensesdomain.ErrForbidden
		}
	}

	budget, err := uc.repo.GetBudget(ctx, projectID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	sum, err := uc.repo.SummaryByProject(ctx, projectID, nil, &firstOfMonth, &now)
	if err != nil {
		return nil, err
	}

	return &expensesports.ProjectBudgetStatus{
		MonthlyAmount:        budget,
		CurrentMonthApproved: sum.TotalApproved,
		Month:                now.Format("2006-01"),
	}, nil
}

type SetProjectBudget struct {
	repo expensesports.ExpenseQueryRepository
}

func NewSetProjectBudget(repo expensesports.ExpenseQueryRepository) *SetProjectBudget {
	return &SetProjectBudget{repo: repo}
}

func (uc *SetProjectBudget) Execute(ctx context.Context, projectID string, amount *decimal.Decimal) error {
	if projectID == "" {
		return expensesdomain.ErrMissingRequiredField
	}
	if amount != nil && amount.IsNegative() {
		return expensesdomain.ErrInvalidAmount
	}
	return uc.repo.SetBudget(ctx, projectID, amount)
}
