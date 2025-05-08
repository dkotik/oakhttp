package oakhttp

// TODO: ponder https://vladimir.varank.in/notes/2022/05/a-real-life-use-case-for-generics-in-go-api-for-client-side-pagination/?utm_source=pocket_mylist

type HashableByPage[T any] interface {
	*T
	Hash(page, perPage int) string
}

type PageRequest[T any, H HashableByPage[T]] struct {
	Query   H
	Page    int
	PerPage int
}

type Page[T any] struct {
	Items []*T
	Page  int
	Total int
}
