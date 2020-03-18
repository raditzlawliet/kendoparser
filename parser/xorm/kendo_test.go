package kpxorm

import (
	"strings"
	"testing"

	"github.com/raditzlawliet/kendoparser"
	"github.com/stretchr/testify/require"

	"xorm.io/builder"
)

func Test_FilterParser(t *testing.T) {
	kendoFilter := kendoparser.Filter{
		Filters: []kendoparser.Filter{
			kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter := kendoFilter.Parse(FilterParser)
	expectedFilter := builder.And(
		builder.Eq{"_id": "val"},
	)
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	sql, params, err := builder.ToSQL(resultFilter)
	require.Nil(t, err)
	t.Log(sql, params)

	kendoFilter = kendoparser.Filter{
		Filters: []kendoparser.Filter{
			kendoparser.Filter{
				Filters: []kendoparser.Filter{
					kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
					kendoparser.Filter{Field: "_id", Operator: "neq", Value: "val"},
				},
				Logic: "or",
			},
			kendoparser.Filter{
				Filters: []kendoparser.Filter{
					kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val2"},
					kendoparser.Filter{Field: "_id", Operator: "neq", Value: "val2"},
				},
				Logic: "or",
			},
		},
		Logic: "and",
	}
	resultFilter = kendoFilter.Parse(FilterParser)

	expectedFilter = builder.And(
		builder.Or(
			builder.Eq{"_id": "val"},
			builder.Neq{"_id": "val"},
		),
		builder.Or(
			builder.Eq{"_id": "val2"},
			builder.Neq{"_id": "val2"},
		),
	)
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	sql, params, err = builder.ToSQL(resultFilter)
	require.Nil(t, err)
	t.Log(sql, params)

	// // operator check
	// kendoFilter = kendoparser.Filter{
	// 	Filters: []kendoparser.Filter{
	// 		kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
	// 		kendoparser.Filter{Field: "_id", Operator: "neq", Value: "val"},
	// 		kendoparser.Filter{Field: "_id", Operator: "contains", Value: "val"},
	// 		kendoparser.Filter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
	// 		kendoparser.Filter{Field: "_id", Operator: "gte", Value: "val"},
	// 		kendoparser.Filter{Field: "_id", Operator: "lte", Value: "val"},
	// 		kendoparser.Filter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
	// 		kendoparser.Filter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
	// 		kendoparser.Filter{Field: "_id", Operator: "unknown", Value: "val"},
	// 	},
	// 	Logic: "and",
	// }
	// resultFilter = kendoFilter.Parse(FilterParser).(bson.D)
	// testTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")

	// expectedFilter = bson.D{
	// 	{
	// 		"$and", []bson.D{
	// 			bson.D{{"_id", bson.M{"$eq": "val"}}},
	// 			bson.D{{"_id", bson.M{"$ne": "val"}}},
	// 			// bson.D{{"_id", RegexContains("val", false)}},
	// 			bson.D{{"_id", bson.M{"$in": []interface{}{"val"}}}},
	// 			bson.D{{"_id", bson.M{"$gte": "val"}}},
	// 			bson.D{{"_id", bson.M{"$lte": "val"}}},
	// 			bson.D{{"time", bson.M{"$gte": testTime}}},
	// 			bson.D{{"time", bson.M{"$lte": testTime}}},
	// 			bson.D{{"_id", bson.M{"$eq": "val"}}},
	// 		},
	// 	},
	// }
	// require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
}

// func Test_PreFilterHandler(t *testing.T) {
// 	// transform single filter
// 	// ID => _id
// 	{
// 		kendoFilter := kendoparser.Filter{
// 			Filters: []kendoparser.Filter{
// 				kendoparser.Filter{
// 					Filters: []kendoparser.Filter{
// 						kendoparser.Filter{Field: "ID", Operator: "eq", Value: "val"},
// 						kendoparser.Filter{Field: "STATUS", Operator: "eq", Value: "true"},
// 					},
// 					Logic: "or",
// 				},
// 				kendoparser.Filter{
// 					Filters: []kendoparser.Filter{
// 						kendoparser.Filter{Field: "ID", Operator: "eq", Value: "val2"},
// 						kendoparser.Filter{Field: "ID", Operator: "neq", Value: "val2"},
// 					},
// 					Logic: "or",
// 				},
// 			},
// 			Logic: "and",
// 		}

// 		// try dbox pipe
// 		resultFilter := kendoFilter.TransformAllField(strings.ToLower).
// 			TransformAll(func(kf *kendoparser.Filter) {
// 				if kf.Field == "id" {
// 					kf.Field = "_id"
// 				}
// 			}).
// 			AddAllParser(func(kf *kendoparser.Filter) interface{} {
// 				if kf.Field == "status" {
// 					// return your custom handler
// 					return bson.M{kf.Field: helper.StringToBool(kf.Value, false)}
// 				}
// 				return nil // pas nil to continue original filter
// 			}).Parse(FilterParser).(bson.D)

// 		expectedFilter := bson.D{
// 			{
// 				"$and", []bson.D{
// 					bson.D{{"$or", []bson.D{
// 						bson.D{{"_id", bson.M{"$eq": "val"}}},
// 						bson.D{{"status", true}},
// 					}}},
// 					bson.D{{"$or", []bson.D{
// 						bson.D{{"_id", bson.M{"$eq": "val2"}}},
// 						bson.D{{"_id", bson.M{"$ne": "val2"}}},
// 					}}},
// 				},
// 			},
// 		}

// 		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
// 	}
// }

func Test_Sort(t *testing.T) {
	{
		kData := kendoparser.Data{
			Sort: kendoparser.Sort{
				kendoparser.SortDetail{
					Field: "foo",
					Dir:   "DESC",
				},
				kendoparser.SortDetail{
					Field: "bar",
					Dir:   "ASC",
				},
				kendoparser.SortDetail{
					Field: "_id",
					Dir:   "desc",
				},
			},
		}

		// try dbox filter
		result := kData.Sort.Parse(SortParser)
		expectedPipe := strings.Join([]string{
			"foo DESC", "bar ASC", "_id DESC",
		}, ", ")
		require.Equal(t, expectedPipe, result, "Result must same")
	}
}
