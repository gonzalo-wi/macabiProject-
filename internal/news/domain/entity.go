package newsdomain

import (
	"fmt"
	"strings"
	"time"
)

type NewsStatus string

const (
	StatusDraft     NewsStatus = "draft"
	StatusPublished NewsStatus = "published"
)

// News es un posteo global (no ligado a proyecto) publicado por un admin.
type News struct {
	ID               string
	Title            string
	Body             string
	ImageStoragePath *string
	Status           NewsStatus
	AuthorID         string
	// ProjectIDs: proyectos a los que se dirige la noticia. Vacío = todos.
	ProjectIDs  []string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewsNotification es la entrada de campana in-app para un usuario.
type NewsNotification struct {
	ID        string
	UserID    string
	NewsID    string
	Message   string
	ReadAt    *time.Time
	CreatedAt time.Time
}

func NewNews(title, body, authorID string) (*News, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	authorID = strings.TrimSpace(authorID)
	if title == "" || body == "" || authorID == "" {
		return nil, ErrMissingRequiredField
	}
	return &News{
		Title:    title,
		Body:     body,
		AuthorID: authorID,
		Status:   StatusDraft,
	}, nil
}

// Publish marca la noticia como publicada y fija PublishedAt la primera vez.
func (n *News) Publish(at time.Time) {
	n.Status = StatusPublished
	if n.PublishedAt == nil {
		t := at.UTC()
		n.PublishedAt = &t
	}
}

func (n *News) IsPublished() bool {
	return n.Status == StatusPublished
}

// ImageKeyPrefix es el prefijo del object key en el bucket de imágenes.
func ImageKeyPrefix(newsID string) string {
	return fmt.Sprintf("news/%s", newsID)
}

// ValidateImageObjectKey evita rutas fuera del prefijo esperado de la noticia.
func ValidateImageObjectKey(newsID, objectKey string) bool {
	prefix := ImageKeyPrefix(newsID) + "/"
	if !strings.HasPrefix(objectKey, prefix) {
		return false
	}
	rest := strings.TrimPrefix(objectKey, prefix)
	if rest == "" || strings.Contains(rest, "/") || strings.Contains(rest, "..") {
		return false
	}
	return true
}
