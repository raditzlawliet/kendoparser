package kpdboxpipe

import (
	"strings"
	"testing"
	"time"

	tk "github.com/eaciit/toolkit"
	"github.com/raditzlawliet/gokendoparser"
	"github.com/raditzlawliet/gokendoparser/helper"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func Test_ParseFilter(t *testing.T) {
	kendoFilter := gokendoparser.KendoFilter{
		Filters: []gokendoparser.KendoFilter{
			gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter := kendoFilter.Parse(ParseFilter).(tk.M)
	expectedFilter := tk.M{"$and": []tk.M{tk.M{"_id": tk.M{"$eq": "val"}}}}
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")

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
	resultFilter = kendoFilter.Parse(ParseFilter).(tk.M)
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
	kendoFilter = gokendoparser.KendoFilter{
		Filters: []gokendoparser.KendoFilter{
			gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
			gokendoparser.KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
			gokendoparser.KendoFilter{Field: "_id", Operator: "contains", Value: "val"},
			gokendoparser.KendoFilter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
			gokendoparser.KendoFilter{Field: "_id", Operator: "gte", Value: "val"},
			gokendoparser.KendoFilter{Field: "_id", Operator: "lte", Value: "val"},
			gokendoparser.KendoFilter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
			gokendoparser.KendoFilter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
			gokendoparser.KendoFilter{Field: "_id", Operator: "unknown", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter = kendoFilter.Parse(ParseFilter).(tk.M)
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
		kendoFilter := gokendoparser.KendoFilter{
			Filters: []gokendoparser.KendoFilter{
				gokendoparser.KendoFilter{
					Filters: []gokendoparser.KendoFilter{
						gokendoparser.KendoFilter{Field: "ID", Operator: "eq", Value: "val"},
						gokendoparser.KendoFilter{Field: "STATUS", Operator: "eq", Value: "true"},
					},
					Logic: "or",
				},
				gokendoparser.KendoFilter{
					Filters: []gokendoparser.KendoFilter{
						gokendoparser.KendoFilter{Field: "ID", Operator: "eq", Value: "val2"},
						gokendoparser.KendoFilter{Field: "ID", Operator: "neq", Value: "val2"},
					},
					Logic: "or",
				},
			},
			Logic: "and",
		}

		// try dbox pipe
		resultFilterPipe := kendoFilter.TransformAllField(strings.ToLower).
			TransformAll(func(kf *gokendoparser.KendoFilter) {
				if kf.Field == "id" {
					kf.Field = "_id"
				}
			}).
			BeforeParseAll(func(kf *gokendoparser.KendoFilter) interface{} {
				if kf.Field == "status" {
					// return your custom handler
					return tk.M{kf.Field: helper.StringToBool(kf.Value, false)}
				}
				return nil // pas nil to continue original filter
			}).Parse(ParseFilter).(tk.M)

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
		result := kData.Sort.Parse(ParseSort).(bson.D)

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
