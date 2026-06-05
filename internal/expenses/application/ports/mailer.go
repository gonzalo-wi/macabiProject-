package expensesports

import "context"

type ExpenseMailer interface {
	NotifyCoordinatorsNewExpense(ctx context.Context, coordinatorEmails []string, amount, description, projectName, expenseID string) error
	NotifySubmitterApproved(ctx context.Context, submitterEmail, amount, description, expenseID string) error
	NotifySubmitterRejected(ctx context.Context, submitterEmail, amount, description, reason, expenseID string) error
}
