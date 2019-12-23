// helper parser for dboxpipe
package kpxorm

import (
	"fmt"
	"strings"

	"github.com/raditzlawliet/gokendoparser"

	"xorm.io/builder"
)

// Parser Parser
type Parser struct{}

// ParseFilter convert KendoFilter into bson.M/D/A for pipe combination automaticly
// return can @Nullable if filter and filters empty
func (parser Parser) ParseFilter(kf *gokendoparser.KendoFilter) interface{} {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
	if len(kf.Filters) == 0 {
		if kf.Operator == "" {
			return builder.NewCond()
		}

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
				return builder.NewCond()
			}
		}

		// global operator
		if opHandler, ok := operatorManager.RegisteredOperators[kf.Operator]; ok {
			f := opHandler.Filter(ckFilter)
			if f != nil {
				return f
			}
			return builder.NewCond()
		}

		// defaultx
		f := operatorManager.DefaultOperator.Filter(ckFilter)
		if f != nil {
			return f
		}
		return builder.NewCond()
	}

	// so filters has some values
	filters := []builder.Cond{}
	for _, kFilterChild := range kf.Filters {
		filter := parser.ParseFilter(&kFilterChild)
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

// ParserSort ParserSort
func (parser Parser) ParserSort(ksa *gokendoparser.KendoSortArray) interface{} {
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
