package newsservices

import (
	"context"
	"fmt"
	"strings"

	newsports "macabi-back/internal/news/application/ports"
	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/notifications"
	userports "macabi-back/internal/user/application/ports"
)

// NewsNotificationService unifica campana in-app, Web Push y email para noticias.
type NewsNotificationService struct {
	notifs   newsports.NewsNotificationRepository
	members  userports.MemberDirectory
	audience newsports.ProjectAudienceReader
	mailer   newsports.NewsMailer
	push     notifications.PushNotifier
}

func NewNewsNotificationService(
	notifs newsports.NewsNotificationRepository,
	members userports.MemberDirectory,
	audience newsports.ProjectAudienceReader,
	mailer newsports.NewsMailer,
	push notifications.PushNotifier,
) *NewsNotificationService {
	return &NewsNotificationService{
		notifs:   notifs,
		members:  members,
		audience: audience,
		mailer:   mailer,
		push:     push,
	}
}

// NotifyAllMembersNewNews envía los 3 canales a todos los miembros activos menos el autor.
// Cada canal es independiente: los errores se ignoran porque la noticia ya está publicada.
// Pensado para ejecutarse en una goroutine con context.Background().
func (s *NewsNotificationService) NotifyAllMembersNewNews(ctx context.Context, news *newsdomain.News) {
	// Audiencia: si la noticia no tiene proyectos, va a todos los activos; si tiene,
	// solo a los miembros (madrijim + coordinadores) de esos proyectos.
	var members []userports.Member
	var err error
	if len(news.ProjectIDs) == 0 {
		members, err = s.members.FindAllActiveMembers(ctx)
	} else {
		members, err = s.audience.FindActiveMembersOfProjects(ctx, news.ProjectIDs)
	}
	if err != nil || len(members) == 0 {
		return
	}

	title := "Nueva noticia"
	body := news.Title
	url := notifications.NewsDetail(news.ID)
	msg := fmt.Sprintf("Nueva noticia: %s", news.Title)

	notifs := make([]*newsdomain.NewsNotification, 0, len(members))
	emails := make([]string, 0, len(members))
	for _, m := range members {
		if m.ID == news.AuthorID {
			continue
		}
		notifs = append(notifs, &newsdomain.NewsNotification{
			UserID:  m.ID,
			NewsID:  news.ID,
			Message: msg,
		})
		notifications.PushToUser(ctx, s.push, m.ID, title, body, url)
		if e := strings.TrimSpace(m.Email); e != "" {
			emails = append(emails, e)
		}
	}

	if len(notifs) > 0 {
		_ = s.notifs.SaveNotifications(ctx, notifs)
	}
	if len(emails) > 0 {
		_ = s.mailer.NotifyMembersNewNews(ctx, emails, news.Title, summarize(news.Body), news.ID)
	}
}

// summarize recorta el cuerpo para el resumen del email.
func summarize(body string) string {
	body = strings.TrimSpace(body)
	const max = 200
	r := []rune(body)
	if len(r) <= max {
		return body
	}
	return strings.TrimSpace(string(r[:max])) + "…"
}
