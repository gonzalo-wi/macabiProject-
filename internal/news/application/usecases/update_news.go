package newsusecases

import (
	"context"
	"strings"
	"time"

	newsports "macabi-back/internal/news/application/ports"
	newsdomain "macabi-back/internal/news/domain"
)

type UpdateNews struct {
	repo     newsports.NewsRepository
	notifier NewsNotifier
}

func NewUpdateNews(repo newsports.NewsRepository, notifier NewsNotifier) *UpdateNews {
	return &UpdateNews{repo: repo, notifier: notifier}
}

type UpdateNewsInput struct {
	ID    string
	Title *string
	Body  *string
	// Publish: true => publicar; false => volver a borrador; nil => sin cambio de estado.
	Publish *bool
	// Renotify reenvía las notificaciones aunque la noticia ya estuviera publicada.
	Renotify bool
	// ProjectIDs: si != nil, reemplaza los proyectos destino (vacío = todos).
	ProjectIDs *[]string
}

func (uc *UpdateNews) Execute(ctx context.Context, in UpdateNewsInput) (*newsdomain.News, error) {
	n, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		if t == "" {
			return nil, newsdomain.ErrMissingRequiredField
		}
		n.Title = t
	}
	if in.Body != nil {
		b := strings.TrimSpace(*in.Body)
		if b == "" {
			return nil, newsdomain.ErrMissingRequiredField
		}
		n.Body = b
	}

	wasPublished := n.IsPublished()
	firstPublish := false
	if in.Publish != nil {
		if *in.Publish {
			if !wasPublished {
				n.Publish(time.Now())
				firstPublish = true
			}
		} else {
			n.Status = newsdomain.StatusDraft
		}
	}

	if err := uc.repo.Update(ctx, n); err != nil {
		return nil, err
	}

	// Reemplazar proyectos destino si vinieron (antes de notificar, para la audiencia correcta).
	if in.ProjectIDs != nil {
		if err := uc.repo.SetNewsProjects(ctx, n.ID, *in.ProjectIDs); err != nil {
			return nil, err
		}
		n.ProjectIDs = *in.ProjectIDs
	}

	// Notifica al publicar por primera vez, o si el admin pidió reenviar y sigue publicada.
	if firstPublish || (n.IsPublished() && in.Renotify) {
		DispatchNewsNotify(uc.notifier, n)
	}

	return n, nil
}
