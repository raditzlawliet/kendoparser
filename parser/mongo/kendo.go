// helper parser for dboxpipe
package kpmongo

import (
	"strings"

	"github.com/raditzlawliet/kendoparser"
	"go.mongodb.org/mongo-driver/bson"
)

// FilterParser convert Filter into bson.M/D/A for pipe combination automaticly
// return can @Nullable if filter and filters empty
func FilterParser(kf *kendoparser.Filter) interface{} {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
	if len(kf.Filters) == 0 {
		if kf.Operator == "" {
			return bson.D{{}}
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
				return bson.D{{}}
			}
		}

		// (global)
		if opHandler, ok := OperatorManager.OperatorFilters[kf.Operator]; ok {
			f := opHandler(ckFilter)
			if f != nil {
				return f
			}
			return bson.D{{}}
		}

		// default (global)
		f := OperatorManager.DefaultOperatorFilter(ckFilter)
		if f != nil {
			return f
		}
		return bson.D{{}}
	}

	// so filters has some values
	filters := []bson.D{}
	for _, kFilterChild := range kf.Filters {
		filter := FilterParser(&kFilterChild)
		if filter != nil {
			switch filterAssertion := filter.(type) {
			case bson.D:
				filters = append(filters, filterAssertion)
			case bson.M:
				d := bson.D{}
				for k, v := range filterAssertion {
					d = append(d, bson.E{k, v})
				}
				filters = append(filters, d)
			}
		}
	}

	if len(filters) > 0 {
		if strings.ToLower(kf.Logic) == "and" {
			return bson.D{{"$and", filters}}
		}
		return bson.D{{"$or", filters}}
	}

	return bson.D{{}}
}

// SortParser SortParser
func SortParser(ksa *kendoparser.Sort) interface{} {
	sorter := bson.D{}
	for _, ks := range *ksa {
		sort := 1
		if strings.ToLower(ks.Dir) == "desc" {
			sort = -1
		}
		sorter = append(sorter, bson.E{ks.Field, sort})
	}
	return sorter
}
