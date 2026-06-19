package newshttp

import (
	"context"
	"mime"
	"strings"

	"github.com/gin-gonic/gin"

	newsports "macabi-back/internal/news/application/ports"
	newsusecases "macabi-back/internal/news/application/usecases"
	newsdomain "macabi-back/internal/news/domain"
	newsdto "macabi-back/internal/news/infrastructure/http/dto"
	"macabi-back/internal/shared/pagination"
)

type Handler struct {
	create          *newsusecases.CreateNews
	createWithImage *newsusecases.CreateNewsWithImage
	getNews         *newsusecases.GetNews
	listPublished   *newsusecases.ListPublishedNews
	listAll         *newsusecases.ListAllNews
	getLatest       *newsusecases.GetLatestNews
	upd             *newsusecases.UpdateNews
	deleteUC        *newsusecases.DeleteNews
	imageUpload     *newsusecases.ImageUploadFile
	imageView       *newsusecases.ImageDownloadURL
	removeImage     *newsusecases.RemoveImage
	listNotifs      *newsusecases.ListNewsNotifications
	markNotifRead   *newsusecases.MarkNewsNotificationRead
	markAllNotifs   *newsusecases.MarkAllNewsNotificationsRead
	unreadNotifs    *newsusecases.UnreadNewsNotificationCount

	// imageStore resuelve URLs firmadas inline al serializar (nil si no hay storage).
	imageStore newsports.ImageStorage
}

func NewHandler(
	create *newsusecases.CreateNews,
	createWithImage *newsusecases.CreateNewsWithImage,
	getNews *newsusecases.GetNews,
	listPublished *newsusecases.ListPublishedNews,
	listAll *newsusecases.ListAllNews,
	getLatest *newsusecases.GetLatestNews,
	upd *newsusecases.UpdateNews,
	deleteUC *newsusecases.DeleteNews,
	imageUpload *newsusecases.ImageUploadFile,
	imageView *newsusecases.ImageDownloadURL,
	removeImage *newsusecases.RemoveImage,
	listNotifs *newsusecases.ListNewsNotifications,
	markNotifRead *newsusecases.MarkNewsNotificationRead,
	markAllNotifs *newsusecases.MarkAllNewsNotificationsRead,
	unreadNotifs *newsusecases.UnreadNewsNotificationCount,
	imageStore newsports.ImageStorage,
) *Handler {
	return &Handler{
		create:          create,
		createWithImage: createWithImage,
		getNews:         getNews,
		listPublished:   listPublished,
		listAll:         listAll,
		getLatest:       getLatest,
		upd:             upd,
		deleteUC:        deleteUC,
		imageUpload:     imageUpload,
		imageView:       imageView,
		removeImage:     removeImage,
		listNotifs:      listNotifs,
		markNotifRead:   markNotifRead,
		markAllNotifs:   markAllNotifs,
		unreadNotifs:    unreadNotifs,
		imageStore:      imageStore,
	}
}

// toResp serializa una noticia resolviendo su URL firmada de imagen (best-effort).
func (h *Handler) toResp(ctx context.Context, n newsdomain.News) newsdto.NewsResponse {
	imageURL := ""
	if h.imageStore != nil && n.ImageStoragePath != nil && *n.ImageStoragePath != "" {
		if u, err := h.imageStore.CreateSignedDownloadURL(ctx, *n.ImageStoragePath, 3600); err == nil {
			imageURL = u
		}
	}
	return newsdto.ToNewsResponse(n, imageURL)
}

func (h *Handler) toListResp(ctx context.Context, res pagination.Result[newsdomain.News]) pagination.Result[newsdto.NewsResponse] {
	items := make([]newsdto.NewsResponse, len(res.Data))
	for i := range res.Data {
		items[i] = h.toResp(ctx, res.Data[i])
	}
	return newsdto.NewListResponse(items, res)
}

func isMultipartRequest(c *gin.Context) bool {
	ct := strings.TrimSpace(c.GetHeader("Content-Type"))
	if ct == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return strings.Contains(strings.ToLower(ct), "multipart/form-data")
	}
	return mediaType == "multipart/form-data"
}
