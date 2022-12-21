package oakhttp

type HashableByPage[T any] interface {
	*T
	Hash(page, perPage int) string
}

type PageRequest[T any, H HashableByPage[T]] struct {
	Query   H
	Page    int
	PerPage int
}

type Page [T]struct {
	Items []*T
	Page  int
	Total int
}
