package kpdbox

import (
	"regexp"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
	"github.com/raditzlawliet/gokendoparser"
	"github.com/raditzlawliet/gokendoparser/helper"
)

func init() {
	RegisterOperator()
}

// RegisterOperator RegisterOperator
func RegisterOperator() {
	equalOp := EqualOp{}
	notEqualOp := NotEqualOp{}
	containOp := ContainOp{}
	notContainOp := NotContainOp{}
	inOp := InOp{}
	gteOp := GteOp{}
	lteOp := LteOp{}
	gteDateOp := GteDateOp{}
	lteDateOp := LteDateOp{}
	existsOp := ExistsOp{}

	gokendoparser.SetDefaultOperator(equalOp)
	gokendoparser.RegisterOperator(equalOp, "eq", "equal")
	gokendoparser.RegisterOperator(notEqualOp, "ne", "neq", "notequal")
	gokendoparser.RegisterOperator(containOp, "contain", "contains")
	gokendoparser.RegisterOperator(notContainOp, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains")
	gokendoparser.RegisterOperator(inOp, "in")
	gokendoparser.RegisterOperator(gteOp, "gte")
	gokendoparser.RegisterOperator(lteOp, "lte")
	gokendoparser.RegisterOperator(gteDateOp, "gteDate")
	gokendoparser.RegisterOperator(lteDateOp, "lteDate")
	gokendoparser.RegisterOperator(existsOp, "exists")
}

func RegisterTo(k *gokendoparser.KendoRequest) {
	equalOp := EqualOp{}
	notEqualOp := NotEqualOp{}
	containOp := ContainOp{}
	notContainOp := NotContainOp{}
	inOp := InOp{}
	gteOp := GteOp{}
	lteOp := LteOp{}
	gteDateOp := GteDateOp{}
	lteDateOp := LteDateOp{}
	existsOp := ExistsOp{}

	k.RegisterOperatorAll(equalOp, "eq", "equal")
	k.RegisterOperatorAll(notEqualOp, "ne", "neq", "notequal")
	k.RegisterOperatorAll(containOp, "contain", "contains")
	k.RegisterOperatorAll(notContainOp, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains")
	k.RegisterOperatorAll(inOp, "in")
	k.RegisterOperatorAll(gteOp, "gte")
	k.RegisterOperatorAll(lteOp, "lte")
	k.RegisterOperatorAll(gteDateOp, "gteDate")
	k.RegisterOperatorAll(lteDateOp, "lteDate")
	k.RegisterOperatorAll(existsOp, "exists")
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
