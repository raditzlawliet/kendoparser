package kpdbflex

import (
	"strings"
	"testing"
	"time"

	"git.eaciitapp.com/sebar/dbflex"

	tk "github.com/eaciit/toolkit"
	"github.com/raditzlawliet/gokendoparser"
	"github.com/stretchr/testify/require"
)

func Test_ParseFilter(t *testing.T) {
	// Single filter
	{
		kendoFilter := gokendoparser.KendoFilter{
			Field: "_id", Operator: "eq", Value: "val",
		}

		resultFilter := kendoFilter.Parse(Parser{}).(*dbflex.Filter)
		expectedFilter := dbflex.Eq("_id", "val")
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := gokendoparser.KendoRequest{}
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
		resultFilterJSON := kendoRequest.Data.Filter.Parse(Parser{}).(*dbflex.Filter)
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}

	{
		kendoFilter := gokendoparser.KendoFilter{
			Filters: []gokendoparser.KendoFilter{
				gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.Parse(Parser{}).(*dbflex.Filter)
		expectedFilter := dbflex.And(dbflex.Eq("_id", "val"))
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	{
		kendoFilter := gokendoparser.KendoFilter{
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
		resultFilter := kendoFilter.Parse(Parser{})
		expectedFilter := dbflex.And(
			dbflex.Or(
				dbflex.Eq("_id", "val"),
				dbflex.Ne("_id", "val"),
			),
			dbflex.Or(
				dbflex.Eq("_id", "val2"),
				dbflex.Ne("_id", "val2"),
			),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	// operator check
	{
		kendoFilter := gokendoparser.KendoFilter{
			Filters: []gokendoparser.KendoFilter{
				gokendoparser.KendoFilter{Field: "_id", Operator: "eq", Value: "val"},
				gokendoparser.KendoFilter{Field: "_id", Operator: "neq", Value: "val"},
				// gokendoparser.KendoFilter{Field: "_id", Operator: "doesnotcontain", Value: "val"},
				gokendoparser.KendoFilter{Field: "_id", Operator: "contain", Value: "val"},
				gokendoparser.KendoFilter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
				gokendoparser.KendoFilter{Field: "_id", Operator: "gte", Value: "val"},
				gokendoparser.KendoFilter{Field: "_id", Operator: "lte", Value: "val"},
				gokendoparser.KendoFilter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
				gokendoparser.KendoFilter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
				gokendoparser.KendoFilter{Field: "_id", Operator: "unknown", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.Parse(Parser{}).(*dbflex.Filter)
		testTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
		expectedFilter := dbflex.And(
			dbflex.Eq("_id", "val"),
			dbflex.Ne("_id", "val"),
			// dbflex.Eq("_id", tk.M{"$ne": tk.M{"$regex": `` + "val" + ``, "$options": "i"}}),
			dbflex.Contains("_id", "val"),
			dbflex.In("_id", []interface{}{"val"}...),
			dbflex.Gte("_id", "val"),
			dbflex.Lte("_id", "val"),
			dbflex.Gte("time", testTime),
			dbflex.Lte("time", testTime),
			dbflex.Eq("_id", "val"),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result filter must same")
	}
}
func Test_TransformField(t *testing.T) {
	// transform single filter field
	// _ID => _id
	{
		kendoFilter := gokendoparser.KendoFilter{
			Field: "_ID", Operator: "eq", Value: "val",
		}
		kendoFilter.TransformField(strings.ToLower)

		resultFilter := kendoFilter.Parse(Parser{}).(*dbflex.Filter)
		expectedFilter := dbflex.Eq("_id", "val")
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := gokendoparser.KendoRequest{}
		e := tk.UnjsonFromString(`{
			"data": {
				"filter": {
					"field": "_ID",
					"operator": "eq",
					"value": "val"
				}
			}
		}`, &kendoRequest)
		kendoRequest.Data.Filter.TransformField(strings.ToLower)

		require.Nil(t, e, "Json parse must work")
		require.Equal(t, kendoFilter, kendoRequest.Data.Filter, "Filter must same")
		resultFilterJSON := kendoRequest.Data.Filter.Parse(Parser{}).(*dbflex.Filter)
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}

	// test Transform single field, should not affect the child filter
	// _ID => _id
	{
		kendoFilter := gokendoparser.KendoFilter{
			Filters: []gokendoparser.KendoFilter{
				gokendoparser.KendoFilter{
					Filters: []gokendoparser.KendoFilter{
						gokendoparser.KendoFilter{Field: "_ID", Operator: "eq", Value: "val"},
						gokendoparser.KendoFilter{Field: "_ID", Operator: "neq", Value: "val"},
					},
					Logic: "or",
				},
				gokendoparser.KendoFilter{
					Filters: []gokendoparser.KendoFilter{
						gokendoparser.KendoFilter{Field: "_ID", Operator: "eq", Value: "val2"},
						gokendoparser.KendoFilter{Field: "_ID", Operator: "neq", Value: "val2"},
					},
					Logic: "or",
				},
			},
			Logic: "and",
		}
		kendoFilter.TransformField(strings.ToLower)
		resultFilter := kendoFilter.Parse(Parser{}).(*dbflex.Filter)
		expectedFilter := dbflex.And(
			dbflex.Or(
				dbflex.Eq("_ID", "val"),
				dbflex.Ne("_ID", "val"),
			),
			dbflex.Or(
				dbflex.Eq("_ID", "val2"),
				dbflex.Ne("_ID", "val2"),
			),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	// test transform all field
	// _ID => _id
	{
		kendoFilter := gokendoparser.KendoFilter{
			Filters: []gokendoparser.KendoFilter{
				gokendoparser.KendoFilter{
					Filters: []gokendoparser.KendoFilter{
						gokendoparser.KendoFilter{Field: "_ID", Operator: "eq", Value: "val"},
						gokendoparser.KendoFilter{Field: "_ID", Operator: "neq", Value: "val"},
					},
					Logic: "or",
				},
				gokendoparser.KendoFilter{
					Filters: []gokendoparser.KendoFilter{
						gokendoparser.KendoFilter{Field: "_ID", Operator: "eq", Value: "val2"},
						gokendoparser.KendoFilter{Field: "_ID", Operator: "neq", Value: "val2"},
					},
					Logic: "or",
				},
			},
			Logic: "and",
		}
		kendoFilter.TransformAllField(strings.ToLower)
		resultFilter := kendoFilter.Parse(Parser{}).(*dbflex.Filter)
		expectedFilter := dbflex.And(
			dbflex.Or(
				dbflex.Eq("_id", "val"),
				dbflex.Ne("_id", "val"),
			),
			dbflex.Or(
				dbflex.Eq("_id", "val2"),
				dbflex.Ne("_id", "val2"),
			),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
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
		result := kData.Sort.Parse(Parser{}).([]string)
		expected := []string{"-foo", "bar", "-_id"}
		require.Equal(t, expected, result, "Result must same")
	}
}
