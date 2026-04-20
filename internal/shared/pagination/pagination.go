package pagination

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
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return Params{Page: page, PageSize: pageSize}
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
