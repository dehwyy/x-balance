package storage

type Pagination struct {
	limit  int
	offset int
}

func NewPagination(limit, offset int) Pagination {
	if limit <= 0 {
		limit = 20
	}
	return Pagination{limit: limit, offset: offset}
}

func (p Pagination) Limit() int  { return p.limit }
func (p Pagination) Offset() int { return p.offset }

type PageResult[T any] struct {
	Items      []T
	Pagination Pagination
}

func NewPageResult[T any](items []T, pagination Pagination) *PageResult[T] {
	return &PageResult[T]{Items: items, Pagination: pagination}
}
