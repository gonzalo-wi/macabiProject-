package expensesports

import "context"

// ExpenseMailer sends email notifications related to project expenses.
// expenseID is used to build a deep link to the expense detail.
type ExpenseMailer interface {
	NotifyCoordinatorsNewExpense(ctx context.Context, coordinatorEmails []string, amount, description, projectName, expenseID string) error
	NotifySubmitterApproved(ctx context.Context, submitterEmail, amount, description, expenseID string) error
	NotifySubmitterRejected(ctx context.Context, submitterEmail, amount, description, reason, expenseID string) error
}
