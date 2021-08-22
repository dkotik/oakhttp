package oakquery

type Query struct {
	Filters []QueryFilter
	Flags   []string
	Tags    []string
	Page    uint32
	PerPage uint32
	Total   uint64 // can be mutated
	// is range needed?
	Range QueryRange // Range will be mutated, limiting PerPage, updating Total
}

type QueryRange struct {
	Page    uint32
	PerPage uint32
	Total   uint64
}

func (q *Query) Is(flags ...string) bool {
	found := 0
	for _, flag := range flags {
		for _, hasFlag := range q.Flags {
			if flag == hasFlag {
				found++
				break
			}
		}
	}
	return found == len(flags)
}
