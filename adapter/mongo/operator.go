package kpmongo

import (
	"regexp"
	"strings"
	"time"

	"github.com/raditzlawliet/gokendoparser"
	"github.com/raditzlawliet/gokendoparser/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/spf13/cast"
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

// BetweenOp BetweenOp
type BetweenOp struct{}

// Filter Filter
func (EqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	if kf.IgnoreCase {
		value := regexp.QuoteMeta(cast.ToString(kf.Value))
		return bson.M{kf.Field: primitive.Regex{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
	}
	return bson.M{kf.Field: bson.M{"$eq": kf.Value}}
}

// Filter Filter
func (NotEqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return bson.M{kf.Field: bson.M{"$ne": kf.Value}}
}

// Filter Filter
func (ContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return bson.M{kf.Field: RegexContains(cast.ToString(kf.Value), kf.IgnoreCase)}
}

// Filter Filter
func (NotContainOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return bson.M{kf.Field: bson.M{"$ne": RegexContains(cast.ToString(kf.Value), kf.IgnoreCase)}}
}

// Filter Filter
func (InOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return bson.M{kf.Field: bson.M{"$in": kf.Values}}
}

// Filter Filter
func (GteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return bson.M{kf.Field: bson.M{"$gte": kf.Value}}
}

// Filter Filter
func (LteOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return bson.M{kf.Field: bson.M{"$lte": kf.Value}}
}

// Filter Filter
func (GteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return bson.M{kf.Field: bson.M{"$gte": dtVariable}}
}

// Filter Filter
func (LteDateOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return bson.M{kf.Field: bson.M{"$lte": dtVariable}}
}

// Filter Filter
func (ExistsOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	return bson.M{kf.Field: bson.M{"$exists": helper.StringToBool(cast.ToString(kf.Value), false)}}
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
	return bson.M{kf.Field: bson.M{"$gte": v0, "$lte": v1}}
}

// RegexCaseInsensitive Generate bson.RegEx for case insensitive
func RegexCaseInsensitive(value string) primitive.Regex {
	value = regexp.QuoteMeta(value)
	return primitive.Regex{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}
}

// RegexContains Generate bson.RegEx for contains
func RegexContains(value string, ignoreCase bool) primitive.Regex {
	value = regexp.QuoteMeta(value)
	if ignoreCase {
		return primitive.Regex{Pattern: "" + strings.ToLower(value) + "", Options: "i"}
	} else {
		return primitive.Regex{Pattern: "" + strings.ToLower(value) + "", Options: ""}
	}
}
