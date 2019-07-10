package gokendoparser

/*
 * @Author
 * Radityo <radityohernanda@gmail.com>
 */

import (
	"strings"

	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
)

// KendoRequest option variable to struct (each apps has different format, defined/extend yourself if needed)
type KendoRequest struct {
	Data KendoData `json:"data"`
}

// KendoData datasource payload
type KendoData struct {
	Filter   KendoFilter `json:"filter"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Skip     int         `json:"skip"`
	Take     int         `json:"take"`
	Sort     []KendoSort `json:"sort"`
}

// KendoSort struct sort
type KendoSort struct {
	Dir   string `json:"dir"`
	Field string `json:"field"`
}

// KendoFilter struct filters
type KendoFilter struct {
	Filters             []KendoFilter `json:"filters"`
	Logic               string        `json:"logic"`
	Field               string        `json:"field"`
	Operator            string        `json:"operator"`
	IgnoreCase          bool          `json:"ignoreCase"`
	Value               string        `json:"value"`
	Values              []interface{} `json:"values"`
	registeredOperators map[string]Operator
}

// ToDboxFilter convert KendoFilter into *dbox.Filter combination automaticly
// return can @Nullable if filter and filters empty
func (kf *KendoFilter) ToDboxFilter() *dbox.Filter {
	// single filter
	if len(kf.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kf

		// local scope operator
		if kf.registeredOperators != nil {
			if opHandler, ok := kf.registeredOperators[kf.Operator]; ok && opHandler != nil {
				return opHandler.ToDboxFilter(ckFilter)
			}
		}

		if opHandler, ok := RegisteredOperators[kf.Operator]; ok {
			return opHandler.ToDboxFilter(ckFilter)
		}

		return DefaultOperator.ToDboxFilter(ckFilter)
	}

	// so filters has some values
	dboxFilters := []*dbox.Filter{}
	for _, kFilterChild := range kf.Filters {
		dboxFilter := kFilterChild.ToDboxFilter()
		if dboxFilter != nil {
			dboxFilters = append(dboxFilters, dboxFilter)
		}
	}

	if len(dboxFilters) > 0 {
		if strings.ToLower(kf.Logic) == "and" {
			return dbox.And(dboxFilters...)
		}
		return dbox.Or(dboxFilters...)
	}

	return nil // can return nil if filter & filters are meh ...
}

// ToAggregationFilter convert KendoFilter into M for pipe combination automaticly
// return can @Nullable if filter and filters empty
func (kf *KendoFilter) ToAggregationFilter() toolkit.M {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
	if len(kf.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kf

		// local scope operator
		if kf.registeredOperators != nil {
			if opHandler, ok := kf.registeredOperators[kf.Operator]; ok && opHandler != nil {
				return opHandler.ToAggregationFilter(ckFilter)
			}
		}

		if opHandler, ok := RegisteredOperators[kf.Operator]; ok {
			return opHandler.ToAggregationFilter(ckFilter)
		}
		return DefaultOperator.ToAggregationFilter(ckFilter)
	}

	// so filters has some values
	filters := []toolkit.M{}
	for _, kFilterChild := range kf.Filters {
		filter := kFilterChild.ToAggregationFilter()
		if filter != nil {
			filters = append(filters, filter)
		}
	}

	if len(filters) > 0 {
		if strings.ToLower(kf.Logic) == "and" {
			return toolkit.M{"$and": filters}
		}
		return toolkit.M{"$or": filters}
	}

	return nil
}

// Transform your filter
func (kf *KendoFilter) Transform(t func(*KendoFilter)) {
	t(kf)
}

// TransformField only transform field
func (kf *KendoFilter) TransformField(t func(string) string) {
	kf.Field = t(kf.Field)
}

// TransformAll your filter include all childs
func (kf *KendoFilter) TransformAll(t func(*KendoFilter)) {
	for i := range kf.Filters {
		kf.Filters[i].TransformAll(t)
	}
	kf.Transform(t)
}

// TransformAllField only transform field include all childs
func (kf *KendoFilter) TransformAllField(t func(string) string) {
	for i := range kf.Filters {
		kf.Filters[i].TransformAllField(t)
	}
	kf.TransformField(t)
}

// RegisterOperator register operator local scope
func (kf *KendoFilter) RegisterOperator(op string, f Operator) {
	if kf.registeredOperators == nil {
		kf.registeredOperators = map[string]Operator{}
	}
	kf.registeredOperators[op] = f
}

// RegisterOperatorAll register operator local scope include childs
func (kf *KendoFilter) RegisterOperatorAll(op string, f Operator) {
	for i := range kf.Filters {
		kf.Filters[i].RegisterOperatorAll(op, f)
	}
	kf.RegisterOperator(op, f)
}
