// helper parser for dbox
package kpdbflex

import (
	"strings"

	"git.eaciitapp.com/sebar/dbflex"

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
				return opHandler.Filter(ckFilter)
			}
		}

		// parser scope registered
		if opHandler, ok := operatorManager.RegisteredOperators[kf.Operator]; ok {
			return opHandler.Filter(ckFilter)
		}

		// global defult
		return operatorManager.DefaultOperator.Filter(ckFilter)
	}

	// so filters has some values
	filters := []*dbflex.Filter{}
	for _, kFilterChild := range kf.Filters {
		filter := parser.ParseFilter(&kFilterChild)
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
