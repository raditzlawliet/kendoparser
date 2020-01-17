package dbflex

// FilterOp is string represent enumeration of supported filter command
type FilterOp string

const (
	// OpAnd is AND
	OpAnd FilterOp = "$and"
	// OpOr is OR
	OpOr = "$or"
	// OpNot is Not
	OpNot = "$not"
	// OpEq is Equal
	OpEq = "$eq"
	// OpNe is Not Equal
	OpNe = "$ne"
	// OpGte is Greater than or Equal
	OpGte = "$gte"
	// OpGt is Greater than
	OpGt = "$gt"
	// OpLt is Less than
	OpLt = "$lt"
	// OpLte is Less than or equal
	OpLte = "$lte"
	// OpRange is range from until
	OpRange = "$range"
	// OpContains is Contains
	OpContains = "$contains"
	// OpStartWith is Start with
	OpStartWith = "$startwith"
	// OpEndWith is End with	
	OpEndWith = "$endwith"
	// OpIn is In
	OpIn = "$in"
	// OpNin is Not in
	OpNin = "$nin"
)

// Filter holding Items, Field, Operation, and Value
type Filter struct {
	Items []*Filter
	Field string
	Op    FilterOp
	Value interface{}
}

// NewFilter create new filter with given parameter
func NewFilter(field string, op FilterOp, v interface{}, items []*Filter) *Filter {
	f := new(Filter)
	f.Field = field
	f.Op = op
	f.Value = v
	if items != nil {
		f.Items = items
	}
	return f
}

// And create new filter with And operation
func And(items ...*Filter) *Filter {
	return NewFilter("", OpAnd, nil, items)
}

// Or create new filter with Or operation
func Or(items ...*Filter) *Filter {
	return NewFilter("", OpOr, nil, items)
}

// Eq create new filter with Eq operation
func Eq(field string, v interface{}) *Filter {
	return NewFilter(field, OpEq, v, nil)
}

// Not create new filter with Eq operation
func Not(item *Filter) *Filter {
	return NewFilter("", OpNot, nil, []*Filter{item})
}

// Ne create new filter with Ne operation
func Ne(field string, v interface{}) *Filter {
	return NewFilter(field, OpNe, v, nil)
}

// Gte create new filter with Gte operation
func Gte(field string, v interface{}) *Filter {
	return NewFilter(field, OpGte, v, nil)
}

// Gt create new filter with Gt operation
func Gt(field string, v interface{}) *Filter {
	return NewFilter(field, OpGt, v, nil)
}

// Lt create new filter with Lt operation
func Lt(field string, v interface{}) *Filter {
	return NewFilter(field, OpLt, v, nil)
}

// Lte create new filter with Lte operation
func Lte(field string, v interface{}) *Filter {
	return NewFilter(field, OpLte, v, nil)
}

// Range create new filter with Range operation
func Range(field string, from, to interface{}) *Filter {
	f := NewFilter(field, OpRange, nil, nil)
	f.Value = []interface{}{from, to}
	return f
}

// In create new filter with In operation
func In(field string, invalues ...interface{}) *Filter {
	f := new(Filter)
	f.Field = field
	f.Op = OpIn
	f.Value = invalues
	return f
}

// Nin create new filter with Nin operation
func Nin(field string, invalues ...interface{}) *Filter {
	f := new(Filter)
	f.Field = field
	f.Op = OpNin
	f.Value = invalues
	return f
}

// Contains create new filter with Contains operation
func Contains(field string, values ...string) *Filter {
	f := new(Filter)
	f.Field = field
	f.Op = OpContains
	f.Value = values
	return f
}

// StartWith create new filter with StartWith operation
func StartWith(field string, values string) *Filter {
	f := new(Filter)
	f.Field = field
	f.Op = OpStartWith
	f.Value = values
	return f
}

// EndWith create new filter with EndWith operation
func EndWith(field string, values string) *Filter {
	f := new(Filter)
	f.Field = field
	f.Op = OpEndWith
	f.Value = values
	return f
}
