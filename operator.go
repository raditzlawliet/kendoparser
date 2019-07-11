package gokendoparser

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

// RegisteredOperators a global list registered operator, will use in all kendo, you can overwrite in struct scope if needed
var RegisteredOperators = map[string]Operator{}

// DefaultOperator will call if registerd not found, default operator is equa
var DefaultOperator Operator = EqualOp{}

var mutex = &sync.Mutex{}

// Operator basic interface of OperatorHander will have 2 of this func
type Operator interface {
	ToDboxFilter(KendoFilter) *dbox.Filter
	ToDboxPipe(KendoFilter) toolkit.M
}

func init() {
	// you can assign as easy like this
	eqOp := EqualOp{}
	RegisterOperator("eq", eqOp)
	RegisterOperator("equal", eqOp)

	neqOp := NotEqualOp{}
	RegisterOperator("ne", neqOp)
	RegisterOperator("neq", neqOp)
	RegisterOperator("notequal", neqOp)

	containOp := ContainOp{}
	RegisterOperator("contain", containOp)
	RegisterOperator("contains", containOp)

	notcontainOp := NotContainOp{}
	RegisterOperator("notcontain", notcontainOp)
	RegisterOperator("notcontains", notcontainOp)
	RegisterOperator("doesnotcontain", notcontainOp)
	RegisterOperator("doesnotcontains", notcontainOp)

	inOp := InOp{}
	RegisterOperator("in", inOp)

	gteOp := GteOp{}
	RegisterOperator("gte", gteOp)

	lteOp := LteOp{}
	RegisterOperator("lte", lteOp)

	gtedateOp := GteDateOp{}
	RegisterOperator("gtedate", gtedateOp)

	ltedateOp := LteDateOp{}
	RegisterOperator("ltedate", ltedateOp)

	existsOp := ExistsOp{}
	RegisterOperator("exists", existsOp)
}

// RegisterOperator register operator with safe
func RegisterOperator(op string, f Operator) {
	mutex.Lock()
	RegisteredOperators[op] = f
	mutex.Unlock()
}

// ResetRegisterOperator resetting global register (if needed)
func ResetRegisterOperator() {
	mutex.Lock()
	RegisteredOperators = map[string]Operator{}
	mutex.Unlock()
}

//SetDefaultOperator by default, if no operator found, will use this instead
func SetDefaultOperator(f Operator) {
	mutex.Lock()
	DefaultOperator = f
	mutex.Unlock()
}

/*
 * basic operator available
 */
type EqualOp struct{}
type NotEqualOp struct{}
type ContainOp struct{}
type NotContainOp struct{}
type InOp struct{}
type GteOp struct{}
type LteOp struct{}
type GteDateOp struct{}
type LteDateOp struct{}
type ExistsOp struct{}

func (EqualOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	return dbox.Eq(kf.Field, kf.Value)
}
func (EqualOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	if kf.IgnoreCase {
		value := regexp.QuoteMeta(kf.Value)
		return toolkit.M{kf.Field: bson.RegEx{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
	}
	return toolkit.M{kf.Field: toolkit.M{"$eq": kf.Value}}
}
func (NotEqualOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	return dbox.Ne(kf.Field, kf.Value)
}
func (NotEqualOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	return toolkit.M{kf.Field: toolkit.M{"$ne": kf.Value}}
}
func (ContainOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	return dbox.Contains(kf.Field, kf.Value)
}
func (ContainOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	return toolkit.M{kf.Field: RegexContains(kf.Value, kf.IgnoreCase)}
}
func (NotContainOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
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
func (NotContainOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	return toolkit.M{kf.Field: toolkit.M{"$ne": RegexContains(kf.Value, kf.IgnoreCase)}}
}
func (InOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	return dbox.In(kf.Field, kf.Values...)
}
func (InOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	return toolkit.M{kf.Field: toolkit.M{"$in": kf.Values}}
}
func (GteOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	return dbox.Gte(kf.Field, kf.Value)
}
func (GteOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	return toolkit.M{kf.Field: toolkit.M{"$gte": kf.Value}}
}
func (LteOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	return dbox.Lte(kf.Field, kf.Value)
}
func (LteOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	return toolkit.M{kf.Field: toolkit.M{"$lte": kf.Value}}
}
func (GteDateOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return dbox.Gte(kf.Field, dtVariable)
}
func (GteDateOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return toolkit.M{kf.Field: toolkit.M{"$gte": dtVariable}}
}
func (LteDateOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return dbox.Lte(kf.Field, dtVariable)
}
func (LteDateOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	dtVariable, _ := time.Parse(time.RFC3339, kf.Value)
	return toolkit.M{kf.Field: toolkit.M{"$lte": dtVariable}}
}
func (ExistsOp) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	return &dbox.Filter{
		Field: kf.Field,
		Op:    dbox.FilterOpEqual,
		Value: toolkit.M{
			"$exists": StringToBool(kf.Value, false),
		},
	}
}
func (ExistsOp) ToDboxPipe(kf KendoFilter) toolkit.M {
	return toolkit.M{kf.Field: toolkit.M{"$exists": StringToBool(kf.Value, false)}}
}
