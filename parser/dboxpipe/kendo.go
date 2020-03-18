// helper parser for dboxpipe
package kpdboxpipe

import (
	"strings"

	"github.com/eaciit/toolkit"
	"github.com/raditzlawliet/kendoparser"
	"gopkg.in/mgo.v2/bson"
)

// Parser Parser
// type Parser struct{}

// FilterParser convert Filter into M for pipe combination automaticly
// return can @Nullable if filter and filters empty
func FilterParser(kf *kendoparser.Filter) interface{} {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
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
					return f.(toolkit.M)
				}
				return nil
			}
		}

		// (global)
		if opHandler, ok := OperatorManager.OperatorFilters[kf.Operator]; ok {
			f := opHandler(ckFilter)
			if f != nil {
				return f.(toolkit.M)
			}
			return nil
		}

		// default (global)
		f := OperatorManager.DefaultOperatorFilter(ckFilter)
		if f != nil {
			return f.(toolkit.M)
		}
		return nil
	}

	// so filters has some values
	filters := []toolkit.M{}
	for _, kFilterChild := range kf.Filters {
		filter := FilterParser(&kFilterChild)
		if filter != nil {
			filters = append(filters, filter.(toolkit.M))
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

// SortParser SortParser
func SortParser(ksa *kendoparser.Sort) interface{} {
	sorter := bson.D{}
	for _, ks := range *ksa {
		sort := 1
		if strings.ToLower(ks.Dir) == "desc" {
			sort = -1
		}
		sorter = append(sorter, bson.DocElem{
			Name:  ks.Field,
			Value: sort,
		})
	}
	return sorter
}
