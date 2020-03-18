package kendoparser

// FilterParser FilterParser
type FilterParser func(f *Filter) interface{}

// SortParser SortParser
type SortParser func(s *Sort) interface{}
