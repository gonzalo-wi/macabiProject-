package expensesusecases

import (
	"context"
	"strings"
	"time"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

// SubmitterNotifier notifies the expense submitter of approval or rejection.
type SubmitterNotifier interface {
	NotifySubmitterApproved(ctx context.Context, exp *expensesdomain.Expense)
	NotifySubmitterRejected(ctx context.Context, exp *expensesdomain.Expense, reason string)
}

type ApproveExpense struct {
	repo     expensesports.ExpenseRepository
	projects projectports.ProjectMembership
	notifier SubmitterNotifier
}

func NewApproveExpense(
	repo expensesports.ExpenseRepository,
	projects projectports.ProjectMembership,
	notifier SubmitterNotifier,
) *ApproveExpense {
	return &ApproveExpense{repo: repo, projects: projects, notifier: notifier}
}

type MutationInput struct {
	ExpenseID string
	ActorID   string
	UserRole  string
}

func (uc *ApproveExpense) Execute(ctx context.Context, in MutationInput) error {
	exp, err := uc.repo.FindByID(ctx, in.ExpenseID)
	if err != nil {
		return err
	}
	if exp.Status != expensesdomain.StatusPending {
		return expensesdomain.ErrInvalidStatusTransition
	}
	if ok, err := canApproveOrRejectExpense(ctx, uc.projects, exp.ProjectID, in.ActorID, in.UserRole); err != nil {
		return err
	} else if !ok {
		return expensesdomain.ErrForbidden
	}
	now := time.Now().UTC()
	u := in.ActorID
	exp.Status = expensesdomain.StatusApproved
	exp.ApprovedByUserID = &u
	exp.ApprovedAt = &now
	exp.RejectionReason = nil
	if err := uc.repo.Update(ctx, exp); err != nil {
		return err
	}

	uc.notifier.NotifySubmitterApproved(ctx, exp)
	return nil
}

type RejectExpense struct {
	repo     expensesports.ExpenseRepository
	projects projectports.ProjectMembership
	notifier SubmitterNotifier
}

func NewRejectExpense(
	repo expensesports.ExpenseRepository,
	projects projectports.ProjectMembership,
	notifier SubmitterNotifier,
) *RejectExpense {
	return &RejectExpense{repo: repo, projects: projects, notifier: notifier}
}

type RejectExpenseInput struct {
	MutationInput
	Reason string
}

func (uc *RejectExpense) Execute(ctx context.Context, in RejectExpenseInput) error {
	exp, err := uc.repo.FindByID(ctx, in.ExpenseID)
	if err != nil {
		return err
	}
	if exp.Status != expensesdomain.StatusPending {
		return expensesdomain.ErrInvalidStatusTransition
	}
	if ok, err := canApproveOrRejectExpense(ctx, uc.projects, exp.ProjectID, in.ActorID, in.UserRole); err != nil {
		return err
	} else if !ok {
		return expensesdomain.ErrForbidden
	}
	r := strings.TrimSpace(in.Reason)
	exp.Status = expensesdomain.StatusRejected
	exp.RejectionReason = &r
	exp.ApprovedByUserID = nil
	exp.ApprovedAt = nil
	if err := uc.repo.Update(ctx, exp); err != nil {
		return err
	}

	uc.notifier.NotifySubmitterRejected(ctx, exp, r)
	return nil
}

func canApproveOrRejectExpense(ctx context.Context, projects projectports.ProjectMembership, projectID, actorID, role string) (bool, error) {
	if userdomain.Role(role) == userdomain.RoleAdmin {
		return true, nil
	}
	return projects.IsProjectCoordinator(ctx, projectID, actorID)
}
