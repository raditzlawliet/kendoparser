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
	Filters    []KendoFilter `json:"filters"`
	Logic      string        `json:"logic"`
	Field      string        `json:"field"`
	Operator   string        `json:"operator"`
	IgnoreCase bool          `json:"ignoreCase"`
	Value      string        `json:"value"`
	Values     []interface{} `json:"values"`
}

// ToDboxFilter convert KendoFilter into *dbox.Filter combination automaticly
// return can @Nullable if filter and filters empty
func (kFilter *KendoFilter) ToDboxFilter() *dbox.Filter {
	// single filter
	if len(kFilter.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kFilter

		if opHandler, ok := RegisteredOperators[kFilter.Operator]; ok {
			return opHandler.ToDboxFilter(ckFilter)
		}
		return DefaultOperator.ToDboxFilter(ckFilter)
	}

	// so filters has some values
	dboxFilters := []*dbox.Filter{}
	for _, kFilterChild := range kFilter.Filters {
		dboxFilter := kFilterChild.ToDboxFilter()
		if dboxFilter != nil {
			dboxFilters = append(dboxFilters, dboxFilter)
		}
	}

	if len(dboxFilters) > 0 {
		if strings.ToLower(kFilter.Logic) == "and" {
			return dbox.And(dboxFilters...)
		}
		return dbox.Or(dboxFilters...)
	}

	return nil // can return nil if filter & filters are meh ...
}

// ToAggregationFilter convert KendoFilter into M for pipe combination automaticly
// return can @Nullable if filter and filters empty
func (kFilter *KendoFilter) ToAggregationFilter() toolkit.M {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
	if len(kFilter.Filters) == 0 {
		ckFilter := *kFilter
		if opHandler, ok := RegisteredOperators[kFilter.Operator]; ok {
			return opHandler.ToAggregationFilter(ckFilter)
		}
		return DefaultOperator.ToAggregationFilter(ckFilter)
	}

	// so filters has some values
	filters := []toolkit.M{}
	for _, kFilterChild := range kFilter.Filters {
		filter := kFilterChild.ToAggregationFilter()
		if filter != nil {
			filters = append(filters, filter)
		}
	}

	if len(filters) > 0 {
		if strings.ToLower(kFilter.Logic) == "and" {
			return toolkit.M{"$and": filters}
		}
		return toolkit.M{"$or": filters}
	}

	return nil
}
