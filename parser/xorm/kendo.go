// helper parser for dboxpipe
package kpxorm

import (
	"fmt"
	"strings"

	"github.com/raditzlawliet/kendoparser"

	"xorm.io/builder"
)

// Parser Parser
// type Parser struct{}

// FilterParser convert Filter into bson.M/D/A for pipe combination automaticly
// return can @Nullable if filter and filters empty
func FilterParser(kf *kendoparser.Filter) interface{} {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
	if len(kf.Filters) == 0 {
		if kf.Operator == "" {
			return builder.NewCond()
		}

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
					return f
				}
				return builder.NewCond()
			}
		}

		// (global)
		if opHandler, ok := OperatorManager.OperatorFilters[kf.Operator]; ok {
			f := opHandler(ckFilter)
			if f != nil {
				return f
			}
			return builder.NewCond()
		}

		// default (global)
		f := OperatorManager.DefaultOperatorFilter(ckFilter)
		if f != nil {
			return f
		}
		return builder.NewCond()
	}

	// so filters has some values
	filters := []builder.Cond{}
	for _, kFilterChild := range kf.Filters {
		filter := FilterParser(&kFilterChild)
		if filter != nil {
			filters = append(filters, filter.(builder.Cond))
		}
	}

	if len(filters) > 0 {
		if strings.ToLower(kf.Logic) == "and" {
			return builder.And(filters...)
		}
		return builder.Or(filters...)
	}
	return builder.NewCond()
}

// SortParser SortParser
func SortParser(ksa *kendoparser.Sort) interface{} {
	sorter := []string{}
	for _, ks := range *ksa {
		sort := "ASC"
		if strings.ToLower(ks.Dir) == "desc" {
			sort = "DESC"
		}
		sorter = append(sorter, fmt.Sprintf("%s %s", ks.Field, sort))
	}
	return strings.Join(sorter, ", ")
}
