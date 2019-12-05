package gokendoparser

// Parser Parser
type Parser interface {
	ParseFilter(kf *KendoFilter) interface{}
	ParserSort(ksa *KendoSortArray) interface{}
}

// ParseFilter ParseFilter
type ParseFilter func(kf *KendoFilter) interface{}

// ParseSort ParseSort
type ParseSort func(ksa *KendoSortArray) interface{}
