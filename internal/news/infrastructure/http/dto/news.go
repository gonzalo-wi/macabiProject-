package newsdto

import (
	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"
)

const TimeFormat = "2006-01-02T15:04:05Z07:00"

type CreateNewsBody struct {
	Title      string   `json:"title" binding:"required"`
	Body       string   `json:"body" binding:"required"`
	Publish    bool     `json:"publish"`
	ProjectIDs []string `json:"project_ids"`
}

type PatchNewsBody struct {
	Title      *string   `json:"title"`
	Body       *string   `json:"body"`
	Publish    *bool     `json:"publish"`
	Renotify   bool      `json:"renotify"`
	ProjectIDs *[]string `json:"project_ids"`
}

type NewsResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	Status      string   `json:"status"`
	AuthorID    string   `json:"author_id"`
	ImageURL    string   `json:"image_url,omitempty"`
	HasImage    bool     `json:"has_image"`
	ProjectIDs  []string `json:"project_ids"`
	PublishedAt *string  `json:"published_at"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// ToNewsResponse serializa una noticia. imageURL es la URL firmada (vacía si no hay imagen).
func ToNewsResponse(n newsdomain.News, imageURL string) NewsResponse {
	pids := n.ProjectIDs
	if pids == nil {
		pids = []string{}
	}
	resp := NewsResponse{
		ID:         n.ID,
		Title:      n.Title,
		Body:       n.Body,
		Status:     string(n.Status),
		AuthorID:   n.AuthorID,
		ImageURL:   imageURL,
		HasImage:   n.ImageStoragePath != nil && *n.ImageStoragePath != "",
		ProjectIDs: pids,
		CreatedAt:  n.CreatedAt.Format(TimeFormat),
		UpdatedAt:  n.UpdatedAt.Format(TimeFormat),
	}
	if n.PublishedAt != nil {
		s := n.PublishedAt.Format(TimeFormat)
		resp.PublishedAt = &s
	}
	return resp
}

// NewListResponse arma una página de NewsResponse a partir de los items ya serializados.
func NewListResponse(items []NewsResponse, src pagination.Result[newsdomain.News]) pagination.Result[NewsResponse] {
	return pagination.Result[NewsResponse]{
		Data:       items,
		Total:      src.Total,
		Page:       src.Page,
		PageSize:   src.PageSize,
		TotalPages: src.TotalPages,
	}
}
