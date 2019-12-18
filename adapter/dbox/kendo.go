// helper parser for dbox
package kpdbox

import (
	"strings"

	"github.com/eaciit/dbox"
	"github.com/raditzlawliet/gokendoparser"
)

// Parser Parser
type Parser struct{}

// ParseFilter convert KendoFilter into *dbox.Filter combination automaticly
// return can @Nullable if filter and filters empty
func (parser Parser) ParseFilter(kf *gokendoparser.KendoFilter) interface{} {
	// single filter
	if len(kf.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kf

		// very customable handler
		if kf.GetBeforeParse() != nil {
			for _, handler := range kf.GetBeforeParse() {
				if r := handler(&ckFilter); r != nil {
					return r
				}
			}
		}

		// local scope operator
		if kf.GetRegisteredOperators() != nil {
			if opHandler, ok := kf.GetRegisteredOperators()[kf.Operator]; ok && opHandler != nil {
				f := opHandler.Filter(ckFilter)
				if f != nil {
					return f.(*dbox.Filter)
				}
				return nil
			}
		}

		// parser scope registered
		if opHandler, ok := operatorManager.RegisteredOperators[kf.Operator]; ok {
			f := opHandler.Filter(ckFilter)
			if f != nil {
				return f.(*dbox.Filter)
			}
			return nil
		}

		// // global defult
		f := operatorManager.DefaultOperator.Filter(ckFilter)
		if f != nil {
			return f.(*dbox.Filter)
		}
		return nil
	}

	// so filters has some values
	dboxFilters := []*dbox.Filter{}
	for _, kFilterChild := range kf.Filters {
		dboxFilter := parser.ParseFilter(&kFilterChild)
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

// ParserSort return []string
func (parser Parser) ParserSort(ksa *gokendoparser.KendoSortArray) interface{} {
	sorter := []string{}
	for _, ks := range *ksa {
		if strings.ToLower(ks.Dir) == "desc" {
			ks.Field = "-" + ks.Field
		}
		sorter = append(sorter, ks.Field)
	}
	return sorter
}
