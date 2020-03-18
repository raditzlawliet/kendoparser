package kpdboxpipe

import (
	"strings"
	"testing"
	"time"

	tk "github.com/eaciit/toolkit"
	"github.com/raditzlawliet/kendoparser"
	"github.com/raditzlawliet/kendoparser/helper"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func Test_FilterParser(t *testing.T) {
	kendoFilter := kendoparser.Filter{
		Filters: []kendoparser.Filter{
			kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter := kendoFilter.Parse(FilterParser).(tk.M)
	expectedFilter := tk.M{"$and": []tk.M{tk.M{"_id": tk.M{"$eq": "val"}}}}
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
	resultFilter = kendoFilter.Parse(FilterParser).(tk.M)
	expectedFilter = tk.M{"$and": []tk.M{
		tk.M{"$or": []tk.M{
			tk.M{"_id": tk.M{"$eq": "val"}},
			tk.M{"_id": tk.M{"$ne": "val"}},
		}},
		tk.M{"$or": []tk.M{
			tk.M{"_id": tk.M{"$eq": "val2"}},
			tk.M{"_id": tk.M{"$ne": "val2"}},
		}},
	}}
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
	resultFilter = kendoFilter.Parse(FilterParser).(tk.M)
	testTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
	expectedFilter = tk.M{"$and": []tk.M{
		tk.M{"_id": tk.M{"$eq": "val"}},
		tk.M{"_id": tk.M{"$ne": "val"}},
		tk.M{"_id": helper.RegexContains("val", false)},
		tk.M{"_id": tk.M{"$in": []interface{}{"val"}}},
		tk.M{"_id": tk.M{"$gte": "val"}},
		tk.M{"_id": tk.M{"$lte": "val"}},
		tk.M{"time": tk.M{"$gte": testTime}},
		tk.M{"time": tk.M{"$lte": testTime}},
		tk.M{"_id": tk.M{"$eq": "val"}},
	}}
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
}

func Test_PreFilterHandler(t *testing.T) {
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

		// try dbox pipe
		resultFilterPipe := kendoFilter.TransformAllField(strings.ToLower).
			TransformAll(func(kf *kendoparser.Filter) {
				if kf.Field == "id" {
					kf.Field = "_id"
				}
			}).
			AddAllParser(func(kf *kendoparser.Filter) interface{} {
				if kf.Field == "status" {
					// return your custom handler
					return tk.M{kf.Field: helper.StringToBool(cast.ToString(kf.Value), false)}
				}
				return nil // pas nil to continue original filter
			}).Parse(FilterParser).(tk.M)

		expectedFilterPipe := tk.M{"$and": []tk.M{
			tk.M{"$or": []tk.M{
				tk.M{"_id": tk.M{"$eq": "val"}},
				tk.M{"status": true},
			}},
			tk.M{"$or": []tk.M{
				tk.M{"_id": tk.M{"$eq": "val2"}},
				tk.M{"_id": tk.M{"$ne": "val2"}},
			}},
		}}
		require.Equal(t, expectedFilterPipe, resultFilterPipe, "Result dbox filter must same")
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

		// try dbox filter
		result := kData.Sort.Parse(SortParser).(bson.D)

		expectedPipe := bson.D{
			bson.DocElem{
				Name:  "foo",
				Value: -1,
			},
			bson.DocElem{
				Name:  "bar",
				Value: 1,
			},
			bson.DocElem{
				Name:  "_id",
				Value: -1,
			},
		}

		require.Equal(t, expectedPipe, result, "Result must same")
	}
}
