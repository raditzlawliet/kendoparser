// helper parser for dbox
package kpdbox

import (
	"strings"

	"github.com/eaciit/dbox"
	"github.com/raditzlawliet/gokendoparser"
)

// ParseFilter convert KendoFilter into *dbox.Filter combination automaticly
// return can @Nullable if filter and filters empty
func ParseFilter(kf *gokendoparser.KendoFilter) *dbox.Filter {
	// single filter
	if len(kf.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kf

		// very customable handler
		if kf.GetPreFilter() != nil {
			for _, handler := range kf.GetPreFilter() {
				handler(&ckFilter)
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

		if opHandler, ok := gokendoparser.RegisteredOperators[kf.Operator]; ok {
			f := opHandler.Filter(ckFilter)
			if f != nil {
				return f.(*dbox.Filter)
			}
			return nil
		}

		f := gokendoparser.DefaultOperator.Filter(ckFilter)
		if f != nil {
			return f.(*dbox.Filter)
		}
		return nil
	}

	// so filters has some values
	dboxFilters := []*dbox.Filter{}
	for _, kFilterChild := range kf.Filters {
		dboxFilter := ParseFilter(&kFilterChild)
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

// ParserSort ParserSort
func ParserSort(ksa *gokendoparser.KendoSortArray) []string {
	sorter := []string{}
	for _, ks := range *ksa {
		if strings.ToLower(ks.Dir) == "desc" {
			ks.Field = "-" + ks.Field
		}
		sorter = append(sorter, ks.Field)
	}
	return sorter
}
