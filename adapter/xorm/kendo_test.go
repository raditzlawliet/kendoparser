package kpxorm

import (
	"strings"
	"testing"

	"github.com/raditzlawliet/gokendoparser"
	"github.com/stretchr/testify/require"

	"xorm.io/builder"
)

func Test_ParseFilter(t *testing.T) {
	kendoFilter := gokendoparser.KendoFilter{
		Filters: []gokendoparser.KendoFilter{
			gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter := kendoFilter.Parse(Parser{})
	expectedFilter := builder.And(
		builder.Eq{"_id": "val"},
	)
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	sql, params, err := builder.ToSQL(resultFilter)
	require.Nil(t, err)
	t.Log(sql, params)

	kendoFilter = gokendoparser.KendoFilter{
		Filters: []gokendoparser.KendoFilter{
			gokendoparser.KendoFilter{
				Filters: []gokendoparser.KendoFilter{
					gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
					gokendoparser.KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
				},
				Logic: "or",
			},
			gokendoparser.KendoFilter{
				Filters: []gokendoparser.KendoFilter{
					gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val2"},
					gokendoparser.KendoFilter{Field: "_id", Operator: "neq", Value: "val2"},
				},
				Logic: "or",
			},
		},
		Logic: "and",
	}
	resultFilter = kendoFilter.Parse(Parser{})

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
	// kendoFilter = gokendoparser.KendoFilter{
	// 	Filters: []gokendoparser.KendoFilter{
	// 		gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
	// 		gokendoparser.KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
	// 		gokendoparser.KendoFilter{Field: "_id", Operator: "contains", Value: "val"},
	// 		gokendoparser.KendoFilter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
	// 		gokendoparser.KendoFilter{Field: "_id", Operator: "gte", Value: "val"},
	// 		gokendoparser.KendoFilter{Field: "_id", Operator: "lte", Value: "val"},
	// 		gokendoparser.KendoFilter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
	// 		gokendoparser.KendoFilter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
	// 		gokendoparser.KendoFilter{Field: "_id", Operator: "unknown", Value: "val"},
	// 	},
	// 	Logic: "and",
	// }
	// resultFilter = kendoFilter.Parse(Parser{}).(bson.D)
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
// 		kendoFilter := gokendoparser.KendoFilter{
// 			Filters: []gokendoparser.KendoFilter{
// 				gokendoparser.KendoFilter{
// 					Filters: []gokendoparser.KendoFilter{
// 						gokendoparser.KendoFilter{Field: "ID", Operator: "eq", Value: "val"},
// 						gokendoparser.KendoFilter{Field: "STATUS", Operator: "eq", Value: "true"},
// 					},
// 					Logic: "or",
// 				},
// 				gokendoparser.KendoFilter{
// 					Filters: []gokendoparser.KendoFilter{
// 						gokendoparser.KendoFilter{Field: "ID", Operator: "eq", Value: "val2"},
// 						gokendoparser.KendoFilter{Field: "ID", Operator: "neq", Value: "val2"},
// 					},
// 					Logic: "or",
// 				},
// 			},
// 			Logic: "and",
// 		}

// 		// try dbox pipe
// 		resultFilter := kendoFilter.TransformAllField(strings.ToLower).
// 			TransformAll(func(kf *gokendoparser.KendoFilter) {
// 				if kf.Field == "id" {
// 					kf.Field = "_id"
// 				}
// 			}).
// 			BeforeParseAll(func(kf *gokendoparser.KendoFilter) interface{} {
// 				if kf.Field == "status" {
// 					// return your custom handler
// 					return bson.M{kf.Field: helper.StringToBool(kf.Value, false)}
// 				}
// 				return nil // pas nil to continue original filter
// 			}).Parse(Parser{}).(bson.D)

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
		kData := gokendoparser.KendoData{
			Sort: gokendoparser.KendoSortArray{
				gokendoparser.KendoSort{
					Field: "foo",
					Dir:   "DESC",
				},
				gokendoparser.KendoSort{
					Field: "bar",
					Dir:   "ASC",
				},
				gokendoparser.KendoSort{
					Field: "_id",
					Dir:   "desc",
				},
			},
		}

		// try dbox filter
		result := kData.Sort.Parse(Parser{})
		expectedPipe := strings.Join([]string{
			"foo DESC", "bar ASC", "_id DESC",
		}, ", ")
		require.Equal(t, expectedPipe, result, "Result must same")
	}
}
