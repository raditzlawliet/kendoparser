package gokendoparser

/*
 * SCB-DMT Project Team
 * Kendo data source payload to dbox
 */

import (
	"regexp"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

// KendoRequest option variable to struct (each apps has different format, defined yourself if needed)
type KendoRequest struct {
	Data KendoData `json:"data"`
}

// KendoData datasource payload
type KendoData struct {
	Filter   KendoFilter `json:"filter"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Skip     int         `json:"skip"`
	Take     int         `json:"take"`
	Sort     []KendoSort `json:"sort"`
}

// KendoSort struct sort
type KendoSort struct {
	Dir   string `json:"dir"`
	Field string `json:"field"`
}

// KendoFilter struct filters
type KendoFilter struct {
	Filters    []KendoFilter `json:"filters"`
	Logic      string        `json:"logic"`
	Field      string        `json:"field"`
	Operator   string        `json:"operator"`
	IgnoreCase bool          `json:"ignoreCase"`
	Value      string        `json:"value"`
	Values     []interface{} `json:"values"`
}

// ToDboxFilter convert KendoFilter into *dbox.Filter combination automaticly
// return can @Nullable if filter and filters empty
func (kFilter *KendoFilter) ToDboxFilter(fieldValidator func(field *string, filter *KendoFilter) *dbox.Filter) *dbox.Filter {
	// defaultDboxFilter := &dbox.Filter{
	// 	Field: "_id", // assume all collection always has field _id
	// 	Op:    dbox.FilterOpEqual,
	// 	Value: tk.M{
	// 		"$exists": true,
	// 	},
	// }

	if len(kFilter.Filters) == 0 {
		var f *dbox.Filter
		if !tk.IsNilOrEmpty(kFilter.Operator) && !tk.IsNilOrEmpty(kFilter.Field) {
			field := strings.ToLower(kFilter.Field)
			if fieldValidator != nil {
				f = fieldValidator(&field, kFilter)
				if f != nil {
					return f
				}
			}
			switch kFilter.Operator {
			case "eq":
				f = dbox.Eq(field, kFilter.Value)
			case "ne":
				fallthrough
			case "neq":
				f = dbox.Ne(field, kFilter.Value)
			case "doesnotcontain":
				value := regexp.QuoteMeta(kFilter.Value)
				f = &dbox.Filter{
					Field: field,
					Op:    dbox.FilterOpEqual, // equal are field = value and can be manipulate for others
					Value: tk.M{"$ne": tk.M{
						"$regex":   `` + value + ``,
						"$options": "i",
					}},
				}
			case "contain":
				fallthrough
			case "contains":
				f = dbox.Contains(field, kFilter.Value)
			case "in":
				f = dbox.In(field, kFilter.Values...)
			case "gte":
				f = dbox.Gte(field, kFilter.Value)
			case "lte":
				f = dbox.Lte(field, kFilter.Value)
			case "gtedate":
				_dtVariable, _ := time.Parse(time.RFC3339, kFilter.Value)
				f = dbox.Gte(field, _dtVariable)
			case "ltedate":
				_dtVariable, _ := time.Parse(time.RFC3339, kFilter.Value)
				f = dbox.Lte(field, _dtVariable)
			case "exists":
				f = &dbox.Filter{
					Field: field,
					Op:    dbox.FilterOpEqual,
					Value: tk.M{
						"$exists": true,
					},
				}
			default:
				f = dbox.Eq(field, kFilter.Value)
			}
			return f
		}
		return nil
	}

	// so filters has some values
	dboxFilters := []*dbox.Filter{}
	for _, kFilterChild := range kFilter.Filters {
		dboxFilter := kFilterChild.ToDboxFilter(fieldValidator)
		if dboxFilter != nil {
			dboxFilters = append(dboxFilters, dboxFilter)
		}
	}

	if len(dboxFilters) > 0 {
		if strings.ToLower(kFilter.Logic) == "and" {
			return dbox.And(dboxFilters...)
		}
		return dbox.Or(dboxFilters...)
	}

	return nil // can return nil if filter & filters are meh ...
}

// ToAggregationFilter convert KendoFilter into M for pipe combination automaticly
// return can @Nullable if filter and filters empty
func (kFilter *KendoFilter) ToAggregationFilter(fieldValidator func(field *string, filter *KendoFilter) tk.M) tk.M {
	// defaultFilter := tk.M{"_id": tk.M{"$exists": true}}
	if len(kFilter.Filters) == 0 {
		var f tk.M
		if !tk.IsNilOrEmpty(kFilter.Operator) && !tk.IsNilOrEmpty(kFilter.Field) {
			field := strings.ToLower(kFilter.Field)
			if fieldValidator != nil {
				f = fieldValidator(&field, kFilter)
				if f != nil {
					return f
				}
			}

			switch kFilter.Operator {
			case "eq":
				if kFilter.IgnoreCase {
					value := regexp.QuoteMeta(kFilter.Value)
					f = tk.M{field: bson.RegEx{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
				} else {
					f = tk.M{field: tk.M{"$eq": kFilter.Value}}
				}
			case "ne":
				fallthrough
			case "neq":
				f = tk.M{field: tk.M{"$ne": kFilter.Value}}
			case "doesnotcontain":
				f = tk.M{field: tk.M{"$ne": RegexContains(kFilter.Value, kFilter.IgnoreCase)}}
			case "contain":
				fallthrough
			case "contains":
				f = tk.M{field: RegexContains(kFilter.Value, kFilter.IgnoreCase)}
			case "in":
				f = tk.M{field: tk.M{"$in": kFilter.Values}}
			case "gte":
				f = tk.M{field: tk.M{"$gte": kFilter.Value}}
			case "lte":
				f = tk.M{field: tk.M{"$lte": kFilter.Value}}
			case "gtedate":
				_dtVariable, _ := time.Parse(time.RFC3339, kFilter.Value)
				f = tk.M{field: tk.M{"$gte": _dtVariable}}
			case "ltedate":
				_dtVariable, _ := time.Parse(time.RFC3339, kFilter.Value)
				f = tk.M{field: tk.M{"$lte": _dtVariable}}
			case "exists":
				f = tk.M{field: tk.M{"$exists": StringToBool(kFilter.Value, false)}}
			default:
				if kFilter.IgnoreCase {
					f = tk.M{field: RegexCaseInsensitive(kFilter.Value)}
				} else {
					f = tk.M{field: tk.M{"$eq": kFilter.Value}}
				}
			}
			return f
		}
		return nil
	}

	// so filters has some values
	filters := []tk.M{}
	for _, kFilterChild := range kFilter.Filters {
		filter := kFilterChild.ToAggregationFilter(fieldValidator)
		if filter != nil {
			filters = append(filters, filter)
		}
	}

	if len(filters) > 0 {
		if strings.ToLower(kFilter.Logic) == "and" {
			return tk.M{"$and": filters}
		}
		return tk.M{"$or": filters}
	}

	return nil
}
