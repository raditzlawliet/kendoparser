package kpdbflex

import (
	"time"

	"git.eaciitapp.com/sebar/dbflex"

	"github.com/raditzlawliet/kendoparser"

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
	OperatorManager.RegisterOperator(Operator.In, "in")
	OperatorManager.RegisterOperator(Operator.Gt, "gt")
	OperatorManager.RegisterOperator(Operator.Gte, "gte")
	OperatorManager.RegisterOperator(Operator.Lt, "lt")
	OperatorManager.RegisterOperator(Operator.Lte, "lte")
	OperatorManager.RegisterOperator(Operator.GtDate, "gtdate")
	OperatorManager.RegisterOperator(Operator.GteDate, "gtedate")
	OperatorManager.RegisterOperator(Operator.LtDate, "ltdate")
	OperatorManager.RegisterOperator(Operator.LteDate, "ltedate")
	OperatorManager.RegisterOperator(Operator.Range, "range", "between")
}

func (o *OperatorBundle) Equal(kf kendoparser.Filter) interface{} {
	return dbflex.Eq(kf.Field, kf.Value)
}

func (o *OperatorBundle) NotEqual(kf kendoparser.Filter) interface{} {
	return dbflex.Ne(kf.Field, kf.Value)
}

func (o *OperatorBundle) Contain(kf kendoparser.Filter) interface{} {
	return dbflex.Contains(kf.Field, cast.ToString(kf.Value))
}

func (o *OperatorBundle) In(kf kendoparser.Filter) interface{} {
	return dbflex.In(kf.Field, kf.Values...)
}

func (o *OperatorBundle) Gt(kf kendoparser.Filter) interface{} {
	return dbflex.Gt(kf.Field, kf.Value)
}

func (o *OperatorBundle) Gte(kf kendoparser.Filter) interface{} {
	return dbflex.Gte(kf.Field, kf.Value)
}

func (o *OperatorBundle) Lt(kf kendoparser.Filter) interface{} {
	return dbflex.Lt(kf.Field, kf.Value)
}

func (o *OperatorBundle) Lte(kf kendoparser.Filter) interface{} {
	return dbflex.Lte(kf.Field, kf.Value)
}

func (o *OperatorBundle) GtDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Gt(kf.Field, dtVariable)
}

func (o *OperatorBundle) GteDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Gte(kf.Field, dtVariable)
}

func (o *OperatorBundle) LtDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Lt(kf.Field, dtVariable)
}

func (o *OperatorBundle) LteDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return dbflex.Lte(kf.Field, dtVariable)
}

func (o *OperatorBundle) Range(kf kendoparser.Filter) interface{} {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return dbflex.Range(kf.Field, v0, v1)
}
