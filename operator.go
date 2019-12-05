package gokendoparser

import (
	"sync"
)

// Operator basic interface of OperatorHander will have 2 of this func
type Operator interface {
	Filter(KendoFilter) interface{}
}

// OperatorManager OperatorManager
type OperatorManager struct {
	DefaultOperator     Operator
	RegisteredOperators map[string]Operator
	mutex               sync.Mutex
}

// RegisterOperator register operator with safe
func (om *OperatorManager) RegisterOperator(f Operator, ops ...string) {
	om.mutex.Lock()
	for _, op := range ops {
		if op != "" {
			if om.RegisteredOperators == nil {
				om.RegisteredOperators = map[string]Operator{}
			}
			om.RegisteredOperators[op] = f
		}
	}
	om.mutex.Unlock()
}

// Reset resetting global register (if needed)
func (om *OperatorManager) Reset() {
	om.mutex.Lock()
	om.RegisteredOperators = map[string]Operator{}
	om.mutex.Unlock()
}

//SetDefaultOperator by default, if no operator found, will use this instead
func (om *OperatorManager) SetDefaultOperator(f Operator) {
	om.mutex.Lock()
	om.DefaultOperator = f
	om.mutex.Unlock()
}

// RegisterTo RegisterTo
func (om *OperatorManager) RegisterTo(k *KendoRequest) {
	for op, fop := range om.RegisteredOperators {
		k.RegisterOperatorAll(fop, op)
	}
}
