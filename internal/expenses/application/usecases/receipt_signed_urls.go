package expensesusecases

import (
	"context"
	"io"
	"strings"

	expensesports "macabi-back/internal/expenses/application/ports"
	expensesdomain "macabi-back/internal/expenses/domain"
	projectports "macabi-back/internal/project/application/ports"
	userdomain "macabi-back/internal/user/domain"
)

var ReceiptMIME = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
}

const MaxReceiptBytes = 2 * 1024 * 1024

type ReceiptUploadURL struct {
	signer expensesports.ReceiptSigner
	repo   expensesports.ExpenseRepository
	proj   projectports.ProjectMembership
}

func NewReceiptUploadURL(signer expensesports.ReceiptSigner, repo expensesports.ExpenseRepository, proj projectports.ProjectMembership) *ReceiptUploadURL {
	return &ReceiptUploadURL{signer: signer, repo: repo, proj: proj}
}

type ReceiptUploadInput struct {
	ExpenseID   string
	ActorID     string
	UserRole    string
	ContentType string
}

type ReceiptUploadResult struct {
	UploadURL     string            `json:"upload_url"`
	Path          string            `json:"path"`
	UploadHeaders map[string]string `json:"upload_headers,omitempty"`
}

func (uc *ReceiptUploadURL) Execute(ctx context.Context, in ReceiptUploadInput) (*ReceiptUploadResult, error) {
	if uc.signer == nil {
		return nil, expensesdomain.ErrReceiptStorageUnavailable
	}

	ct := strings.TrimSpace(strings.ToLower(in.ContentType))
	suffix, ok := ReceiptMIME[ct]
	if !ok {
		return nil, expensesdomain.ErrInvalidMimeType
	}

	exp, err := uc.repo.FindByID(ctx, in.ExpenseID)
	if err != nil {
		return nil, err
	}

	isAdmin := userdomain.Role(in.UserRole) == userdomain.RoleAdmin
	coord, err := uc.proj.IsProjectCoordinator(ctx, exp.ProjectID, in.ActorID)
	if err != nil {
		return nil, err
	}
	owner := exp.SubmittedByUserID == in.ActorID

	can := false
	if isAdmin || coord {
		can = true
	}
	if owner && exp.Status == expensesdomain.StatusPending {
		can = true
	}
	if !can {
		return nil, expensesdomain.ErrForbidden
	}

	objectKey := expensesdomain.ReceiptKeyPrefix(exp.ProjectID, exp.ID) + "/receipt" + suffix

	uploadURL, err := uc.signer.CreateSignedUploadURL(ctx, objectKey, ct)
	if err != nil {
		return nil, expensesdomain.ErrReceiptStorageUnavailable
	}

	return &ReceiptUploadResult{
		UploadURL:     uploadURL,
		Path:          objectKey,
		UploadHeaders: map[string]string{"Content-Type": ct},
	}, nil
}

type ReceiptDownloadURL struct {
	signer expensesports.ReceiptSigner
	repo   expensesports.ExpenseRepository
	proj   projectports.ProjectMembership
}

func NewReceiptDownloadURL(signer expensesports.ReceiptSigner, repo expensesports.ExpenseRepository, proj projectports.ProjectMembership) *ReceiptDownloadURL {
	return &ReceiptDownloadURL{signer: signer, repo: repo, proj: proj}
}

type ReceiptDownloadInput struct {
	ExpenseID string
	UserID    string
	UserRole  string
}

type ReceiptDownloadResult struct {
	ViewURL string `json:"view_url"`
}

func (uc *ReceiptDownloadURL) Execute(ctx context.Context, in ReceiptDownloadInput) (*ReceiptDownloadResult, error) {
	if uc.signer == nil {
		return nil, expensesdomain.ErrReceiptStorageUnavailable
	}

	exp, err := uc.repo.FindByID(ctx, in.ExpenseID)
	if err != nil {
		return nil, err
	}
	if exp.ReceiptStoragePath == nil || *exp.ReceiptStoragePath == "" {
		return nil, expensesdomain.ErrNoReceiptAttached
	}
	if !expensesdomain.ValidateReceiptObjectKey(exp.ProjectID, exp.ID, *exp.ReceiptStoragePath) {
		return nil, expensesdomain.ErrInvalidReceiptPath
	}

	isAdmin := userdomain.Role(in.UserRole) == userdomain.RoleAdmin
	if !isAdmin && exp.SubmittedByUserID != in.UserID {
		coord, err := uc.proj.IsProjectCoordinator(ctx, exp.ProjectID, in.UserID)
		if err != nil {
			return nil, err
		}
		if !coord {
			return nil, expensesdomain.ErrForbidden
		}
	}

	viewURL, err := uc.signer.CreateSignedDownloadURL(ctx, *exp.ReceiptStoragePath, 3600)
	if err != nil {
		return nil, expensesdomain.ErrReceiptStorageUnavailable
	}

	return &ReceiptDownloadResult{ViewURL: viewURL}, nil
}

type ReceiptUploadFile struct {
	signer expensesports.ReceiptSigner
	repo   expensesports.ExpenseRepository
	proj   projectports.ProjectMembership
}

func NewReceiptUploadFile(signer expensesports.ReceiptSigner, repo expensesports.ExpenseRepository, proj projectports.ProjectMembership) *ReceiptUploadFile {
	return &ReceiptUploadFile{signer: signer, repo: repo, proj: proj}
}

type ReceiptUploadFileInput struct {
	ExpenseID   string
	ActorID     string
	UserRole    string
	ContentType string
	Size        int64
	Body        io.Reader
}

type ReceiptUploadFileResult struct {
	Path string `json:"path"`
}

func (uc *ReceiptUploadFile) Execute(ctx context.Context, in ReceiptUploadFileInput) (*ReceiptUploadFileResult, error) {
	if uc.signer == nil {
		return nil, expensesdomain.ErrReceiptStorageUnavailable
	}
	if in.Body == nil {
		return nil, expensesdomain.ErrMissingRequiredField
	}

	data, ct, err := ValidateReceiptBytes(in.Body, in.Size, in.ContentType)
	if err != nil {
		return nil, err
	}
	suffix, ok := ReceiptMIME[ct]
	if !ok {
		return nil, expensesdomain.ErrInvalidMimeType
	}

	exp, err := uc.repo.FindByID(ctx, in.ExpenseID)
	if err != nil {
		return nil, err
	}

	isAdmin := userdomain.Role(in.UserRole) == userdomain.RoleAdmin
	coord, err := uc.proj.IsProjectCoordinator(ctx, exp.ProjectID, in.ActorID)
	if err != nil {
		return nil, err
	}
	owner := exp.SubmittedByUserID == in.ActorID

	can := false
	if isAdmin || coord {
		can = true
	}
	if owner && exp.Status == expensesdomain.StatusPending {
		can = true
	}
	if !can {
		return nil, expensesdomain.ErrForbidden
	}

	objectKey := expensesdomain.ReceiptKeyPrefix(exp.ProjectID, exp.ID) + "/receipt" + suffix

	if err := uc.signer.UploadObject(ctx, objectKey, ct, bytesReader(data)); err != nil {
		return nil, expensesdomain.ErrReceiptStorageUnavailable
	}

	exp.ReceiptStoragePath = &objectKey
	if err := uc.repo.Update(ctx, exp); err != nil {
		return nil, err
	}

	return &ReceiptUploadFileResult{Path: objectKey}, nil
}
