package xorm

import (
	"time"

	"github.com/raditzlawliet/gokendoparser"
	"github.com/raditzlawliet/gokendoparser/helper"
	"xorm.io/builder"
)

var (
	operatorManager = gokendoparser.OperatorManager{}
	equalOp         = EqualOp{}
	notEqualOp      = NotEqualOp{}
	containOp       = ContainOp{}
	notContainOp    = NotContainOp{}
	inOp            = InOp{}
	gtOp            = GtOp{}
	ltOp            = LtOp{}
	gteOp           = GteOp{}
	lteOp           = LteOp{}
	gteDateOp       = GteDateOp{}
	lteDateOp       = LteDateOp{}
	existsOp        = ExistsOp{}
	betweenOp       = BetweenOp{}
)

func init() {
	RegisterOperator()
}

// RegisterOperator RegisterOperator
func RegisterOperator() {
	operatorManager.SetDefaultOperator(equalOp)
	operatorManager.RegisterOperator(equalOp, "eq", "equal")
	operatorManager.RegisterOperator(notEqualOp, "ne", "neq", "notequal")
	operatorManager.RegisterOperator(containOp, "contain", "contains")
	operatorManager.RegisterOperator(notContainOp, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains")
	operatorManager.RegisterOperator(inOp, "in")
	operatorManager.RegisterOperator(gtOp, "gt")
	operatorManager.RegisterOperator(ltOp, "lt")
	operatorManager.RegisterOperator(gteOp, "gte")
	operatorManager.RegisterOperator(lteOp, "lte")
	operatorManager.RegisterOperator(gteDateOp, "gtedate")
	operatorManager.RegisterOperator(lteDateOp, "ltedate")
	operatorManager.RegisterOperator(existsOp, "exist", "exists")
	operatorManager.RegisterOperator(betweenOp, "between")
}

//EqualOp EqualOp
type EqualOp struct{}

// NotEqualOp NotEqualOp
type NotEqualOp struct{}

// ContainOp ContainOp
type ContainOp struct{}

// NotContainOp NotContainOp
type NotContainOp struct{}

// InOp InOp
type InOp struct{}

// GtOp GtOp
type GtOp struct{}

// LtOp LtOp
type LtOp struct{}

// GteOp GteOp
type GteOp struct{}

// LteOp LteOp
type LteOp struct{}

// GteDateOp GteDateOp
type GteDateOp struct{}

// LteDateOp LteDateOp
type LteDateOp struct{}

// ExistsOp ExistsOp
type ExistsOp struct{}

// BetweenOp BetweenOp
type BetweenOp struct{}

// Filter Filter
func (EqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	if kf.IgnoreCase {
		return EqCi{kf.Field: kf.Value}
	}
	return builder.Eq{kf.Field: kf.Value}
}

// Filter Filter
func (NotEqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return builder.Neq{kf.Field: kf.Value}
}

// Filter Filter
func (ContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	if kf.IgnoreCase {
		return LikeCi{kf.Field, kf.Value}
	}
	return builder.Like{kf.Value, kf.Value}
}

// Filter Filter
func (NotContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	if kf.IgnoreCase {
		return builder.Not{
			LikeCi{kf.Field, kf.Value},
		}
	}
	return builder.Not{
		builder.Like{kf.Value, kf.Value},
	}
}

// Filter Filter
func (InOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return builder.In(kf.Field, kf.Values...)
}

// Filter Filter
func (GtOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return builder.Gt{kf.Field: kf.Value}
}

// Filter Filter
func (LtOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return builder.Lt{kf.Field: kf.Value}
}

// Filter Filter
func (GteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return builder.Gte{kf.Field: kf.Value}
}

// Filter Filter
func (LteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return builder.Lte{kf.Field: kf.Value}
}

// Filter Filter
func (GteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return builder.Gte{kf.Field: dtVariable}
}

// Filter Filter
func (LteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return builder.Lte{kf.Field: dtVariable}
}

// Filter Filter
func (ExistsOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	if helper.StringToBool(kf.Value, false) {
		return builder.NotNull{kf.Field}
	}
	return builder.IsNull{kf.Field}
}

// Filter Filter
func (BetweenOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return builder.Between{kf.Field, v0, v1}
}
