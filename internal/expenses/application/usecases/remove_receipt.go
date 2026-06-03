package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	userdomain "macabi-back/internal/user/domain"
)

type RemoveReceipt struct {
	repo     expensesports.ExpenseRepository
	projects expensesports.ProjectMembership
	signer   expensesports.ReceiptSigner
}

func NewRemoveReceipt(
	repo expensesports.ExpenseRepository,
	projects expensesports.ProjectMembership,
	signer expensesports.ReceiptSigner,
) *RemoveReceipt {
	return &RemoveReceipt{repo: repo, projects: projects, signer: signer}
}

type RemoveReceiptInput struct {
	ExpenseID string
	ActorID   string
	UserRole  string
}

func (uc *RemoveReceipt) Execute(ctx context.Context, in RemoveReceiptInput) error {
	exp, err := uc.repo.FindByID(ctx, in.ExpenseID)
	if err != nil {
		return err
	}

	isAdmin := userdomain.Role(in.UserRole) == userdomain.RoleAdmin
	coord, err := uc.projects.IsProjectCoordinator(ctx, exp.ProjectID, in.ActorID)
	if err != nil {
		return err
	}
	owner := exp.SubmittedByUserID == in.ActorID
	// Mismo criterio que editar el comprobante en UpdateExpense.
	if !(isAdmin || coord || (owner && exp.Status == expensesdomain.StatusPending)) {
		return expensesdomain.ErrForbidden
	}

	if exp.ReceiptStoragePath != nil && *exp.ReceiptStoragePath != "" && uc.signer != nil {
		_ = uc.signer.DeleteObject(ctx, *exp.ReceiptStoragePath)
	}
	exp.ReceiptStoragePath = nil
	return uc.repo.Update(ctx, exp)
}
