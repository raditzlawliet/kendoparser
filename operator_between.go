package gokendoparser

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
)

// BetweenOperator sampling operator for between value 0 and value 1
type BetweenOperator struct {
}

func (BetweenOperator) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return dbox.And(dbox.Gte(kf.Field, v0), dbox.Lte(kf.Field, v1))
}
func (BetweenOperator) ToDboxPipe(kf KendoFilter) toolkit.M {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return toolkit.M{kf.Field: toolkit.M{"$gte": v0, "$lte": v1}}
}
