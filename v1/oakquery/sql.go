package oakquery

import (
	"fmt"
	"strconv"
)

func (qr *QueryRange) AsSQL() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", qr.PerPage, qr.PerPage*qr.Page)
}

func (qf *QueryFilter) AsStringSQL(prefix string) string {
	field := QuoteIdentifier(prefix + qf.Field)
	value := Quote(qf.Value)
	// switch qf.Comparison {
	// case QueryFilterComparisonLessThan:
	// 	return fmt.Sprintf("%s < %d", field, number)
	// case QueryFilterComparisonGreaterThan:
	// 	return fmt.Sprintf("%s > %d", field, number)
	// default:
	// }
	return field + "=" + value
}

func (qf *QueryFilter) AsNumericSQL(prefix string) string {
	field := QuoteIdentifier(prefix + qf.Field)
	number, _ := strconv.Atoi(qf.Value) // TODO: eats floats?
	switch qf.Comparison {
	case QueryFilterComparisonLessThan:
		return fmt.Sprintf("%s<%d", field, number)
	case QueryFilterComparisonGreaterThan:
		return fmt.Sprintf("%s>%d", field, number)
	default:
		return fmt.Sprintf("%s=%d", field, number)
	}
}
