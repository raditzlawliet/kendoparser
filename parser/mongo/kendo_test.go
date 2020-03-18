package kpmongo

import (
	"github.com/raditzlawliet/kendoparser"
	"github.com/raditzlawliet/kendoparser/helper"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"testing"
	"time"
)

func Test_Parser(t *testing.T) {
	kendoFilter := kendoparser.Filter{
		Filters: []kendoparser.Filter{
			kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter := kendoFilter.Parse(FilterParser).(bson.D)
	expectedFilter := bson.D{
		{
			"$and", []bson.D{
			bson.D{{"_id", bson.M{"$eq": "val"}}},
		},
		},
	}
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")

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
	resultFilter = kendoFilter.Parse(FilterParser).(bson.D)

	expectedFilter = bson.D{
		{
			"$and", []bson.D{
			bson.D{{"$or", []bson.D{
				bson.D{{"_id", bson.M{"$eq": "val"}}},
				bson.D{{"_id", bson.M{"$ne": "val"}}},
			}}},
			bson.D{{"$or", []bson.D{
				bson.D{{"_id", bson.M{"$eq": "val2"}}},
				bson.D{{"_id", bson.M{"$ne": "val2"}}},
			}}},
		},
		},
	}
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")

	// operator check
	kendoFilter = kendoparser.Filter{
		Filters: []kendoparser.Filter{
			kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
			kendoparser.Filter{Field: "_id", Operator: "neq", Value: "val"},
			kendoparser.Filter{Field: "_id", Operator: "contains", Value: "val"},
			kendoparser.Filter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
			kendoparser.Filter{Field: "_id", Operator: "gte", Value: "val"},
			kendoparser.Filter{Field: "_id", Operator: "lte", Value: "val"},
			kendoparser.Filter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
			kendoparser.Filter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
			kendoparser.Filter{Field: "_id", Operator: "unknown", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter = kendoFilter.Parse(FilterParser).(bson.D)
	testTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")

	expectedFilter = bson.D{
		{
			"$and", []bson.D{
			bson.D{{"_id", bson.M{"$eq": "val"}}},
			bson.D{{"_id", bson.M{"$ne": "val"}}},
			bson.D{{"_id", RegexContains("val", false)}},
			bson.D{{"_id", bson.M{"$in": []interface{}{"val"}}}},
			bson.D{{"_id", bson.M{"$gte": "val"}}},
			bson.D{{"_id", bson.M{"$lte": "val"}}},
			bson.D{{"time", bson.M{"$gte": testTime}}},
			bson.D{{"time", bson.M{"$lte": testTime}}},
			bson.D{{"_id", bson.M{"$eq": "val"}}},
		},
		},
	}
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
}

func Test_AdditionalParser(t *testing.T) {
	// transform single filter
	// ID => _id
	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "ID", Operator: "eq", Value: "val"},
						kendoparser.Filter{Field: "STATUS", Operator: "eq", Value: "true"},
					},
					Logic: "or",
				},
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "ID", Operator: "eq", Value: "val2"},
						kendoparser.Filter{Field: "ID", Operator: "neq", Value: "val2"},
					},
					Logic: "or",
				},
			},
			Logic: "and",
		}

		// try
		resultFilter := kendoFilter.TransformAllField(strings.ToLower).
			TransformAll(func(kf *kendoparser.Filter) {
				if kf.Field == "id" {
					kf.Field = "_id"
				}
			}).
			AddAllParser(func(kf *kendoparser.Filter) interface{} {
				if kf.Field == "status" {
					// return your custom handler
					return bson.M{kf.Field: helper.StringToBool(cast.ToString(kf.Value), false)}
				}
				return nil // pas nil to continue original filter
			}).Parse(FilterParser).(bson.D)
		expectedFilter := bson.D{
			{
				"$and", []bson.D{
				bson.D{{"$or", []bson.D{
					bson.D{{"_id", bson.M{"$eq": "val"}}},
					bson.D{{"status", true}},
				}}},
				bson.D{{"$or", []bson.D{
					bson.D{{"_id", bson.M{"$eq": "val2"}}},
					bson.D{{"_id", bson.M{"$ne": "val2"}}},
				}}},
			},
			},
		}
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")

		// reset again, so additional parser will gone
		resultFilter = kendoFilter.TransformAllField(strings.ToLower).
			TransformAll(func(kf *kendoparser.Filter) {
				if kf.Field == "id" {
					kf.Field = "_id"
				}
			}).
			AddAllParser(func(kf *kendoparser.Filter) interface{} {
				if kf.Field == "status" {
					// return your custom handler
					return bson.M{kf.Field: helper.StringToBool(cast.ToString(kf.Value), false)}
				}
				return nil // pas nil to continue original filter
			}).ResetAllAdditionalParsers().Parse(FilterParser).(bson.D)
		expectedFilter = bson.D{
			{
				"$and", []bson.D{
				bson.D{{"$or", []bson.D{
					bson.D{{"_id", bson.M{"$eq": "val"}}},
					bson.D{{"status", bson.M{"$eq": "true"}}},
				}}},
				bson.D{{"$or", []bson.D{
					bson.D{{"_id", bson.M{"$eq": "val2"}}},
					bson.D{{"_id", bson.M{"$ne": "val2"}}},
				}}},
			},
			},
		}
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")

		// Reset only the wrapper (but not the childs), so reset will not work
		resultFilter = kendoFilter.TransformAllField(strings.ToLower).
			TransformAll(func(kf *kendoparser.Filter) {
				if kf.Field == "id" {
					kf.Field = "_id"
				}
			}).
			AddAllParser(func(kf *kendoparser.Filter) interface{} {
				if kf.Field == "status" {
					// return your custom handler
					return bson.M{kf.Field: helper.StringToBool(cast.ToString(kf.Value), false)}
				}
				return nil // pas nil to continue original filter
			}).ResetAdditionalParsers().Parse(FilterParser).(bson.D)
		expectedFilter = bson.D{
			{
				"$and", []bson.D{
				bson.D{{"$or", []bson.D{
					bson.D{{"_id", bson.M{"$eq": "val"}}},
					bson.D{{"status", true}},
				}}},
				bson.D{{"$or", []bson.D{
					bson.D{{"_id", bson.M{"$eq": "val2"}}},
					bson.D{{"_id", bson.M{"$ne": "val2"}}},
				}}},
			},
			},
		}
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	}
}

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

		// try filter
		result := kData.Sort.Parse(SortParser).(bson.D)

		expectedPipe := bson.D{
			bson.E{"foo", -1},
			bson.E{
				"bar",
				1,
			},
			bson.E{
				"_id",
				-1,
			},
		}

		require.Equal(t, expectedPipe, result, "Result must same")
	}
}

func Test_OperatorManager(t *testing.T) {
	om := kendoparser.OperatorManager{}
	om.RegisterOperator(func(kf kendoparser.Filter) interface{} {
		return bson.M{kf.Field: bson.M{"$eq": 1}}
	}, "eq")

	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.Parse(FilterParser).(bson.D)
		expectedFilter := bson.D{
			{
				"$and", []bson.D{
				bson.D{{"_id", bson.M{"$eq": "val"}}},
			},
			},
		}
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	}

	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.SetOperatorManager(&om).Parse(FilterParser).(bson.D)
		expectedFilter := bson.D{
			{
				"$and", []bson.D{
				bson.D{{"_id", bson.M{"$eq": 1}}},
			},
			},
		}
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	}
	om.Reset()

	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.Parse(FilterParser).(bson.D)
		expectedFilter := bson.D{
			{
				"$and", []bson.D{
				bson.D{{"_id", bson.M{"$eq": "val"}}},
			},
			},
		}
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	}
	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.SetOperatorManager(nil).Parse(FilterParser).(bson.D)
		expectedFilter := bson.D{
			{
				"$and", []bson.D{
				bson.D{{"_id", bson.M{"$eq": "val"}}},
			},
			},
		}
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	}
}
