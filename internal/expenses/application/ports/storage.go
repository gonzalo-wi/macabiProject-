package expensesports

import (
	"context"
	"io"
)

type ReceiptSigner interface {
	CreateSignedUploadURL(ctx context.Context, objectKey, contentType string) (uploadURL string, err error)
	CreateSignedDownloadURL(ctx context.Context, objectKey string, expiresSec int) (downloadURL string, err error)
	UploadObject(ctx context.Context, objectKey, contentType string, body io.Reader) error
	DeleteObject(ctx context.Context, objectKey string) error
}
