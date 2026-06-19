package newsusecases

import (
	"context"
	"io"
	"time"

	newsports "macabi-back/internal/news/application/ports"
	newsdomain "macabi-back/internal/news/domain"
)

// CreateNewsWithImage crea la noticia, adjunta la imagen y recién entonces publica,
// para no notificar una noticia cuya imagen falló (en ese caso se elimina el borrador).
type CreateNewsWithImage struct {
	repo        newsports.NewsRepository
	imageUpload *ImageUploadFile
	notifier    NewsNotifier
}

func NewCreateNewsWithImage(repo newsports.NewsRepository, imageUpload *ImageUploadFile, notifier NewsNotifier) *CreateNewsWithImage {
	return &CreateNewsWithImage{repo: repo, imageUpload: imageUpload, notifier: notifier}
}

type CreateNewsWithImageInput struct {
	CreateNewsInput
	ImageContentType string
	ImageSize        int64
	ImageBody        io.Reader
}

func (uc *CreateNewsWithImage) Execute(ctx context.Context, in CreateNewsWithImageInput) (*newsdomain.News, error) {
	// 1. Crear siempre como borrador (no se notifica hasta tener la imagen lista).
	n, err := newsdomain.NewNews(in.Title, in.Body, in.AuthorID)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, n); err != nil {
		return nil, err
	}

	// 2. Adjuntar imagen; si falla, deshacer la noticia.
	if in.ImageBody != nil {
		if _, err := uc.imageUpload.Execute(ctx, ImageUploadFileInput{
			NewsID:      n.ID,
			ContentType: in.ImageContentType,
			Size:        in.ImageSize,
			Body:        in.ImageBody,
		}); err != nil {
			_ = uc.repo.DeleteByID(ctx, n.ID)
			return nil, newsdomain.ErrImageAttachFailed
		}
	}

	// 3. Asociar proyectos antes de publicar (para acertar la audiencia al notificar).
	if len(in.ProjectIDs) > 0 {
		if err := uc.repo.SetNewsProjects(ctx, n.ID, in.ProjectIDs); err != nil {
			_ = uc.repo.DeleteByID(ctx, n.ID)
			return nil, err
		}
	}

	// 4. Publicar (si corresponde) con la imagen y los proyectos ya asociados.
	fresh, err := uc.repo.FindByID(ctx, n.ID)
	if err != nil {
		return nil, err
	}
	if in.Publish {
		fresh.Publish(time.Now())
		if err := uc.repo.Update(ctx, fresh); err != nil {
			return nil, err
		}
		DispatchNewsNotify(uc.notifier, fresh)
	}
	return fresh, nil
}
