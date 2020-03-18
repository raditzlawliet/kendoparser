// helper parser for dbox
package kpdbox

import (
	"strings"

	"github.com/eaciit/dbox"
	"github.com/raditzlawliet/kendoparser"
)

// Parser Parser
// type Parser struct{}

// FilterParser convert Filter into *dbox.Filter combination automaticly
// return can @Nullable if filter and filters empty
func FilterParser(kf *kendoparser.Filter) interface{} {
	// single filter
	if len(kf.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kf

		// very customable handler
		if kf.AdditionalParsers() != nil {
			for _, handler := range kf.AdditionalParsers() {
				if r := handler(&ckFilter); r != nil {
					return r
				}
			}
		}

		// (scoped)
		if om := kf.GetOperatorManager(); om != nil {
			if opHandler, ok := om.OperatorFilters[kf.Operator]; ok {
				f := opHandler(ckFilter)
				if f != nil {
					return f.(*dbox.Filter)
				}
				return nil
			}
		}

		// (global)
		if opHandler, ok := OperatorManager.OperatorFilters[kf.Operator]; ok {
			f := opHandler(ckFilter)
			if f != nil {
				return f.(*dbox.Filter)
			}
			return nil
		}

		// default (global)
		f := OperatorManager.DefaultOperatorFilter(ckFilter)
		if f != nil {
			return f.(*dbox.Filter)
		}
		return nil
	}

	// so filters has some values
	dboxFilters := []*dbox.Filter{}
	for _, kFilterChild := range kf.Filters {
		dboxFilter := FilterParser(&kFilterChild)
		if dboxFilter != nil {
			dboxFilters = append(dboxFilters, dboxFilter.(*dbox.Filter))
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

// SortParser return []string
func SortParser(ksa *kendoparser.Sort) interface{} {
	sorter := []string{}
	for _, ks := range *ksa {
		if strings.ToLower(ks.Dir) == "desc" {
			ks.Field = "-" + ks.Field
		}
		sorter = append(sorter, ks.Field)
	}
	return sorter
}
