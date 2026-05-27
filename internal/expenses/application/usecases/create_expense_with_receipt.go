package expensesusecases

import (
	"context"
	"io"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
)

type CreateExpenseWithReceipt struct {
	create  *CreateExpense
	receipt *ReceiptUploadFile
	repo    expensesports.ExpenseRepository
}

func NewCreateExpenseWithReceipt(
	create *CreateExpense,
	receipt *ReceiptUploadFile,
	repo expensesports.ExpenseRepository,
) *CreateExpenseWithReceipt {
	return &CreateExpenseWithReceipt{create: create, receipt: receipt, repo: repo}
}

type CreateExpenseWithReceiptInput struct {
	CreateExpenseInput
	ReceiptContentType string
	ReceiptSize        int64
	ReceiptBody        io.Reader
}

func (uc *CreateExpenseWithReceipt) Execute(ctx context.Context, in CreateExpenseWithReceiptInput) (*expensesdomain.Expense, error) {
	exp, err := uc.create.Execute(ctx, in.CreateExpenseInput)
	if err != nil {
		return nil, err
	}

	if in.ReceiptBody == nil {
		return exp, nil
	}

	_, err = uc.receipt.Execute(ctx, ReceiptUploadFileInput{
		ExpenseID:   exp.ID,
		ActorID:     in.SubmittedBy,
		UserRole:    in.UserRole,
		ContentType: in.ReceiptContentType,
		Size:        in.ReceiptSize,
		Body:        in.ReceiptBody,
	})
	if err != nil {
		_ = uc.repo.DeleteByID(ctx, exp.ID)
		return nil, expensesdomain.ErrReceiptAttachFailed
	}

	updated, err := uc.repo.FindByID(ctx, exp.ID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
