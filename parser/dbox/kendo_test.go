package kpdbox

import (
	"strings"
	"testing"
	"time"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"github.com/raditzlawliet/kendoparser"
	"github.com/raditzlawliet/kendoparser/helper"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func Test_FilterParser(t *testing.T) {
	// Single filter
	{
		kendoFilter := kendoparser.Filter{
			Field: "_id", Operator: "eq", Value: "val",
		}

		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
		expectedFilter := dbox.Eq("_id", "val")
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := kendoparser.KendoRequest{}
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
		resultFilterJSON := kendoRequest.Data.Filter.Parse(FilterParser).(*dbox.Filter)
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}

	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
		expectedFilter := dbox.And(dbox.Eq("_id", "val"))
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	{
		kendoFilter := kendoparser.Filter{
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
		resultFilter := kendoFilter.Parse(FilterParser)
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
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{Field: "_id", Operator: "eq", Value: "val"},
				kendoparser.Filter{Field: "_id", Operator: "neq", Value: "val"},
				kendoparser.Filter{Field: "_id", Operator: "doesnotcontain", Value: "val"},
				kendoparser.Filter{Field: "_id", Operator: "contain", Value: "val"},
				kendoparser.Filter{Field: "_id", Operator: "in", Values: []interface{}{"val"}},
				kendoparser.Filter{Field: "_id", Operator: "gte", Value: "val"},
				kendoparser.Filter{Field: "_id", Operator: "lte", Value: "val"},
				kendoparser.Filter{Field: "time", Operator: "gtedate", Value: "2006-01-02T15:04:05Z07:00"},
				kendoparser.Filter{Field: "time", Operator: "ltedate", Value: "2006-01-02T15:04:05Z07:00"},
				kendoparser.Filter{Field: "_id", Operator: "unknown", Value: "val"},
				kendoparser.Filter{Field: "_id", Operator: "exists", Value: "val"},
			},
			Logic: "and",
		}
		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
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
			&dbox.Filter{
				Field: "_id",
				Op:    dbox.FilterOpEqual,
				Value: tk.M{
					"$exists": helper.StringToBool("val", false),
				},
			},
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}
}
func Test_TransformField(t *testing.T) {
	// transform single filter field
	// _ID => _id
	{
		kendoFilter := kendoparser.Filter{
			Field: "_ID", Operator: "eq", Value: "val",
		}
		kendoFilter.TransformField(strings.ToLower)

		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
		expectedFilter := dbox.Eq("_id", "val")
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := kendoparser.KendoRequest{}
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
		resultFilterJSON := kendoRequest.Data.Filter.Parse(FilterParser).(*dbox.Filter)
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}

	// test Transform single field, should not affect the child filter
	// _ID => _id
	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "_ID", Operator: "eq", Value: "val"},
						kendoparser.Filter{Field: "_ID", Operator: "neq", Value: "val"},
					},
					Logic: "or",
				},
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "_ID", Operator: "eq", Value: "val2"},
						kendoparser.Filter{Field: "_ID", Operator: "neq", Value: "val2"},
					},
					Logic: "or",
				},
			},
			Logic: "and",
		}
		kendoFilter.TransformField(strings.ToLower)
		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
		expectedFilter := dbox.And(
			dbox.Or(
				dbox.Eq("_ID", "val"),
				dbox.Ne("_ID", "val"),
			),
			dbox.Or(
				dbox.Eq("_ID", "val2"),
				dbox.Ne("_ID", "val2"),
			),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	// test transform all field
	// _ID => _id
	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "_ID", Operator: "eq", Value: "val"},
						kendoparser.Filter{Field: "_ID", Operator: "neq", Value: "val"},
					},
					Logic: "or",
				},
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "_ID", Operator: "eq", Value: "val2"},
						kendoparser.Filter{Field: "_ID", Operator: "neq", Value: "val2"},
					},
					Logic: "or",
				},
			},
			Logic: "and",
		}
		kendoFilter.TransformAllField(strings.ToLower)
		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
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
}
func Test_Transform(t *testing.T) {
	// transform single filter
	// ID => _id
	{
		kendoFilter := kendoparser.Filter{
			Field: "ID", Operator: "eq", Value: "val",
		}
		kendoFilter.TransformField(strings.ToLower)
		kendoFilter.Transform(func(kf *kendoparser.Filter) {
			if kf.Field == "id" {
				kf.Field = "_id"
			}
		})

		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
		expectedFilter := dbox.Eq("_id", "val")
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := kendoparser.KendoRequest{}
		e := tk.UnjsonFromString(`{
			"data": {
				"filter": {
					"field": "ID",
					"operator": "eq",
					"value": "val"
				}
			}
		}`, &kendoRequest)
		kendoRequest.Data.Filter.TransformField(strings.ToLower)
		kendoRequest.Data.Filter.Transform(func(kf *kendoparser.Filter) {
			if kf.Field == "id" {
				kf.Field = "_id"
			}
		})

		require.Nil(t, e, "Json parse must work")
		require.Equal(t, kendoFilter, kendoRequest.Data.Filter, "Filter must same")
		resultFilterJSON := kendoRequest.Data.Filter.Parse(FilterParser).(*dbox.Filter)
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}

	// test Transform single, should not affect the child filter
	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "ID", Operator: "eq", Value: "val"},
						kendoparser.Filter{Field: "ID", Operator: "neq", Value: "val"},
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
		kendoFilter.TransformField(strings.ToLower)
		kendoFilter.Transform(func(kf *kendoparser.Filter) {
			if kf.Field == "id" {
				kf.Field = "_id"
			}
		})
		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
		expectedFilter := dbox.And(
			dbox.Or(
				dbox.Eq("ID", "val"),
				dbox.Ne("ID", "val"),
			),
			dbox.Or(
				dbox.Eq("ID", "val2"),
				dbox.Ne("ID", "val2"),
			),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	// test transform all
	// ID => _id
	{
		kendoFilter := kendoparser.Filter{
			Filters: []kendoparser.Filter{
				kendoparser.Filter{
					Filters: []kendoparser.Filter{
						kendoparser.Filter{Field: "ID", Operator: "eq", Value: "val"},
						kendoparser.Filter{Field: "ID", Operator: "neq", Value: "val"},
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
		kendoFilter.TransformAllField(strings.ToLower)
		kendoFilter.TransformAll(func(kf *kendoparser.Filter) {
			if kf.Field == "id" {
				kf.Field = "_id"
			}
		})
		resultFilter := kendoFilter.Parse(FilterParser).(*dbox.Filter)
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
		// try dbox filter
		resultFilter := kendoFilter.TransformAllField(strings.ToLower).
			TransformAll(func(kf *kendoparser.Filter) {
				if kf.Field == "id" {
					kf.Field = "_id"
				}
			}).
			AddAllParser(func(kf *kendoparser.Filter) interface{} {
				if kf.Field == "status" {
					// return your custom handler
					return dbox.Eq(kf.Field, helper.StringToBool(cast.ToString(kf.Value), false))
				}
				return nil // pas nil to continue original filter
			}).
			Parse(FilterParser).(*dbox.Filter)

		// reset if needed another
		kendoFilter.ResetAllAdditionalParsers()

		expectedFilter := dbox.And(
			dbox.Or(
				dbox.Eq("_id", "val"),
				dbox.Eq("status", true),
			),
			dbox.Or(
				dbox.Eq("_id", "val2"),
				dbox.Ne("_id", "val2"),
			),
		)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
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
		result := kData.Sort.Parse(SortParser).([]string)
		expected := []string{"-foo", "bar", "-_id"}
		require.Equal(t, expected, result, "Result must same")
	}
}
