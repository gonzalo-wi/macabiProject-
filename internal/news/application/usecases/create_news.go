package newsusecases

import (
	"context"
	"time"

	newsports "macabi-back/internal/news/application/ports"
	newsdomain "macabi-back/internal/news/domain"
)

// NewsNotifier dispara los 3 canales a todos los miembros cuando se publica una noticia.
type NewsNotifier interface {
	NotifyAllMembersNewNews(ctx context.Context, news *newsdomain.News)
}

// DispatchNewsNotify ejecuta el notifier en una goroutine independiente: la
// publicación nunca se bloquea ni falla por las notificaciones.
func DispatchNewsNotify(notifier NewsNotifier, n *newsdomain.News) {
	if notifier == nil || n == nil {
		return
	}
	snapshot := *n
	go notifier.NotifyAllMembersNewNews(context.Background(), &snapshot)
}

type CreateNews struct {
	repo     newsports.NewsRepository
	notifier NewsNotifier
}

func NewCreateNews(repo newsports.NewsRepository, notifier NewsNotifier) *CreateNews {
	return &CreateNews{repo: repo, notifier: notifier}
}

type CreateNewsInput struct {
	Title    string
	Body     string
	AuthorID string
	Publish  bool
	// ProjectIDs: proyectos destino (vacío = todos).
	ProjectIDs []string
}

func (uc *CreateNews) Execute(ctx context.Context, in CreateNewsInput) (*newsdomain.News, error) {
	n, err := newsdomain.NewNews(in.Title, in.Body, in.AuthorID)
	if err != nil {
		return nil, err
	}
	if in.Publish {
		n.Publish(time.Now())
	}
	if err := uc.repo.Save(ctx, n); err != nil {
		return nil, err
	}
	// Asociar proyectos (tras obtener el ID) antes de notificar, para acertar la audiencia.
	if len(in.ProjectIDs) > 0 {
		if err := uc.repo.SetNewsProjects(ctx, n.ID, in.ProjectIDs); err != nil {
			return nil, err
		}
	}
	n.ProjectIDs = in.ProjectIDs
	if n.IsPublished() {
		DispatchNewsNotify(uc.notifier, n)
	}
	return n, nil
}
