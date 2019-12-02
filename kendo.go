package gokendoparser

/*
 * @Author
 * Radityo <radityohernanda@gmail.com>
 */

// KendoRequest option variable to struct (each apps has different format, defined/extend yourself if needed)
type KendoRequest struct {
	Data KendoData `json:"data"`
}

// RegisterOperatorAll register operator local scope include childs
func (kr *KendoRequest) RegisterOperatorAll(f Operator, op ...string) *KendoRequest {
	kr.Data.Filter.RegisterOperatorAll(f, op...)
	return kr
}

// KendoData datasource payload
type KendoData struct {
	Filter   KendoFilter    `json:"filter"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Skip     int            `json:"skip"`
	Take     int            `json:"take"`
	Sort     KendoSortArray `json:"sort"`
}

// KendoFilter struct filters
type KendoFilter struct {
	Filters    []KendoFilter `json:"filters"`
	Logic      string        `json:"logic"`
	Field      string        `json:"field"`
	Operator   string        `json:"operator"`
	IgnoreCase bool          `json:"ignoreCase"`
	Value      string        `json:"value"`
	Values     []interface{} `json:"values"`

	// for extension
	preFilter []func(*KendoFilter)

	// will not change the original value
	registeredOperators map[string]Operator
}

// GetRegisteredOperators GetRegisteredOperators
func (kf *KendoFilter) GetRegisteredOperators() map[string]Operator {
	return kf.registeredOperators
}

// AddRegisteredOperator AddRegisteredOperator
func (kf *KendoFilter) AddRegisteredOperator(k string, op Operator) *KendoFilter {
	if kf.registeredOperators == nil {
		kf.registeredOperators = map[string]Operator{}
	}
	if op != nil {
		kf.registeredOperators[k] = op
	}
	return kf
}

// GetPreFilter GetPreFilter
func (kf *KendoFilter) GetPreFilter() []func(*KendoFilter) {
	return kf.preFilter
}

// AddPreFilter AddPreFilter
func (kf *KendoFilter) AddPreFilter(f func(*KendoFilter)) *KendoFilter {
	if kf.preFilter == nil {
		kf.preFilter = []func(*KendoFilter){}
	}
	if f != nil {
		kf.preFilter = append(kf.preFilter, f)
	}
	return kf
}

// PreDboxFilterAll add custom handler pre-filtering apply to nested struct
func (kf *KendoFilter) AddPreFilterAll(f func(*KendoFilter)) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].AddPreFilterAll(f)
	}
	kf.AddPreFilter(f)
	return kf
}

// ResetPreFilter reset all pre-filter available
func (kf *KendoFilter) ResetPreFilter() *KendoFilter {
	kf.preFilter = []func(*KendoFilter){}
	return kf
}

// Transform your filter
func (kf *KendoFilter) Transform(t func(*KendoFilter)) *KendoFilter {
	t(kf)
	return kf
}

// TransformField only transform field
func (kf *KendoFilter) TransformField(t func(string) string) *KendoFilter {
	kf.Field = t(kf.Field)
	return kf
}

// TransformAll your filter include all childs
func (kf *KendoFilter) TransformAll(t func(*KendoFilter)) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].TransformAll(t)
	}
	kf.Transform(t)
	return kf
}

// TransformAllField only transform field include all childs
func (kf *KendoFilter) TransformAllField(t func(string) string) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].TransformAllField(t)
	}
	kf.TransformField(t)
	return kf
}

// RegisterOperator register operator local scope
func (kf *KendoFilter) RegisterOperator(f Operator, ops ...string) *KendoFilter {
	if kf.registeredOperators == nil {
		kf.registeredOperators = map[string]Operator{}
	}
	for _, op := range ops {
		kf.registeredOperators[op] = f
	}
	return kf
}

// RegisterOperatorAll register operator local scope include childs
func (kf *KendoFilter) RegisterOperatorAll(f Operator, ops ...string) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].RegisterOperatorAll(f, ops...)
	}
	kf.RegisterOperator(f, ops...)
	return kf
}

// KendoSort struct sort
type KendoSort struct {
	Dir   string `json:"dir"`
	Field string `json:"field"`
}

// KendoSortArray alias []KendoSort
type KendoSortArray []KendoSort