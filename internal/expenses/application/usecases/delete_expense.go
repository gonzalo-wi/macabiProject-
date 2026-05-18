package expensesusecases

import (
	"context"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	userdomain "macabi-back/internal/user/domain"
)

type DeleteExpense struct {
	repo     expensesports.ExpenseRepository
	projects expensesports.ProjectMembership
	signer   expensesports.ReceiptSigner
}

func NewDeleteExpense(
	repo expensesports.ExpenseRepository,
	projects expensesports.ProjectMembership,
	signer expensesports.ReceiptSigner,
) *DeleteExpense {
	return &DeleteExpense{repo: repo, projects: projects, signer: signer}
}

type DeleteExpenseInput struct {
	ExpenseID string
	ActorID   string
	UserRole  string
}

func (uc *DeleteExpense) Execute(ctx context.Context, in DeleteExpenseInput) error {
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
	allowed := false
	if owner && exp.Status == expensesdomain.StatusPending {
		allowed = true
	}
	if (isAdmin || coord) && (exp.Status == expensesdomain.StatusApproved || exp.Status == expensesdomain.StatusRejected) {
		allowed = true
	}
	if !allowed {
		return expensesdomain.ErrForbidden
	}

	if exp.ReceiptStoragePath != nil && *exp.ReceiptStoragePath != "" && uc.signer != nil {
		_ = uc.signer.DeleteObject(ctx, *exp.ReceiptStoragePath)
	}
	_ = uc.repo.DeleteNotificationsByExpenseID(ctx, exp.ID)
	return uc.repo.DeleteByID(ctx, exp.ID)
}
