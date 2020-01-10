package kpdbox

import (
	"regexp"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
	"github.com/raditzlawliet/gokendoparser"
	"github.com/raditzlawliet/gokendoparser/helper"
)

var (
	operatorManager = gokendoparser.OperatorManager{}
	equalOp         = EqualOp{}
	notEqualOp      = NotEqualOp{}
	containOp       = ContainOp{}
	notContainOp    = NotContainOp{}
	inOp            = InOp{}
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
	operatorManager.RegisterOperator(containOp, "contain", "contains", "include", "includes")
	operatorManager.RegisterOperator(notContainOp, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains", "notinclude", "notincludes", "doesnotinclude", "doesnotincludes")
	operatorManager.RegisterOperator(inOp, "in")
	operatorManager.RegisterOperator(gteOp, "gte")
	operatorManager.RegisterOperator(lteOp, "lte")
	operatorManager.RegisterOperator(gteDateOp, "gtedate")
	operatorManager.RegisterOperator(lteDateOp, "ltedate")
	operatorManager.RegisterOperator(existsOp, "exist", "exists")
	operatorManager.RegisterOperator(betweenOp, "between")
}

// EqualOp EqualOp
type EqualOp struct{}

// NotEqualOp NotEqualOp
type NotEqualOp struct{}

// ContainOp ContainOp
type ContainOp struct{}

// NotContainOp NotContainOp
type NotContainOp struct{}

// InOp InOp
type InOp struct{}

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
	return dbox.Eq(kf.Field, kf.Value)
}

// Filter Filter
func (NotEqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbox.Ne(kf.Field, kf.Value)
}

// Filter Filter
func (ContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbox.Contains(kf.Field, kf.Value)
}

// Filter Filter
func (NotContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	value := regexp.QuoteMeta(kf.Value)
	return &dbox.Filter{
		Field: kf.Field,
		Op:    dbox.FilterOpEqual, // equal are field = value and can be manipulate for others
		Value: toolkit.M{"$ne": toolkit.M{
			"$regex":   `` + value + ``,
			"$options": "i",
		}},
	}
}

// Filter Filter
func (InOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbox.In(kf.Field, kf.Values...)
}

// Filter Filter
func (GteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbox.Gte(kf.Field, kf.Value)
}

// Filter Filter
func (LteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbox.Lte(kf.Field, kf.Value)
}

// Filter Filter
func (GteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return dbox.Gte(kf.Field, dtVariable)
}

// Filter Filter
func (LteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return dbox.Lte(kf.Field, dtVariable)
}

// Filter Filter
func (ExistsOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return &dbox.Filter{
		Field: kf.Field,
		Op:    dbox.FilterOpEqual,
		Value: toolkit.M{
			"$exists": helper.StringToBool(kf.Value, false),
		},
	}
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
	return dbox.And(dbox.Gte(kf.Field, v0), dbox.Lte(kf.Field, v1))
}
