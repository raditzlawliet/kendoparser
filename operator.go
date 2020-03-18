package kendoparser

import (
	"sync"
)

type OperatorFilter func(Filter) interface{}

// OperatorManager OperatorManager
type OperatorManager struct {
	DefaultOperatorFilter OperatorFilter
	OperatorFilters       map[string]OperatorFilter
	mutex                 sync.Mutex
}

// RegisterOperator register operator with safe
func (om *OperatorManager) RegisterOperator(f OperatorFilter, ops ...string) {
	om.mutex.Lock()
	for _, op := range ops {
		if op != "" {
			if om.OperatorFilters == nil {
				om.OperatorFilters = map[string]OperatorFilter{}
			}
			om.OperatorFilters[op] = f
		}
	}
	om.mutex.Unlock()
}

// Reset resetting global register (if needed)
func (om *OperatorManager) Reset() {
	om.mutex.Lock()
	om.OperatorFilters = map[string]OperatorFilter{}
	om.mutex.Unlock()
}

//SetDefaultOperator by default, if no operator found, will use this instead
func (om *OperatorManager) SetDefaultOperator(f OperatorFilter) {
	om.mutex.Lock()
	om.DefaultOperatorFilter = f
	om.mutex.Unlock()
}
