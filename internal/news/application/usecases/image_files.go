package newsusecases

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	newsports "macabi-back/internal/news/application/ports"
	newsdomain "macabi-back/internal/news/domain"
)

// ImageMIME mapea los tipos de imagen aceptados a su extensión. Solo imágenes (sin PDF).
var ImageMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

const MaxImageBytes = 2 * 1024 * 1024

func bytesReader(data []byte) *bytes.Reader { return bytes.NewReader(data) }

// ValidateImageBytes valida tamaño y detecta el MIME real por signature.
func ValidateImageBytes(body io.Reader, declaredSize int64, declaredContentType string) ([]byte, string, error) {
	if body == nil {
		return nil, "", newsdomain.ErrMissingRequiredField
	}
	if declaredSize < 0 || declaredSize > MaxImageBytes {
		return nil, "", newsdomain.ErrImageTooLarge
	}

	limited := io.LimitReader(body, MaxImageBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", newsdomain.ErrImageAttachFailed
	}
	if len(data) == 0 {
		return nil, "", newsdomain.ErrMissingRequiredField
	}
	if len(data) > MaxImageBytes {
		return nil, "", newsdomain.ErrImageTooLarge
	}

	detected := detectImageMIME(data)
	if detected == "" {
		return nil, "", newsdomain.ErrInvalidMimeType
	}

	declared := normalizeDeclaredMIME(declaredContentType)
	if declared != "" && declared != detected {
		if ImageMIME[declared] != ImageMIME[detected] {
			return nil, "", newsdomain.ErrInvalidMimeType
		}
	}

	return data, detected, nil
}

func normalizeDeclaredMIME(ct string) string {
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "" {
		return ""
	}
	if i := strings.Index(ct, ";"); i >= 0 {
		ct = strings.TrimSpace(ct[:i])
	}
	return ct
}

func detectImageMIME(data []byte) string {
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}
	if len(data) >= 8 && data[0] == 0x89 && string(data[1:4]) == "PNG" {
		return "image/png"
	}
	if len(data) >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff {
		return "image/jpeg"
	}
	sniff := normalizeDeclaredMIME(http.DetectContentType(data))
	if _, ok := ImageMIME[sniff]; ok {
		return sniff
	}
	return ""
}

// ImageUploadFile valida y sube la imagen, guardando el path en la noticia.
// La autorización (admin) la garantiza el router.
type ImageUploadFile struct {
	store newsports.ImageStorage
	repo  newsports.NewsRepository
}

func NewImageUploadFile(store newsports.ImageStorage, repo newsports.NewsRepository) *ImageUploadFile {
	return &ImageUploadFile{store: store, repo: repo}
}

type ImageUploadFileInput struct {
	NewsID      string
	ContentType string
	Size        int64
	Body        io.Reader
}

type ImageUploadFileResult struct {
	Path string `json:"path"`
}

func (uc *ImageUploadFile) Execute(ctx context.Context, in ImageUploadFileInput) (*ImageUploadFileResult, error) {
	if uc.store == nil {
		return nil, newsdomain.ErrImageStorageUnavailable
	}
	if in.Body == nil {
		return nil, newsdomain.ErrMissingRequiredField
	}

	data, ct, err := ValidateImageBytes(in.Body, in.Size, in.ContentType)
	if err != nil {
		return nil, err
	}
	suffix, ok := ImageMIME[ct]
	if !ok {
		return nil, newsdomain.ErrInvalidMimeType
	}

	n, err := uc.repo.FindByID(ctx, in.NewsID)
	if err != nil {
		return nil, err
	}

	objectKey := newsdomain.ImageKeyPrefix(n.ID) + "/image" + suffix

	if err := uc.store.UploadObject(ctx, objectKey, ct, bytesReader(data)); err != nil {
		return nil, newsdomain.ErrImageStorageUnavailable
	}

	// Si había una imagen con otra extensión, eliminar la huérfana.
	if n.ImageStoragePath != nil && *n.ImageStoragePath != "" && *n.ImageStoragePath != objectKey {
		_ = uc.store.DeleteObject(ctx, *n.ImageStoragePath)
	}

	n.ImageStoragePath = &objectKey
	if err := uc.repo.Update(ctx, n); err != nil {
		return nil, err
	}

	return &ImageUploadFileResult{Path: objectKey}, nil
}

// ImageDownloadURL genera una URL firmada de descarga para la imagen de una noticia.
type ImageDownloadURL struct {
	store newsports.ImageStorage
	repo  newsports.NewsRepository
}

func NewImageDownloadURL(store newsports.ImageStorage, repo newsports.NewsRepository) *ImageDownloadURL {
	return &ImageDownloadURL{store: store, repo: repo}
}

func (uc *ImageDownloadURL) Execute(ctx context.Context, newsID string) (string, error) {
	if uc.store == nil {
		return "", newsdomain.ErrImageStorageUnavailable
	}
	n, err := uc.repo.FindByID(ctx, newsID)
	if err != nil {
		return "", err
	}
	if n.ImageStoragePath == nil || *n.ImageStoragePath == "" {
		return "", newsdomain.ErrNoImageAttached
	}
	if !newsdomain.ValidateImageObjectKey(n.ID, *n.ImageStoragePath) {
		return "", newsdomain.ErrInvalidImagePath
	}
	return uc.store.CreateSignedDownloadURL(ctx, *n.ImageStoragePath, 3600)
}

// RemoveImage borra la imagen del storage y limpia el path en la noticia.
type RemoveImage struct {
	store newsports.ImageStorage
	repo  newsports.NewsRepository
}

func NewRemoveImage(store newsports.ImageStorage, repo newsports.NewsRepository) *RemoveImage {
	return &RemoveImage{store: store, repo: repo}
}

func (uc *RemoveImage) Execute(ctx context.Context, newsID string) error {
	n, err := uc.repo.FindByID(ctx, newsID)
	if err != nil {
		return err
	}
	if n.ImageStoragePath == nil || *n.ImageStoragePath == "" {
		return nil
	}
	if uc.store != nil {
		_ = uc.store.DeleteObject(ctx, *n.ImageStoragePath)
	}
	n.ImageStoragePath = nil
	return uc.repo.Update(ctx, n)
}
