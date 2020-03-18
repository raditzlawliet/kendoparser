package kendoparser

import (
	"fmt"
	"github.com/eaciit/cast"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"testing"
)

func Test_All(t *testing.T) {
	// Operator Manager
	om := OperatorManager{}
	equal := func(kf Filter) interface{} {
		if kf.IgnoreCase {
			value := regexp.QuoteMeta(cast.ToString(kf.Value))
			return bson.M{kf.Field: primitive.Regex{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
		}
		return bson.M{kf.Field: bson.M{"$eq": kf.Value}}
	}
	om.SetDefaultOperator(equal)
	om.RegisterOperator(equal, "eq")
	require.Equal(t, 1, len(om.OperatorFilters), "Must same")
	om.Reset()
	require.Equal(t, 0, len(om.OperatorFilters), "Must same")

	// Sort
	payloadSort := Sort{
		SortDetail{
			Field: "foo",
			Dir:   "DESC",
		},
		SortDetail{
			Field: "bar",
			Dir:   "ASC",
		},
		SortDetail{
			Field: "_id",
			Dir:   "desc",
		},
	}
	// testing kendo sort
	sortParser := func(s *Sort) interface{} {
		sorter := []string{}
		for _, ks := range *s {
			sort := 1
			if strings.ToLower(ks.Dir) == "desc" {
				sort = -1
			}
			sorter = append(sorter, fmt.Sprintf("%v=%v", ks.Field, sort))
		}
		return sorter
	}

	require.Equal(t, []string{"foo=-1", "bar=1", "_id=-1"}, payloadSort.Parse(sortParser), "Must same")
}
