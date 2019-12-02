package kpdboxpipe

import (
	"regexp"
	"strings"
	"time"

	"github.com/eaciit/toolkit"
	"github.com/raditzlawliet/gokendoparser"
	"github.com/raditzlawliet/gokendoparser/helper"
	"gopkg.in/mgo.v2/bson"
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
	gokendoparser.RegisterOperator(gteDateOp, "gtedate")
	gokendoparser.RegisterOperator(lteDateOp, "ltedate")
	gokendoparser.RegisterOperator(existsOp, "exist", "exists")
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

	k.SetDefaultOperator(equalOp)
	k.RegisterOperatorAll(equalOp, "eq", "equal")
	k.RegisterOperatorAll(notEqualOp, "ne", "neq", "notequal")
	k.RegisterOperatorAll(containOp, "contain", "contains")
	k.RegisterOperatorAll(notContainOp, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains")
	k.RegisterOperatorAll(inOp, "in")
	k.RegisterOperatorAll(gteOp, "gte")
	k.RegisterOperatorAll(lteOp, "lte")
	k.RegisterOperatorAll(gteDateOp, "gtedate")
	k.RegisterOperatorAll(lteDateOp, "ltedate")
	k.RegisterOperatorAll(existsOp, "exist", "exists")
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
	if kf.IgnoreCase {
		value := regexp.QuoteMeta(kf.Value)
		return toolkit.M{kf.Field: bson.RegEx{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
	}
	return toolkit.M{kf.Field: toolkit.M{"$eq": kf.Value}}
}

// Filter Filter
func (NotEqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$ne": kf.Value}}
}

// Filter Filter
func (ContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return toolkit.M{kf.Field: helper.RegexContains(kf.Value, kf.IgnoreCase)}
}

// Filter Filter
func (NotContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$ne": helper.RegexContains(kf.Value, kf.IgnoreCase)}}
}

// Filter Filter
func (InOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$in": kf.Values}}
}

// Filter Filter
func (GteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$gte": kf.Value}}
}

// Filter Filter
func (LteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$lte": kf.Value}}
}

// Filter Filter
func (GteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return toolkit.M{kf.Field: toolkit.M{"$gte": dtVariable}}
}

// Filter Filter
func (LteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return toolkit.M{kf.Field: toolkit.M{"$lte": dtVariable}}
}

// Filter Filter
func (ExistsOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$exists": helper.StringToBool(kf.Value, false)}}
}
