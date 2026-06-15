package expensesdto

import (
	expensesports "macabi-back/internal/expenses/application/ports"
)

type ProjectBudgetResponse struct {
	MonthlyBudget        *string `json:"monthly_budget"`
	CurrentMonthApproved string  `json:"current_month_approved"`
	Month                string  `json:"month"`
}

type SetProjectBudgetRequest struct {
	MonthlyAmount *string `json:"monthly_amount"`
}

func BudgetToResp(s *expensesports.ProjectBudgetStatus) ProjectBudgetResponse {
	var budget *string
	if s.MonthlyAmount != nil {
		v := s.MonthlyAmount.String()
		budget = &v
	}
	return ProjectBudgetResponse{
		MonthlyBudget:        budget,
		CurrentMonthApproved: s.CurrentMonthApproved.String(),
		Month:                s.Month,
	}
}
