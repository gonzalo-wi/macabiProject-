package expensesports

import "context"

// ExpenseMailer sends email notifications related to project expenses.
type ExpenseMailer interface {
	NotifyCoordinatorsNewExpense(ctx context.Context, coordinatorEmails []string, amount, description, projectName string) error
	NotifySubmitterApproved(ctx context.Context, submitterEmail, amount, description string) error
	NotifySubmitterRejected(ctx context.Context, submitterEmail, amount, description, reason string) error
}
