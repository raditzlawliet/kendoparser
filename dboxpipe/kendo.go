// helper parser for dboxpipe
package kpdboxpipe

import (
	"strings"

	"github.com/eaciit/toolkit"
	"github.com/raditzlawliet/gokendoparser"
	"gopkg.in/mgo.v2/bson"
)

// ParseFilter convert KendoFilter into M for pipe combination automaticly
// return can @Nullable if filter and filters empty
func ParseFilter(kf *gokendoparser.KendoFilter) toolkit.M {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
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
					return f.(toolkit.M)
				}
				return nil
			}
		}

		if opHandler, ok := gokendoparser.RegisteredOperators[kf.Operator]; ok {
			f := opHandler.Filter(ckFilter)
			if f != nil {
				return f.(toolkit.M)
			}
			return nil
		}
		f := gokendoparser.DefaultOperator.Filter(ckFilter)
		if f != nil {
			return f.(toolkit.M)
		}
		return nil
	}

	// so filters has some values
	filters := []toolkit.M{}
	for _, kFilterChild := range kf.Filters {
		filter := ParseFilter(&kFilterChild)
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

// ParserSort ParserSort
func ParserSort(ksa *gokendoparser.KendoSortArray) bson.D {
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
