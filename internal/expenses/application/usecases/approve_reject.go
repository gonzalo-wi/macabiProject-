package expensesusecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	userdomain "macabi-back/internal/user/domain"
)

type ApproveExpense struct {
	repo     expensesports.ExpenseRepository
	notifs   expensesports.ExpenseNotificationRepository
	projects expensesports.ProjectMembership
	emails   expensesports.UserEmailReader
	mailer   expensesports.ExpenseMailer
}

func NewApproveExpense(
	repo expensesports.ExpenseRepository,
	notifs expensesports.ExpenseNotificationRepository,
	projects expensesports.ProjectMembership,
	emails expensesports.UserEmailReader,
	mailer expensesports.ExpenseMailer,
) *ApproveExpense {
	return &ApproveExpense{repo: repo, notifs: notifs, projects: projects, emails: emails, mailer: mailer}
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

	_ = uc.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
		UserID:    exp.SubmittedByUserID,
		ExpenseID: exp.ID,
		ProjectID: exp.ProjectID,
		Message:   fmt.Sprintf("Tu gasto fue aprobado: %s — %s", formatExpenseAmount(exp), exp.Description),
	})

	if email, err := uc.emails.FindEmailByID(ctx, exp.SubmittedByUserID); err == nil {
		_ = uc.mailer.NotifySubmitterApproved(ctx, email, formatExpenseAmount(exp), exp.Description, exp.ID)
	}
	return nil
}

type RejectExpense struct {
	repo     expensesports.ExpenseRepository
	notifs   expensesports.ExpenseNotificationRepository
	projects expensesports.ProjectMembership
	emails   expensesports.UserEmailReader
	mailer   expensesports.ExpenseMailer
}

func NewRejectExpense(
	repo expensesports.ExpenseRepository,
	notifs expensesports.ExpenseNotificationRepository,
	projects expensesports.ProjectMembership,
	emails expensesports.UserEmailReader,
	mailer expensesports.ExpenseMailer,
) *RejectExpense {
	return &RejectExpense{repo: repo, notifs: notifs, projects: projects, emails: emails, mailer: mailer}
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

	rejectedMsg := fmt.Sprintf("Tu gasto fue rechazado: %s — %s", formatExpenseAmount(exp), exp.Description)
	if r != "" {
		rejectedMsg += fmt.Sprintf(" (motivo: %s)", r)
	}
	_ = uc.notifs.SaveNotification(ctx, &expensesdomain.ExpenseNotification{
		UserID:    exp.SubmittedByUserID,
		ExpenseID: exp.ID,
		ProjectID: exp.ProjectID,
		Message:   rejectedMsg,
	})

	if email, err := uc.emails.FindEmailByID(ctx, exp.SubmittedByUserID); err == nil {
		_ = uc.mailer.NotifySubmitterRejected(ctx, email, formatExpenseAmount(exp), exp.Description, r, exp.ID)
	}
	return nil
}

func canApproveOrRejectExpense(ctx context.Context, projects expensesports.ProjectMembership, projectID, actorID, role string) (bool, error) {
	if userdomain.Role(role) == userdomain.RoleAdmin {
		return true, nil
	}
	return projects.IsProjectCoordinator(ctx, projectID, actorID)
}
