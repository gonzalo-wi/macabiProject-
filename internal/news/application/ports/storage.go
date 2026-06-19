package newsports

import (
	"context"
	"io"
)

// ImageStorage es satisfecho por sharedstorage.SupabaseSigner.
type ImageStorage interface {
	CreateSignedDownloadURL(ctx context.Context, objectKey string, expiresSec int) (downloadURL string, err error)
	UploadObject(ctx context.Context, objectKey, contentType string, body io.Reader) error
	DeleteObject(ctx context.Context, objectKey string) error
}
