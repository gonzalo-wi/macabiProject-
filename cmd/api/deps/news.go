package deps

import (
	"strings"

	newsports "macabi-back/internal/news/application/ports"
	newsservices "macabi-back/internal/news/application/services"
	newsusecases "macabi-back/internal/news/application/usecases"
	newshttp "macabi-back/internal/news/infrastructure/http"
	newsmail "macabi-back/internal/news/infrastructure/mail"
	newspersistence "macabi-back/internal/news/infrastructure/persistence"
	"macabi-back/internal/shared/config"
	"macabi-back/internal/shared/notifications"
	sharedstorage "macabi-back/internal/shared/storage"
	userports "macabi-back/internal/user/application/ports"

	"gorm.io/gorm"
)

func buildNewsDeps(
	db *gorm.DB,
	cfg *config.Config,
	members userports.MemberDirectory,
	audience newsports.ProjectAudienceReader,
	push notifications.PushNotifier,
) *newshttp.Handler {
	newsRepo := newspersistence.NewNewsRepositoryPG(db)
	if err := newspersistence.RunMigrations(db); err != nil {
		panic("news migrations failed: " + err.Error())
	}

	// Bucket de imágenes: usa el propio de noticias si está seteado; si no, reutiliza
	// el de comprobantes de gastos (mismo proyecto Supabase) para no exigir config nueva.
	// Las imágenes van con prefijo "news/<id>/", sin colisión con los comprobantes.
	newsBucket := cfg.SupabaseNewsImageBucket
	if newsBucket == "" {
		newsBucket = cfg.SupabaseExpenseReceiptBucket
	}
	var imageStore newsports.ImageStorage
	if cfg.SupabaseURL != "" && cfg.SupabaseServiceRoleKey != "" && newsBucket != "" {
		imageStore = &sharedstorage.SupabaseSigner{
			BaseURL: strings.TrimSuffix(cfg.SupabaseURL, "/"),
			APIKey:  cfg.SupabaseServiceRoleKey,
			Bucket:  newsBucket,
		}
	}

	var mailer newsports.NewsMailer
	if cfg.BrevoAPIKey != "" && cfg.BrevoEmailFrom != "" {
		mailer = newsmail.NewBrevoNewsMailer(cfg.BrevoAPIKey, cfg.BrevoEmailFrom, cfg.FrontendPublicURL)
	} else {
		mailer = newsmail.NewNoOpNewsMailer()
	}

	notifier := newsservices.NewNewsNotificationService(newsRepo, members, audience, mailer, push)

	imageUploadUC := newsusecases.NewImageUploadFile(imageStore, newsRepo)

	return newshttp.NewHandler(
		newsusecases.NewCreateNews(newsRepo, notifier),
		newsusecases.NewCreateNewsWithImage(newsRepo, imageUploadUC, notifier),
		newsusecases.NewGetNews(newsRepo),
		newsusecases.NewListPublishedNews(newsRepo, audience),
		newsusecases.NewListAllNews(newsRepo),
		newsusecases.NewGetLatestNews(newsRepo, audience),
		newsusecases.NewUpdateNews(newsRepo, notifier),
		newsusecases.NewDeleteNews(newsRepo, newsRepo, imageStore),
		imageUploadUC,
		newsusecases.NewImageDownloadURL(imageStore, newsRepo),
		newsusecases.NewRemoveImage(imageStore, newsRepo),
		newsusecases.NewListNewsNotifications(newsRepo),
		newsusecases.NewMarkNewsNotificationRead(newsRepo),
		newsusecases.NewMarkAllNewsNotificationsRead(newsRepo),
		newsusecases.NewUnreadNewsNotificationCount(newsRepo),
		imageStore,
	)
}
