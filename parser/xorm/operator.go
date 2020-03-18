package kpxorm

import (
	"time"

	"github.com/raditzlawliet/kendoparser"
	"github.com/raditzlawliet/kendoparser/helper"
	"github.com/spf13/cast"
	"xorm.io/builder"
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
	OperatorManager.RegisterOperator(Operator.Gt, "gt")
	OperatorManager.RegisterOperator(Operator.Lt, "lt")
	OperatorManager.RegisterOperator(Operator.Gte, "gte")
	OperatorManager.RegisterOperator(Operator.Lte, "lte")
	OperatorManager.RegisterOperator(Operator.GteDate, "gtedate")
	OperatorManager.RegisterOperator(Operator.LteDate, "ltedate")
	OperatorManager.RegisterOperator(Operator.Exists, "exist", "exists")
	OperatorManager.RegisterOperator(Operator.Between, "between")
}

func (o *OperatorBundle) Equal(kf kendoparser.Filter) interface{} {
	if kf.IgnoreCase {
		return EqCi{kf.Field: kf.Value}
	}
	return builder.Eq{kf.Field: kf.Value}
}

func (o *OperatorBundle) NotEqual(kf kendoparser.Filter) interface{} {
	return builder.Neq{kf.Field: kf.Value}
}

func (o *OperatorBundle) Contain(kf kendoparser.Filter) interface{} {
	if kf.IgnoreCase {
		return LikeCi{kf.Field, cast.ToString(kf.Value)}
	}
	return builder.Like{cast.ToString(kf.Value), cast.ToString(kf.Value)}
}

func (o *OperatorBundle) NotContain(kf kendoparser.Filter) interface{} {
	if kf.IgnoreCase {
		return builder.Not{
			LikeCi{kf.Field, cast.ToString(kf.Value)},
		}
	}
	return builder.Not{
		builder.Like{cast.ToString(kf.Value), cast.ToString(kf.Value)},
	}
}

func (o *OperatorBundle) In(kf kendoparser.Filter) interface{} {
	return builder.In(kf.Field, kf.Values...)
}

func (o *OperatorBundle) Gt(kf kendoparser.Filter) interface{} {
	return builder.Gt{kf.Field: kf.Value}
}

func (o *OperatorBundle) Lt(kf kendoparser.Filter) interface{} {
	return builder.Lt{kf.Field: kf.Value}
}

func (o *OperatorBundle) Gte(kf kendoparser.Filter) interface{} {
	return builder.Gte{kf.Field: kf.Value}
}

func (o *OperatorBundle) Lte(kf kendoparser.Filter) interface{} {
	return builder.Lte{kf.Field: kf.Value}
}

func (o *OperatorBundle) GteDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return builder.Gte{kf.Field: dtVariable}
}

func (o *OperatorBundle) LteDate(kf kendoparser.Filter) interface{} {
	dtVariable, _ := time.Parse(time.RFC3339, cast.ToString(kf.Value))
	return builder.Lte{kf.Field: dtVariable}
}

func (o *OperatorBundle) Exists(kf kendoparser.Filter) interface{} {
	if helper.StringToBool(cast.ToString(kf.Value), false) {
		return builder.NotNull{kf.Field}
	}
	return builder.IsNull{kf.Field}
}

func (o *OperatorBundle) Between(kf kendoparser.Filter) interface{} {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return builder.Between{kf.Field, v0, v1}
}
