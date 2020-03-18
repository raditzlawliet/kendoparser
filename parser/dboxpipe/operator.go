package kpdboxpipe

import (
	"regexp"
	"strings"
	"time"

	"github.com/eaciit/toolkit"
	"github.com/raditzlawliet/kendoparser"
	"github.com/raditzlawliet/kendoparser/helper"
	"gopkg.in/mgo.v2/bson"

	"github.com/spf13/cast"
)

var (
	// OperatorManager of Mongo Parser
	OperatorManager = new(kendoparser.OperatorManager)
	// Operator bundle of Mongo Parser
	Operator = OperatorBundle{}
)

type OperatorBundle struct{}

func init() {
	RegisterOperator()
}

// RegisterOperator RegisterOperator
func RegisterOperator() {
	OperatorManager.SetDefaultOperator(Operator.Equal)
	OperatorManager.RegisterOperator(Operator.Equal, "eq", "equal")
	OperatorManager.RegisterOperator(Operator.NotEqual, "ne", "neq", "notequal")
	OperatorManager.RegisterOperator(Operator.Contain, "contain", "contains", "include", "includes")
	OperatorManager.RegisterOperator(Operator.NotContain, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains", "notinclude", "notincludes", "doesnotinclude", "doesnotincludes")
	OperatorManager.RegisterOperator(Operator.In, "in")
	OperatorManager.RegisterOperator(Operator.Gte, "gte")
	OperatorManager.RegisterOperator(Operator.Lte, "lte")
	OperatorManager.RegisterOperator(Operator.GteDate, "gtedate")
	OperatorManager.RegisterOperator(Operator.LteDate, "ltedate")
	OperatorManager.RegisterOperator(Operator.Exists, "exist", "exists")
	OperatorManager.RegisterOperator(Operator.Between, "between")
}

func (o *OperatorBundle) Equal(kf kendoparser.Filter) interface{} {
	if kf.IgnoreCase {
		value := regexp.QuoteMeta(cast.ToString(kf.Value))
		return toolkit.M{kf.Field: bson.RegEx{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
	}
	return toolkit.M{kf.Field: toolkit.M{"$eq": kf.Value}}
}

func (o *OperatorBundle) NotEqual(kf kendoparser.Filter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$ne": kf.Value}}
}

func (o *OperatorBundle) Contain(kf kendoparser.Filter) interface{} {
	return toolkit.M{kf.Field: helper.RegexContains(cast.ToString(kf.Value), kf.IgnoreCase)}
}

func (o *OperatorBundle) NotContain(kf kendoparser.Filter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$ne": helper.RegexContains(cast.ToString(kf.Value), kf.IgnoreCase)}}
}

func (o *OperatorBundle) In(kf kendoparser.Filter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$in": kf.Values}}
}

func (o *OperatorBundle) Gte(kf kendoparser.Filter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$gte": kf.Value}}
}

func (o *OperatorBundle) Lte(kf kendoparser.Filter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$lte": kf.Value}}
}

func (o *OperatorBundle) GteDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return toolkit.M{kf.Field: toolkit.M{"$gte": dtVariable}}
}

func (o *OperatorBundle) LteDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return toolkit.M{kf.Field: toolkit.M{"$lte": dtVariable}}
}

func (o *OperatorBundle) Exists(kf kendoparser.Filter) interface{} {
	return toolkit.M{kf.Field: toolkit.M{"$exists": helper.StringToBool(cast.ToString(kf.Value), false)}}
}

func (o *OperatorBundle) Between(kf kendoparser.Filter) interface{} {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return toolkit.M{kf.Field: toolkit.M{"$gte": v0, "$lte": v1}}
}
