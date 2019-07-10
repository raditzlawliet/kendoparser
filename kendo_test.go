package gokendoparser

import (
	"testing"
	"time"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"github.com/stretchr/testify/require"
)

func Test_ToDboxFilter(t *testing.T) {
	// Single filter
	{
		kendoFilter := KendoFilter{
			Field: "_id", Operator: "eq", Value: "val",
		}

		resultFilter := kendoFilter.ToDboxFilter()
		expectedFilter := dbox.Eq("_id", "val")
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := KendoRequest{}
		e := tk.UnjsonFromString(`{
			"data": {
				"filter": {
					"field": "_id",
					"operator": "eq",
					"value": "val"
				}
			}
		}`, &kendoRequest)
		require.Nil(t, e, "Json parse must work")
		require.Equal(t, kendoFilter, kendoRequest.Data.Filter, "Filter must same")
		resultFilterJSON := kendoRequest.Data.Filter.ToDboxFilter()
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}

	{
		kendoFilter := KendoFilter{
			Filters: []KendoFilter{
				KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.ToDboxFilter()
		expectedFilter := dbox.And(dbox.Eq("_id", "val"))
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	{
		kendoFilter := KendoFilter{
			Filters: []KendoFilter{
				KendoFilter{
					Filters: []KendoFilter{
						KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
						KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
					},
					Logic: "or",
				},
				KendoFilter{
					Filters: []KendoFilter{
						KendoFilter{Field: "_id", Operator: "eq", Value: "val2"},
						KendoFilter{Field: "_id", Operator: "neq", Value: "val2"},
					},
					Logic: "or",
				},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.ToDboxFilter()
		expectedFilter := dbox.And(
			dbox.Or(
				dbox.Eq("_id", "val"),
				dbox.Ne("_id", "val"),
			),
			dbox.Or(
				dbox.Eq("_id", "val2"),
				dbox.Ne("_id", "val2"),
			),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	// operator check
	{
		kendoFilter := KendoFilter{
			Filters: []KendoFilter{
				KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
				KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
				KendoFilter{Field: "_id", Operator: "doesnotcontain", Value: "val"},
				KendoFilter{Field: "_id", Operator: "contain", Value: "val"},
				KendoFilter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
				KendoFilter{Field: "_id", Operator: "gte", Value: "val"},
				KendoFilter{Field: "_id", Operator: "lte", Value: "val"},
				KendoFilter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
				KendoFilter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
				KendoFilter{Field: "_id", Operator: "unknown", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.ToDboxFilter()
		testTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
		expectedFilter := dbox.And(
			dbox.Eq("_id", "val"),
			dbox.Ne("_id", "val"),
			dbox.Eq("_id", tk.M{"$ne": tk.M{"$regex": `` + "val" + ``, "$options": "i"}}),
			dbox.Contains("_id", "val"),
			dbox.In("_id", []interface{}{"val"}...),
			dbox.Gte("_id", "val"),
			dbox.Lte("_id", "val"),
			dbox.Gte("time", testTime),
			dbox.Lte("time", testTime),
			dbox.Eq("_id", "val"),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}
}

func Test_ToPipeFilter(t *testing.T) {
	kendoFilter := KendoFilter{
		Filters: []KendoFilter{
			KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter := kendoFilter.ToAggregationFilter()
	expectedFilter := tk.M{"$and": []tk.M{tk.M{"_id": tk.M{"$eq": "val"}}}}
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")

	kendoFilter = KendoFilter{
		Filters: []KendoFilter{
			KendoFilter{
				Filters: []KendoFilter{
					KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
					KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
				},
				Logic: "or",
			},
			KendoFilter{
				Filters: []KendoFilter{
					KendoFilter{Field: "_id", Operator: "eq", Value: "val2"},
					KendoFilter{Field: "_id", Operator: "neq", Value: "val2"},
				},
				Logic: "or",
			},
		},
		Logic: "and",
	}
	resultFilter = kendoFilter.ToAggregationFilter()
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
	kendoFilter = KendoFilter{
		Filters: []KendoFilter{
			KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
			KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
			KendoFilter{Field: "_id", Operator: "contains", Value: "val"},
			KendoFilter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
			KendoFilter{Field: "_id", Operator: "gte", Value: "val"},
			KendoFilter{Field: "_id", Operator: "lte", Value: "val"},
			KendoFilter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
			KendoFilter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
			KendoFilter{Field: "_id", Operator: "unknown", Value: "val"},
		},
		Logic: "and",
	}
	resultFilter = kendoFilter.ToAggregationFilter()
	testTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
	expectedFilter = tk.M{"$and": []tk.M{
		tk.M{"_id": tk.M{"$eq": "val"}},
		tk.M{"_id": tk.M{"$ne": "val"}},
		tk.M{"_id": RegexContains("val", false)},
		tk.M{"_id": tk.M{"$in": []interface{}{"val"}}},
		tk.M{"_id": tk.M{"$gte": "val"}},
		tk.M{"_id": tk.M{"$lte": "val"}},
		tk.M{"time": tk.M{"$gte": testTime}},
		tk.M{"time": tk.M{"$lte": testTime}},
		tk.M{"_id": tk.M{"$eq": "val"}},
	}}
	require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
}
