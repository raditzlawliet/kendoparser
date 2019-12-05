// helper parser for dboxpipe
package kpmongo

import (
	"strings"

	"github.com/raditzlawliet/gokendoparser"
	"gopkg.in/mgo.v2/bson"
)

// Parser Parser
type Parser struct{}

// ParseFilter convert KendoFilter into bson.M/D/A for pipe combination automaticly
// return can @Nullable if filter and filters empty
func (parser Parser) ParseFilter(kf *gokendoparser.KendoFilter) interface{} {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
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
					return f
				}
				return nil
			}
		}

		// global operator
		if opHandler, ok := operatorManager.RegisteredOperators[kf.Operator]; ok {
			f := opHandler.Filter(ckFilter)
			if f != nil {
				return f
			}
			return nil
		}

		// defaultx
		f := operatorManager.DefaultOperator.Filter(ckFilter)
		if f != nil {
			return f
		}
		return nil
	}

	// so filters has some values
	filters := []bson.D{}
	for _, kFilterChild := range kf.Filters {
		filter := parser.ParseFilter(&kFilterChild)
		if filter != nil {
			switch filterAssertion := filter.(type) {
			case bson.D:
				filters = append(filters, filterAssertion)
			case bson.M:
				d := bson.D{}
				for k, v := range filterAssertion {
					d = append(d, bson.DocElem{k, v})
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

	return nil
}

// ParserSort ParserSort
func (parser Parser) ParserSort(ksa *gokendoparser.KendoSortArray) interface{} {
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
