// helper parser for dbox
package kpdbflex

import (
	"strings"

	"git.eaciitapp.com/sebar/dbflex"

	"github.com/raditzlawliet/kendoparser"
)

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
			// parser scope registered
			if opHandler, ok := OperatorManager.OperatorFilters[kf.Operator]; ok {
				return opHandler(ckFilter)
			}
		}

		// (global)
		if opHandler, ok := OperatorManager.OperatorFilters[kf.Operator]; ok {
			return opHandler(ckFilter)
		}

		// default (global)
		return OperatorManager.DefaultOperatorFilter(ckFilter)
	}

	// so filters has some values
	filters := []*dbflex.Filter{}
	for _, kFilterChild := range kf.Filters {
		filter := FilterParser(&kFilterChild)
		if filter != nil {
			filters = append(filters, filter.(*dbflex.Filter))
		}
	}

	if len(filters) > 0 {
		if strings.ToLower(kf.Logic) == "and" {
			return dbflex.And(filters...)
		}
		return dbflex.Or(filters...)
	}
	return nil // can return nil if filter & filters are meh ...
}

// ParserSort return []string
func ParserSort(ksa *kendoparser.Sort) interface{} {
	sorter := []string{}
	for _, ks := range *ksa {
		if strings.ToLower(ks.Dir) == "desc" {
			ks.Field = "-" + ks.Field
		}
		sorter = append(sorter, ks.Field)
	}
	return sorter
}
