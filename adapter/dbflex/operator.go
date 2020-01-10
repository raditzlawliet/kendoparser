package kpdbflex

import (
	"time"

	"git.eaciitapp.com/sebar/dbflex"

	"github.com/raditzlawliet/gokendoparser"

	"github.com/spf13/cast"
)

var (
	operatorManager = gokendoparser.OperatorManager{}
	equalOp         = EqualOp{}
	notEqualOp      = NotEqualOp{}
	containOp       = ContainOp{}
	inOp            = InOp{}
	gtOp            = GtOp{}
	gteOp           = GteOp{}
	ltOp            = LtOp{}
	lteOp           = LteOp{}
	gtDateOp        = GtDateOp{}
	gteDateOp       = GteDateOp{}
	ltDateOp        = LtDateOp{}
	lteDateOp       = LteDateOp{}
	rangeOp         = RangeOp{}
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
	// operatorManager.RegisterOperator(notContainOp, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains", "notinclude", "notincludes", "doesnotinclude", "doesnotincludes") // not supported yet
	operatorManager.RegisterOperator(inOp, "in")
	operatorManager.RegisterOperator(gtOp, "gt")
	operatorManager.RegisterOperator(gteOp, "gte")
	operatorManager.RegisterOperator(ltOp, "lt")
	operatorManager.RegisterOperator(lteOp, "lte")
	operatorManager.RegisterOperator(gtDateOp, "gtdate")
	operatorManager.RegisterOperator(gteDateOp, "gtedate")
	operatorManager.RegisterOperator(ltDateOp, "ltdate")
	operatorManager.RegisterOperator(lteDateOp, "ltedate")
	operatorManager.RegisterOperator(rangeOp, "range", "between")
}

// EqualOp EqualOp
type EqualOp struct{}

// NotEqualOp NotEqualOp
type NotEqualOp struct{}

// ContainOp ContainOp
type ContainOp struct{}

// InOp InOp
type InOp struct{}

// GtOp GtOp
type GtOp struct{}

// GteOp GteOp
type GteOp struct{}

// LtOp LtOp
type LtOp struct{}

// LteOp LteOp
type LteOp struct{}

// GtDateOp GtDateOp
type GtDateOp struct{}

// GteDateOp GteDateOp
type GteDateOp struct{}

// LtDateOp LtDateOp
type LtDateOp struct{}

// LteDateOp LteDateOp
type LteDateOp struct{}

// RangeOp RangeOp
type RangeOp struct{}

// Filter Filter
func (EqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.Eq(kf.Field, kf.Value)
}

// Filter Filter
func (NotEqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.Ne(kf.Field, kf.Value)
}

// Filter Filter
func (ContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.Contains(kf.Field, kf.Value)
}

// Filter Filter
func (InOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.In(kf.Field, kf.Values...)
}

// Filter Filter
func (GtOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.Gt(kf.Field, kf.Value)
}

// Filter Filter
func (GteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.Gte(kf.Field, kf.Value)
}

// Filter Filter
func (LtOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.Lt(kf.Field, kf.Value)
}

// Filter Filter
func (LteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return dbflex.Lte(kf.Field, kf.Value)
}

// Filter Filter
func (GtDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Gt(kf.Field, dtVariable)
}

// Filter Filter
func (GteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Gte(kf.Field, dtVariable)
}

// Filter Filter
func (LtDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Lt(kf.Field, dtVariable)
}

// Filter Filter
func (LteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Lte(kf.Field, dtVariable)
}

// Filter Filter
func (RangeOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return dbflex.Range(kf.Field, v0, v1)
}
