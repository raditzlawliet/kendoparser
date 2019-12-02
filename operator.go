package gokendoparser

import (
	"sync"
)

// RegisteredOperators a global list registered operator, will use in all kendo, you can overwrite in struct scope if needed
var RegisteredOperators = map[string]Operator{}

// DefaultOperator will call if registerd not found, default operator empty
var DefaultOperator Operator = Op{}

var mutex = &sync.Mutex{}

// Operator basic interface of OperatorHander will have 2 of this func
type Operator interface {
	Filter(KendoFilter) interface{}
}

// RegisterOperator register operator with safe
func RegisterOperator(f Operator, ops ...string) {
	mutex.Lock()
	for _, op := range ops {
		RegisteredOperators[op] = f
	}
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

// Op empty op
type Op struct{}

// Filter Filter
func (Op) Filter(kf KendoFilter) interface{} {
	return nil
}
