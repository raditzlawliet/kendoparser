![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/raditzlawliet/gokendoparser.svg?label=release)
![Codecov](https://img.shields.io/codecov/c/github/raditzlawliet/gokendoparser.svg)
[![](https://godoc.org/github.com/raditzlawliet/gokendoparser?status.svg)](http://godoc.org/github.com/raditzlawliet/gokendoparser)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/raditzlawliet/gokendoparser/issues)
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/raditzlawliet/gokendoparser)](https://goreportcard.com/report/github.com/raditzlawliet/gokendoparser) -->

# Go Kendo Parser
Your Golang Kendo parser, parsing Kendo data source request to golang struct immediately. Available parser to
- [eaciit/dbox Filter](https://github.com/eaciit/dbox)
- [eaciit/dbox Pipe or Aggregation type](https://github.com/eaciit/dbox)
- eaciit/dbflux Filter
- [mongo-go-driver Filter](https://github.com/mongodb/mongo-go-driver/)
- [xorm/builder Filter](https://gitea.com/xorm/builder/) This builder have nice ToSQL function, so you can transform into SQL 

## Features
- Convert kendo datasource request into go struct
- Basic Operator
- Transform filter + Chaining
- Plugable operator handler
- Parser Sort
- Custom pre-filter / before parse handler
- Extendable to any result

## Getting Started
Go get, import, use

### Parse Filter
#### JSON Sample
```json
{
    "data": {
        "filter": {
            "field": "id",
            "operator": "eq",
            "value": "val",
        }
    }
}
```

#### GO Implementation 
```go
import kpmgo "github.com/raditzlawliet/gokendoparser/adapter/mgo"

// just for information
// you can use gorilla mux or knot
var k = k *knot.WebContext

// retrieve payload
payload := &KendoRequest{}
k.GetPayload(payload)

// Usually payload the struct will be this
payload := KendoRequest {
    Data: KendoData{
        Filter: KendoFilter{
            Field: "_id", Operator: "eq", Value: "val"
        },
    },
}

filter := payload.Data.Filter.Parse(kpmongo.ParseFilter)
sort := payload.Data.Sort.Parse(kpmongo.ParseSort)
```

#### More JSON
```json
{
    "data": {
        "filters": [
            {
                "filter": {
                    "field": "id",
                    "operator": "eq",
                    "value": "val",
                }
            },
            {
                "filter": {
                    "field": "abc",
                    "operator": "in",
                    "values": ["a", "b"],
                }
            }
        ],
        "logic": "and"
    }
}
```
```go
import kpmgo "github.com/raditzlawliet/gokendoparser/adapter/mgo"
// just for information
// you can use gorilla mux or knot
var k = k *knot.WebContext

// retrieve payload
payload := &KendoRequest{}
k.GetPayload(payload)

// Usually payload the struct will be this
payload := KendoRequest {
    Data: KendoData{
        Filter: KendoFilter{
            Filters: []KendoFilter{
                Filter: KendoFilter{
                    Field: "_id", Operator: "eq", Value: "val"
                },
                Filter: KendoFilter{
                    Field: "abc", Operator: "in", Values: []interface{}{"a", "b"}
                },
            }
        },
    },
}

filter := payload.Data.Filter.Parse(kpmongo.ParseFilter)
sort := payload.Data.Sort.Parse(kpmongo.ParseSort)
```

### Extend & Hook custom operator handler
By default, package already registerd with basic operator such as 
- Equal
- Not Equal
- Contain
- Not Contain
- In
- Gte
- Lte
- Gte Date
- Lte Date
- Exists

But if you want to add custom operator that didn't exists yet, you can register in manually (for sample you can see [operator.go](adapter/mongo/operator.go)). Please ensure the return must be same with other operator each adapter. Each adapter have operatorManager that helping you to register into each adapter parser
```go
type Operator interface {
	Filter(KendoFilter) interface{}
}

type EqualOp struct{}

// Filter Filter
func (EqualOp) Filter(kf gokendoparser.KendoFilter) interface{} {
	if kf.IgnoreCase {
		return EqCi{kf.Field: kf.Value}
	}
	return builder.Eq{kf.Field: kf.Value}
}

func register() {
	operatorManager.RegisterOperator(equalOp, "eq", "equal")
}
```

### Transforming filter
Need modify your field? lowercase all field before processing? don't worry, you can use Transform to modify and apply to your all field. See [kendo_test.go](kendo_test.go)) for more uses
```go
import kpmgo "github.com/raditzlawliet/gokendoparser/adapter/mgo"

kendoFilter := KendoFilter{}

// transform filter field or all field in filters into lower case
kendoFilter.TransformField(strings.ToLower) // only current filter
kendoFilter.TransformFieldAll(strings.ToLower) // include filters

// custom transform include value (if needed), it use pointer KendoFilter, so you can direct modify. No need to return
transformIDMongo := func(kf *KendoFilter) {
    if kf.Field == "id" {
        kf.Field = "_id"
    }
}
kendoFilter.Transform(transformIDMongo) // only current filter
kendoFilter.TransformAll(transformIDMongo) // include filters

// chaining is possible
kendoFilter.TransformFieldAll(strings.ToLower).TransformAll(transformIDMongo).Parse(kpmongo.ParseFilter)

```

### Custom pre-filter
You can also add custom single handler before building filter by registered operator. This approach you can add custom direct filter within loop filter. pre-filter it self act same like parser, so you can chain the filter with other filter. But if your pre-filter already return value, the last parse only act return the value
```go
import kpmgo "github.com/raditzlawliet/gokendoparser/adapter/mgo"

beforeFilter := func (kf *KendoFilter) interface{} {
    return nil
}

// dbox filter
resultFilter := kendoFilter.TransformAllField(strings.ToLower).
    TransformAll(func(kf *KendoFilter) {
        if kf.Field == "id" {
            kf.Field = "_id"
        }
    }).
    BeforeParse(beforeFilter).
    Parse(kpmongo.ParseFilter)

// reset if needed another
kendoFilter.ResetPreFilter()

// dbox pipe
resultFilterPipe := kendoFilter.TransformAllField(strings.ToLower).
    TransformAll(func(kf *KendoFilter) {
        if kf.Field == "id" {
            kf.Field = "_id"
        }
    }).
    BeforeParse(beforeFilter).
    Parse(kpmongo.ParseFilter)
```

## Sort
do you need sort? You can do it easly.
```go
import kpmgo "github.com/raditzlawliet/gokendoparser/adapter/mgo"

sort := kData.Sort.Parse(kpmongo.ParseSort)
```

## Contribute
Feel free to contribute, don't forget to mention if needed

## License
MIT License 

## Author and Contributor
Radityo

