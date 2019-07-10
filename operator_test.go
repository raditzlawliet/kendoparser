package gokendoparser

import (
	"testing"

	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
	tk "github.com/eaciit/toolkit"
	"github.com/stretchr/testify/require"
)

func Test_OperatorHook(t *testing.T) {
	// Single filter
	betOp := BetweenOperator{}
	RegisterOperator("between", betOp)

	// testing eq / global
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

	// testing custom registered between
	{
		kendoFilter := KendoFilter{
			Field: "v", Operator: "between", Values: []interface{}{"1", "2"},
		}

		resultFilter := kendoFilter.ToDboxFilter()
		expectedFilter := dbox.And(dbox.Gte("v", "1"), dbox.Lte("v", "2"))
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := KendoRequest{}
		e := tk.UnjsonFromString(`{
			"data": {
				"filter": {
					"field": "v",
					"operator": "between",
					"values": ["1", "2"]
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
			Field: "v", Operator: "between", Values: []interface{}{"1", "2"},
		}

		resultFilter := kendoFilter.ToAggregationFilter()
		expectedFilter := toolkit.M{"v": toolkit.M{"$gte": "1", "$lte": "2"}}
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := KendoRequest{}
		e := tk.UnjsonFromString(`{
			"data": {
				"filter": {
					"field": "v",
					"operator": "between",
					"values": ["1", "2"]
				}
			}
		}`, &kendoRequest)
		require.Nil(t, e, "Json parse must work")
		require.Equal(t, kendoFilter, kendoRequest.Data.Filter, "Filter must same")
		resultFilterJSON := kendoRequest.Data.Filter.ToAggregationFilter()
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}
}
