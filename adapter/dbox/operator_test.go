package kpdbox

import (
	"testing"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"github.com/raditzlawliet/gokendoparser"
	"github.com/stretchr/testify/require"
)

func Test_OperatorHook(t *testing.T) {
	// Single filter
	betOp := BetweenOp{}
	operatorManager.RegisterOperator(betOp, "between_custom")

	// testing eq / global
	{
		kendoFilter := gokendoparser.KendoFilter{
			Field: "_id", Operator: "eq", Value: "val",
		}

		resultFilter := kendoFilter.Parse(ParseFilter).(*dbox.Filter)
		expectedFilter := dbox.Eq("_id", "val")
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
		resultFilterJSON := kendoRequest.Data.Filter.Parse(ParseFilter).(*dbox.Filter)
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}

	// testing custom registered between
	{
		kendoFilter := gokendoparser.KendoFilter{
			Field: "v", Operator: "between_custom", Values: []interface{}{"1", "2"},
		}

		resultFilter := kendoFilter.Parse(ParseFilter).(*dbox.Filter)
		expectedFilter := dbox.And(dbox.Gte("v", "1"), dbox.Lte("v", "2"))
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

		kendoRequest := gokendoparser.KendoRequest{}
		e := tk.UnjsonFromString(`{
			"data": {
				"filter": {
					"field": "v",
					"operator": "between_custom",
					"values": ["1", "2"]
				}
			}
		}`, &kendoRequest)
		require.Nil(t, e, "Json parse must work")
		require.Equal(t, kendoFilter, kendoRequest.Data.Filter, "Filter must same")
		resultFilterJSON := kendoRequest.Data.Filter.Parse(ParseFilter).(*dbox.Filter)
		require.Equal(t, expectedFilter, resultFilterJSON, "Result dbox filter must same")
	}
}

func Test_OperatorHookLocalScope(t *testing.T) {
	operatorManager.Reset()
	operatorManager.SetDefaultOperator(EqualOp{})
	// testing custom registered between but not registered yet
	{
		kendoFilter := gokendoparser.KendoFilter{
			Field: "v", Operator: "between_custom", Values: []interface{}{"1", "2"},
		}

		resultFilter := kendoFilter.Parse(ParseFilter).(*dbox.Filter)
		expectedFilter := dbox.Eq("v", "") // because between not registered YET
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")

	}

	operatorManager.Reset()
	operatorManager.SetDefaultOperator(EqualOp{})
	{
		// Single filter
		betOp := BetweenOp{}

		kendoFilter := gokendoparser.KendoFilter{
			Field: "v", Operator: "between_custom", Values: []interface{}{"1", "2"},
		}

		resultFilter := kendoFilter.RegisterOperator(betOp, "between_custom").Parse(ParseFilter).(*dbox.Filter)
		expectedFilter := dbox.And(dbox.Gte("v", "1"), dbox.Lte("v", "2"))
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	operatorManager.Reset()
	operatorManager.SetDefaultOperator(EqualOp{})
	{
		// Single filter
		betOp := BetweenOp{}

		kendoFilter := gokendoparser.KendoFilter{
			Filters: []gokendoparser.KendoFilter{
				gokendoparser.KendoFilter{
					Field: "v", Operator: "between_custom", Values: []interface{}{"1", "2"},
				},
				gokendoparser.KendoFilter{
					Field: "v", Operator: "between_custom", Values: []interface{}{"1", "2"},
				},
			},
			Logic: "and",
		}

		resultFilter := kendoFilter.RegisterOperatorAll(betOp, "between_custom").Parse(ParseFilter).(*dbox.Filter)
		expectedFilter :=
			dbox.And(
				dbox.And(dbox.Gte("v", "1"), dbox.Lte("v", "2")),
				dbox.And(dbox.Gte("v", "1"), dbox.Lte("v", "2")),
			)
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

	operatorManager.Reset()
	operatorManager.SetDefaultOperator(EqualOp{})
	{
		// Single filter
		betOp := BetweenOp{}

		kendoFilter := gokendoparser.KendoFilter{
			Field: "v", Operator: "between_custom", Values: []interface{}{"1", "2"},
		}

		resultFilter := kendoFilter.RegisterOperator(betOp, "between_custom").Parse(ParseFilter).(*dbox.Filter)
		expectedFilter := dbox.And(dbox.Gte("v", "1"), dbox.Lte("v", "2"))
		require.Equal(t, expectedFilter, resultFilter, "Result dbox filter must same")
	}

}
