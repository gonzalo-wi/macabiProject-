package newsusecases

import (
	"context"

	newsports "macabi-back/internal/news/application/ports"
	newsdomain "macabi-back/internal/news/domain"
	"macabi-back/internal/shared/pagination"
)

type GetNews struct {
	repo newsports.NewsRepository
}

func NewGetNews(repo newsports.NewsRepository) *GetNews {
	return &GetNews{repo: repo}
}

func (uc *GetNews) Execute(ctx context.Context, id string) (*newsdomain.News, error) {
	return uc.repo.FindByID(ctx, id)
}

type ListPublishedNews struct {
	repo     newsports.NewsRepository
	audience newsports.ProjectAudienceReader
}

func NewListPublishedNews(repo newsports.NewsRepository, audience newsports.ProjectAudienceReader) *ListPublishedNews {
	return &ListPublishedNews{repo: repo, audience: audience}
}

func (uc *ListPublishedNews) Execute(ctx context.Context, userID string, params pagination.Params) (pagination.Result[newsdomain.News], error) {
	pids, err := uc.audience.FindUserProjectIDs(ctx, userID)
	if err != nil {
		return pagination.Result[newsdomain.News]{}, err
	}
	return uc.repo.ListPublished(ctx, pids, params)
}

type ListAllNews struct {
	repo newsports.NewsRepository
}

func NewListAllNews(repo newsports.NewsRepository) *ListAllNews {
	return &ListAllNews{repo: repo}
}

func (uc *ListAllNews) Execute(ctx context.Context, params pagination.Params) (pagination.Result[newsdomain.News], error) {
	return uc.repo.ListAll(ctx, params)
}

type GetLatestNews struct {
	repo     newsports.NewsRepository
	audience newsports.ProjectAudienceReader
}

func NewGetLatestNews(repo newsports.NewsRepository, audience newsports.ProjectAudienceReader) *GetLatestNews {
	return &GetLatestNews{repo: repo, audience: audience}
}

// Execute devuelve la última noticia publicada VISIBLE para el usuario,
// o (nil, nil) si todavía no hay ninguna.
func (uc *GetLatestNews) Execute(ctx context.Context, userID string) (*newsdomain.News, error) {
	pids, err := uc.audience.FindUserProjectIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	return uc.repo.FindLatestPublished(ctx, pids)
}
