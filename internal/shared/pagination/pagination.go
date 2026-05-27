package pagination

import "strconv"

const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

type Params struct {
	Page     int
	PageSize int
}

type Result[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

func NewParams(page, pageSize int) Params {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = DefaultPageSize
	}
	return Params{Page: page, PageSize: pageSize}
}

func ParseParams(pageStr, pageSizeStr string) Params {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = DefaultPage
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = DefaultPageSize
	}
	return NewParams(page, pageSize)
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func NewResult[T any](data []T, total int64, params Params) Result[T] {
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}
	return Result[T]{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}
}
