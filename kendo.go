package kendoparser

/*
 * @Author
 * Radityo <radityohernanda@gmail.com>
 */

// KendoRequest option variable to struct (each apps has different format, defined/extend yourself if needed)
type KendoRequest struct {
	Data Data `json:"data"`
}

// Data Kendo DataSource payload
type Data struct {
	Filter   Filter `json:"filter"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Skip     int    `json:"skip"`
	Take     int    `json:"take"`
	Sort     Sort   `json:"sort"`
}

// Filter struct filters
type Filter struct {
	Filters    []Filter      `json:"filters"`
	Logic      string        `json:"logic"`
	Field      string        `json:"field"`
	Operator   string        `json:"operator"`
	IgnoreCase bool          `json:"ignoreCase"`
	Value      interface{}   `json:"value"`
	Values     []interface{} `json:"values"`

	// logic pre-parser
	additionalParsers []FilterParser
	operatorManager   *OperatorManager
}

// AdditionalParsers AdditionalParsers
func (f *Filter) AdditionalParsers() []FilterParser {
	return f.additionalParsers
}

// SetOperatorManager scoped Operator Manager
func (f *Filter) SetOperatorManager(om *OperatorManager) *Filter {
	f.operatorManager = om
	for i := range f.Filters {
		f.Filters[i].SetOperatorManager(om)
	}
	return f
}

// GetOperatorManager scoped Operator Manager
func (f *Filter) GetOperatorManager() *OperatorManager {
	return f.operatorManager
}

// AddParser AddParser
func (f *Filter) AddParser(parsers ...FilterParser) *Filter {
	if f.additionalParsers == nil {
		f.additionalParsers = []FilterParser{}
	}
	for _, parser := range parsers {
		if parser != nil {
			f.additionalParsers = append(f.additionalParsers, parser)
		}
	}
	return f
}

// AddAllParser AddAllParser
func (f *Filter) AddAllParser(parsers ...FilterParser) *Filter {
	for i := range f.Filters {
		f.Filters[i].AddAllParser(parsers...)
	}
	f.AddParser(parsers...)
	return f
}

// Parse Parse will return interface
func (f *Filter) Parse(parser FilterParser) interface{} {
	return parser(f)
}

// ResetAdditionalParsers reset all pre-filter available
func (f *Filter) ResetAdditionalParsers() *Filter {
	f.additionalParsers = []FilterParser{}
	return f
}

// ResetAllAdditionalParsers reset all pre-filter available
func (f *Filter) ResetAllAdditionalParsers() *Filter {
	for i := range f.Filters {
		f.Filters[i].ResetAllAdditionalParsers()
	}
	f.additionalParsers = []FilterParser{}
	return f
}

// Transform your filter
func (f *Filter) Transform(transform func(*Filter)) *Filter {
	transform(f)
	return f
}

// TransformField only transform field
func (f *Filter) TransformField(transform func(string) string) *Filter {
	f.Field = transform(f.Field)
	return f
}

// TransformAll your filter include all childs
func (f *Filter) TransformAll(transform func(*Filter)) *Filter {
	for i := range f.Filters {
		f.Filters[i].TransformAll(transform)
	}
	f.Transform(transform)
	return f
}

// TransformAllField only transform field include all childs
func (f *Filter) TransformAllField(transform func(string) string) *Filter {
	for i := range f.Filters {
		f.Filters[i].TransformAllField(transform)
	}
	f.TransformField(transform)
	return f
}

// SortDetail struct sort
type SortDetail struct {
	Dir   string `json:"dir"`
	Field string `json:"field"`
}

// Sort alias []SortDetail
type Sort []SortDetail

// Parse Parse
func (s *Sort) Parse(parser SortParser) interface{} {
	return parser(s)
}
