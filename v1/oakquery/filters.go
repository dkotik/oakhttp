package oakquery

type QueryFilter struct {
	Field      string
	Value      string
	Comparison QueryFilterComparison
}

type QueryFilterComparison uint8

func (q QueryFilterComparison) String() string {
	switch q {
	case QueryFilterComparisonLessThan:
		return "less than"
	case QueryFilterComparisonGreaterThan:
		return "greater than"
	case QueryFilterContains:
		return "contains"
	case QueryFilterFuzzy:
		return "fuzzily matches"
	}
	return "equals"
}

const (
	QueryFilterComparisonEquals = iota
	QueryFilterComparisonLessThan
	QueryFilterComparisonGreaterThan
	QueryFilterContains
	QueryFilterFuzzy
)
