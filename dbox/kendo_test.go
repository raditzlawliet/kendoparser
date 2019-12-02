// package kpdbox

// import (
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/eaciit/dbox"
// 	tk "github.com/eaciit/toolkit"
// 	"github.com/stretchr/testify/require"
// 	"gopkg.in/mgo.v2/bson"
// )

// func Test_ToDboxFilter(t *testing.T) {
// 	// Single filter
// 	{
// 		kendoFilter := KendoFilter{
// 			Field: "_id", Operator: "eq", Value: "val",
// 		}

// 		resultFilter := kendoFilter.ToDboxFilter()
// 		expectedFilter := dbox.Eq("_id", "val")
// 		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

// 		kendoRequest := KendoRequest{}
// 		e := tk.UnjsonFromString(`{
// 			"data": {
// 				"filter": {
// 					"field": "_id",
// 					"operator": "eq",
// 					"value": "val"
// 				}
// 			}
// 		}`, &kendoRequest)
// 		require.Nil(t, e, "Json parse must work")
// 		require.Equal(t, kendoFilter, kendoRequest.Data.Filter, "Filter must same")
// 		resultFilterJSON := kendoRequest.Data.Filter.ToDboxFilter()
// 		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
// 	}

// 	{
// 		kendoFilter := KendoFilter{
// 			Filters: []KendoFilter{
// 				KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
// 			},
// 			Logic: "and",
// 		}
// 		resultFilter := kendoFilter.ToDboxFilter()
// 		expectedFilter := dbox.And(dbox.Eq("_id", "val"))
// 		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
// 	}

// 	{
// 		kendoFilter := KendoFilter{
// 			Filters: []KendoFilter{
// 				KendoFilter{
// 					Filters: []KendoFilter{
// 						KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
// 						KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
// 					},
// 					Logic: "or",
// 				},
// 				KendoFilter{
// 					Filters: []KendoFilter{
// 						KendoFilter{Field: "_id", Operator: "eq", Value: "val2"},
// 						KendoFilter{Field: "_id", Operator: "neq", Value: "val2"},
// 					},
// 					Logic: "or",
// 				},
// 			},
// 			Logic: "and",
// 		}
// 		resultFilter := kendoFilter.ToDboxFilter()
// 		expectedFilter := dbox.And(
// 			dbox.Or(
// 				dbox.Eq("_id", "val"),
// 				dbox.Ne("_id", "val"),
// 			),
// 			dbox.Or(
// 				dbox.Eq("_id", "val2"),
// 				dbox.Ne("_id", "val2"),
// 			),
// 		)
// 		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
// 	}

// 	// operator check
// 	{
// 		kendoFilter := KendoFilter{
// 			Filters: []KendoFilter{
// 				KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
// 				KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
// 				KendoFilter{Field: "_id", Operator: "doesnotcontain", Value: "val"},
// 				KendoFilter{Field: "_id", Operator: "contain", Value: "val"},
// 				KendoFilter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
// 				KendoFilter{Field: "_id", Operator: "gte", Value: "val"},
// 				KendoFilter{Field: "_id", Operator: "lte", Value: "val"},
// 				KendoFilter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
// 				KendoFilter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
// 				KendoFilter{Field: "_id", Operator: "unknown", Value: "val"},
// 			},
// 			Logic: "and",
// 		}
// 		resultFilter := kendoFilter.ToDboxFilter()
// 		testTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
// 		expectedFilter := dbox.And(
// 			dbox.Eq("_id", "val"),
// 			dbox.Ne("_id", "val"),
// 			dbox.Eq("_id", tk.M{"$ne": tk.M{"$regex": `` + "val" + ``, "$options": "i"}}),
// 			dbox.Contains("_id", "val"),
// 			dbox.In("_id", []interface{}{"val"}...),
// 			dbox.Gte("_id", "val"),
// 			dbox.Lte("_id", "val"),
// 			dbox.Gte("time", testTime),
// 			dbox.Lte("time", testTime),
// 			dbox.Eq("_id", "val"),
// 		)
// 		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
// 	}
// }

// func Test_PreFilterHandler(t *testing.T) {
// 	// transform single filter
// 	// ID => _id
// 	{
// 		kendoFilter := KendoFilter{
// 			Filters: []KendoFilter{
// 				KendoFilter{
// 					Filters: []KendoFilter{
// 						KendoFilter{Field: "ID", Operator: "eq", Value: "val"},
// 						KendoFilter{Field: "STATUS", Operator: "eq", Value: "true"},
// 					},
// 					Logic: "or",
// 				},
// 				KendoFilter{
// 					Filters: []KendoFilter{
// 						KendoFilter{Field: "ID", Operator: "eq", Value: "val2"},
// 						KendoFilter{Field: "ID", Operator: "neq", Value: "val2"},
// 					},
// 					Logic: "or",
// 				},
// 			},
// 			Logic: "and",
// 		}
// 		// try dbox filter
// 		resultFilter := kendoFilter.TransformAllField(strings.ToLower).
// 			TransformAll(func(kf *KendoFilter) {
// 				if kf.Field == "id" {
// 					kf.Field = "_id"
// 				}
// 			}).
// 			PreDboxFilterAll(func(kf *KendoFilter) *dbox.Filter {
// 				if kf.Field == "status" {
// 					// return your custom handler
// 					return dbox.Eq(kf.Field, StringToBool(kf.Value, false))
// 				}
// 				return nil // pas nil to continue original filter
// 			}).
// 			ToDboxFilter()

// 		// reset if needed another
// 		kendoFilter.ResetPreFilter()

// 		expectedFilter := dbox.And(
// 			dbox.Or(
// 				dbox.Eq("_id", "val"),
// 				dbox.Eq("status", true),
// 			),
// 			dbox.Or(
// 				dbox.Eq("_id", "val2"),
// 				dbox.Ne("_id", "val2"),
// 			),
// 		)
// 		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

// 		// try dbox pipe
// 		resultFilterPipe := kendoFilter.TransformAllField(strings.ToLower).
// 			TransformAll(func(kf *KendoFilter) {
// 				if kf.Field == "id" {
// 					kf.Field = "_id"
// 				}
// 			}).
// 			PreDboxPipeAll(func(kf *KendoFilter) tk.M {
// 				if kf.Field == "status" {
// 					// return your custom handler
// 					return tk.M{kf.Field: StringToBool(kf.Value, false)}
// 				}
// 				return nil // pas nil to continue original filter
// 			}).
// 			ToDboxPipe()

// 		expectedFilterPipe := tk.M{"$and": []tk.M{
// 			tk.M{"$or": []tk.M{
// 				tk.M{"_id": tk.M{"$eq": "val"}},
// 				tk.M{"status": true},
// 			}},
// 			tk.M{"$or": []tk.M{
// 				tk.M{"_id": tk.M{"$eq": "val2"}},
// 				tk.M{"_id": tk.M{"$ne": "val2"}},
// 			}},
// 		}}
// 		require.Equal(t, expectedFilterPipe, resultFilterPipe, "Result dbox filter must same")
// 	}
// }

// func Test_Sort(t *testing.T) {
// 	{
// 		kData := KendoData{
// 			Sort: KendoSortArray{
// 				KendoSort{
// 					Field: "foo",
// 					Dir:   "DESC",
// 				},
// 				KendoSort{
// 					Field: "bar",
// 					Dir:   "ASC",
// 				},
// 				KendoSort{
// 					Field: "_id",
// 					Dir:   "desc",
// 				},
// 			},
// 		}

// 		// try dbox filter
// 		result := kData.Sort.ToDbox()
// 		resultPipe := kData.Sort.ToDboxPipe()

// 		expected := []string{"-foo", "bar", "-_id"}
// 		expectedPipe := bson.D{
// 			bson.DocElem{
// 				Name:  "foo",
// 				Value: -1,
// 			},
// 			bson.DocElem{
// 				Name:  "bar",
// 				Value: 1,
// 			},
// 			bson.DocElem{
// 				Name:  "_id",
// 				Value: -1,
// 			},
// 		}

// 		require.Equal(t, expected, result, "Result must same")
// 		require.Equal(t, expectedPipe, resultPipe, "Result must same")
// 	}
// }
