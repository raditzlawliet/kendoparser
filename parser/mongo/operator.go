package kpmongo

import (
	"regexp"
	"strings"
	"time"

	"github.com/raditzlawliet/kendoparser"
	"github.com/raditzlawliet/kendoparser/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	OperatorManager.RegisterOperator(Operator.Equal, "eq", "equal", "=", "==")
	OperatorManager.RegisterOperator(Operator.NotEqual, "ne", "neq", "notequal", "<>", "!=")
	OperatorManager.RegisterOperator(Operator.EqualDate, "eqdate", "equaldate", "=date", "==date")
	OperatorManager.RegisterOperator(Operator.NotEqualDate, "nedate", "neqdate", "notequaldate", "<>date", "!=date")
	OperatorManager.RegisterOperator(Operator.Contain, "contain", "contains", "include", "includes")
	OperatorManager.RegisterOperator(Operator.NotContain, "notcontains", "notcontains", "doesnotcontain", "doesnotcontains", "notinclude", "notincludes", "doesnotinclude", "doesnotincludes")
	OperatorManager.RegisterOperator(Operator.In, "in")
	OperatorManager.RegisterOperator(Operator.NotIn, "notin", "nin")
	OperatorManager.RegisterOperator(Operator.Gt, "gt", ">")
	OperatorManager.RegisterOperator(Operator.Lt, "lt", "<")
	OperatorManager.RegisterOperator(Operator.Gte, "gte", ">=")
	OperatorManager.RegisterOperator(Operator.Lte, "lte", "<=")
	OperatorManager.RegisterOperator(Operator.GtDate, "gtdate", ">date")
	OperatorManager.RegisterOperator(Operator.LtDate, "ltdate", "<date")
	OperatorManager.RegisterOperator(Operator.GteDate, "gtedate", ">=date")
	OperatorManager.RegisterOperator(Operator.LteDate, "ltedate", "<=date")
	OperatorManager.RegisterOperator(Operator.Exists, "exist", "exists")
	OperatorManager.RegisterOperator(Operator.DateExists, "dateexist", "dateexists")
	OperatorManager.RegisterOperator(Operator.Between, "between")
}

func (o *OperatorBundle) Equal(kf kendoparser.Filter) interface{} {
	// Ignore-case only for string :D
	if value, ok := kf.Value.(string); ok && kf.IgnoreCase {
		value := regexp.QuoteMeta(value)
		return bson.M{kf.Field: primitive.Regex{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
	}
	// other string, it will directly :|
	return bson.M{kf.Field: bson.M{"$eq": kf.Value}}
}

func (o *OperatorBundle) NotEqual(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$ne": kf.Value}}
}

func (o *OperatorBundle) EqualDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return bson.M{kf.Field: bson.M{"$eq": dtVariable}}
}

func (o *OperatorBundle) NotEqualDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return bson.M{kf.Field: bson.M{"$ne": dtVariable}}
}

func (o *OperatorBundle) Contain(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: RegexContains(cast.ToString(kf.Value), kf.IgnoreCase)}
}

func (o *OperatorBundle) ContainNumber(kf kendoparser.Filter) interface{} {
	field := kf.Field
	if !strings.HasPrefix(field, "$") {
		field = "$" + field
	}
	return bson.M{
		"$expr": bson.M{
			"$regexMatch": bson.M{
				"input": bson.M{"$toString": field},
				"regex": RegexContains(cast.ToString(kf.Value), kf.IgnoreCase),
			},
		},
	}
}

func (o *OperatorBundle) NotContain(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$ne": RegexContains(cast.ToString(kf.Value), kf.IgnoreCase)}}
}

func (o *OperatorBundle) In(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$in": kf.Values}}
}

func (o *OperatorBundle) NotIn(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$nin": kf.Values}}
}

func (o *OperatorBundle) Gt(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$gt": kf.Value}}
}

func (o *OperatorBundle) Lt(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$lt": kf.Value}}
}

func (o *OperatorBundle) Gte(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$gte": kf.Value}}
}

func (o *OperatorBundle) Lte(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$lte": kf.Value}}
}

func (o *OperatorBundle) GtDate(kf kendoparser.Filter) interface{} {
	str := cast.ToString(kf.Value)
	if t, ok := kf.Value.(time.Time); ok { // fixing time.Time value not valid converted :|
		str = t.Format(time.RFC3339)
	}
	dtVariable, _ := time.Parse(time.RFC3339, str)
	return bson.M{kf.Field: bson.M{"$gt": dtVariable}}
}

func (o *OperatorBundle) LtDate(kf kendoparser.Filter) interface{} {
	str := cast.ToString(kf.Value)
	if t, ok := kf.Value.(time.Time); ok { // fixing time.Time value not valid converted :|
		str = t.Format(time.RFC3339)
	}
	dtVariable, _ := time.Parse(time.RFC3339, str)
	return bson.M{kf.Field: bson.M{"$lt": dtVariable}}
}

func (o *OperatorBundle) GteDate(kf kendoparser.Filter) interface{} {
	str := cast.ToString(kf.Value)
	if t, ok := kf.Value.(time.Time); ok { // fixing time.Time value not valid converted :|
		str = t.Format(time.RFC3339)
	}
	dtVariable, _ := time.Parse(time.RFC3339, str)
	return bson.M{kf.Field: bson.M{"$gte": dtVariable}}
}

func (o *OperatorBundle) LteDate(kf kendoparser.Filter) interface{} {
	str := cast.ToString(kf.Value)
	if t, ok := kf.Value.(time.Time); ok { // fixing time.Time value not valid converted :|
		str = t.Format(time.RFC3339)
	}
	dtVariable, _ := time.Parse(time.RFC3339, str)
	return bson.M{kf.Field: bson.M{"$lte": dtVariable}}
}

func (o *OperatorBundle) Exists(kf kendoparser.Filter) interface{} {
	return bson.M{kf.Field: bson.M{"$exists": helper.StringToBool(cast.ToString(kf.Value), false)}}
}

// filter for time.Time / DateTime in Mongo exists or not (usefull for safe-delete flag)
func (o *OperatorBundle) DateExists(kf kendoparser.Filter) interface{} {
	if strings.ToLower(cast.ToString(kf.Value)) == "true" {
		return bson.M{"$and": bson.A{
			bson.M{kf.Field: bson.M{"$exists": true}},
			bson.M{kf.Field: bson.M{"$ne": time.Time{}}},
		}}
	}
	return bson.M{"$or": bson.A{
		bson.M{kf.Field: bson.M{"$exists": false}},
		bson.M{kf.Field: bson.M{"$eq": time.Time{}}},
	}}
}

func (o *OperatorBundle) Between(kf kendoparser.Filter) interface{} {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return bson.M{kf.Field: bson.M{"$gte": v0, "$lte": v1}}
}

// RegexContains Generate bson.RegEx for contains
func RegexContains(value string, ignoreCase bool) primitive.Regex {
	value = regexp.QuoteMeta(value)
	if ignoreCase {
		return primitive.Regex{Pattern: "" + strings.ToLower(value) + "", Options: "i"}
	} else {
		return primitive.Regex{Pattern: "" + value + "", Options: ""}
	}
}
