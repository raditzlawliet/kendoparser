package gokendoparser

/*
 * @Author
 * Radityo <radityohernanda@gmail.com>
 */

import (
	"strings"

	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

// KendoRequest option variable to struct (each apps has different format, defined/extend yourself if needed)
type KendoRequest struct {
	Data KendoData `json:"data"`
}

// KendoData datasource payload
type KendoData struct {
	Filter   KendoFilter    `json:"filter"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Skip     int            `json:"skip"`
	Take     int            `json:"take"`
	Sort     KendoSortArray `json:"sort"`
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

	// will not change the original value
	registeredOperators map[string]Operator
	preDboxFilter       []func(*KendoFilter) *dbox.Filter
	preDboxPipe         []func(*KendoFilter) toolkit.M
}

// ToDboxFilter convert KendoFilter into *dbox.Filter combination automaticly
// return can @Nullable if filter and filters empty
func (kf *KendoFilter) ToDboxFilter() *dbox.Filter {
	// single filter
	if len(kf.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kf

		// very customable handler
		if kf.preDboxFilter != nil {
			for _, handler := range kf.preDboxFilter {
				f := handler(&ckFilter)
				if f != nil {
					return f
				}
			}
		}

		// local scope operator
		if kf.registeredOperators != nil {
			if opHandler, ok := kf.registeredOperators[kf.Operator]; ok && opHandler != nil {
				return opHandler.ToDboxFilter(ckFilter)
			}
		}

		if opHandler, ok := RegisteredOperators[kf.Operator]; ok {
			return opHandler.ToDboxFilter(ckFilter)
		}

		return DefaultOperator.ToDboxFilter(ckFilter)
	}

	// so filters has some values
	dboxFilters := []*dbox.Filter{}
	for _, kFilterChild := range kf.Filters {
		dboxFilter := kFilterChild.ToDboxFilter()
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

// ToDboxPipe convert KendoFilter into M for pipe combination automaticly
// return can @Nullable if filter and filters empty
func (kf *KendoFilter) ToDboxPipe() toolkit.M {
	// defaultFilter := toolkit.M{"_id": toolkit.M{"$exists": true}}
	if len(kf.Filters) == 0 {
		// processing will use copy instead to avoid change original value
		ckFilter := *kf

		// very customable handler
		if kf.preDboxPipe != nil {
			for _, handler := range kf.preDboxPipe {
				f := handler(&ckFilter)
				if f != nil {
					return f
				}
			}
		}

		// local scope operator
		if kf.registeredOperators != nil {
			if opHandler, ok := kf.registeredOperators[kf.Operator]; ok && opHandler != nil {
				return opHandler.ToDboxPipe(ckFilter)
			}
		}

		if opHandler, ok := RegisteredOperators[kf.Operator]; ok {
			return opHandler.ToDboxPipe(ckFilter)
		}
		return DefaultOperator.ToDboxPipe(ckFilter)
	}

	// so filters has some values
	filters := []toolkit.M{}
	for _, kFilterChild := range kf.Filters {
		filter := kFilterChild.ToDboxPipe()
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

// PreDboxFilter add custom handler pre-filtering
func (kf *KendoFilter) PreDboxFilter(f func(*KendoFilter) *dbox.Filter) *KendoFilter {
	if kf.preDboxFilter == nil {
		kf.preDboxFilter = []func(*KendoFilter) *dbox.Filter{}
	}
	if f != nil {
		kf.preDboxFilter = append(kf.preDboxFilter, f)
	}
	return kf
}

// PreDboxFilterAll add custom handler pre-filtering apply to nested struct
func (kf *KendoFilter) PreDboxFilterAll(f func(*KendoFilter) *dbox.Filter) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].PreDboxFilterAll(f)
	}
	kf.PreDboxFilter(f)
	return kf
}

// PreDboxPipe add custom handler pre-filtering
func (kf *KendoFilter) PreDboxPipe(f func(*KendoFilter) toolkit.M) *KendoFilter {
	if kf.preDboxPipe == nil {
		kf.preDboxPipe = []func(*KendoFilter) toolkit.M{}
	}
	if f != nil {
		kf.preDboxPipe = append(kf.preDboxPipe, f)
	}
	return kf
}

// PreDboxPipeAll add custom handler pre-filtering apply to nested struct
func (kf *KendoFilter) PreDboxPipeAll(f func(*KendoFilter) toolkit.M) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].PreDboxPipeAll(f)
	}
	kf.PreDboxPipe(f)
	return kf
}

// ResetPreFilter reset all pre-filter available
func (kf *KendoFilter) ResetPreFilter() *KendoFilter {
	kf.preDboxFilter = []func(*KendoFilter) *dbox.Filter{}
	kf.preDboxPipe = []func(*KendoFilter) toolkit.M{}
	return kf
}

// Transform your filter
func (kf *KendoFilter) Transform(t func(*KendoFilter)) *KendoFilter {
	t(kf)
	return kf
}

// TransformField only transform field
func (kf *KendoFilter) TransformField(t func(string) string) *KendoFilter {
	kf.Field = t(kf.Field)
	return kf
}

// TransformAll your filter include all childs
func (kf *KendoFilter) TransformAll(t func(*KendoFilter)) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].TransformAll(t)
	}
	kf.Transform(t)
	return kf
}

// TransformAllField only transform field include all childs
func (kf *KendoFilter) TransformAllField(t func(string) string) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].TransformAllField(t)
	}
	kf.TransformField(t)
	return kf
}

// RegisterOperator register operator local scope
func (kf *KendoFilter) RegisterOperator(op string, f Operator) *KendoFilter {
	if kf.registeredOperators == nil {
		kf.registeredOperators = map[string]Operator{}
	}
	kf.registeredOperators[op] = f
	return kf
}

// RegisterOperatorAll register operator local scope include childs
func (kf *KendoFilter) RegisterOperatorAll(op string, f Operator) *KendoFilter {
	for i := range kf.Filters {
		kf.Filters[i].RegisterOperatorAll(op, f)
	}
	kf.RegisterOperator(op, f)
	return kf
}

// KendoSort struct sort
type KendoSort struct {
	Dir   string `json:"dir"`
	Field string `json:"field"`
}

// KendoSortArray alias []KendoSort
type KendoSortArray []KendoSort

// ToDbox same with ToDboxFilter but for filter
func (ksa *KendoSortArray) ToDbox() []string {
	sorter := []string{}
	for _, ks := range *ksa {
		if strings.ToLower(ks.Dir) == "desc" {
			ks.Field = "-" + ks.Field
		}
		sorter = append(sorter, ks.Field)
	}
	return sorter
}

// ToDboxPipe same with ToAggreagateFilter but for sort
// bson.D can use map but saving the ordering
func (ksa *KendoSortArray) ToDboxPipe() bson.D {
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
