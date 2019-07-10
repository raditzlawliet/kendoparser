![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/raditzlawliet/gokendoparser.svg?label=release)
![Codecov](https://img.shields.io/codecov/c/github/raditzlawliet/gokendoparser.svg)
[![](https://godoc.org/github.com/raditzlawliet/gokendoparser?status.svg)](http://godoc.org/github.com/raditzlawliet/gokendoparser)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/raditzlawliet/gokendoparser/issues)
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/raditzlawliet/gokendoparser)](https://goreportcard.com/report/github.com/raditzlawliet/gokendoparser) -->

# Go Kendo Parser
Your Golang Kendo parser, parsing Kendo data source request to golang struct immediately to eaciit/dbox filter or aggregate filter for **Mongo DB Driver**. Available parser to
- [eaciit/dbox Filter](https://github.com/eaciit/dbox)
- [eaciit/dbox Pipe or Aggregation type](https://github.com/eaciit/dbox)
- eaciit/dbflux Filter (Coming Soon)

## Features
- ~~Convert kendo datasource request into go struct~~
- ~~Basic Operator~~
- ~~Dbox filter & aggregate support for Mongo DB Driver~~
- Extendable & hookable filter (global / struct only)
- ~~Extendable & hookable operator handler~~
- Extendable & hookable database driver (MySQL, SQL, Oracle)
- Parser Sort

## Current Limitation!
- Only working for eaciit/dbox package
- Only working for MongoDB

## Getting Started
Easy to use like title said

### To eaciit/dbox filter
```json
{
    data: {
        filter: {
            field: "id",
            operator: "eq",
            value: "val"
        }
    }
}
```

```go
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
            Field: "id", Operator: "eq", Value: "val",
        },
    },
}

resultFilter := payload.Data.ToDboxFilter()
```

### To eaciit/dbox aggregation filter (return eaciit/toolkit/M)
Same like previously one

```go
resultFilter := payload.Data.ToAggregationFilter()

```

### To eaciit/dbflux filter (Coming soon)
Same like previously one
```go
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

But if you want to add custom operator that didn't exists yet, you can register in global handler (for sample you can see [operator_between.go](operator_between.go)). You must implement all function in Operator interface (or just return nil if you dont want implement other hook)
```go
type Operator interface {
	ToDboxFilter(KendoFilter) *dbox.Filter
	ToAggregationFilter(KendoFilter) toolkit.M
}
```

```go
// extend struct from interface Operator, all interface func must appear
type BetweenOperator struct {}

func (BetweenOperator) ToDboxFilter(kf KendoFilter) *dbox.Filter {
	var v0, v1 interface{}
	if len(kf.Values) > 0 {
		v0 = kf.Values[0]
	}
	if len(kf.Values) > 1 {
		v1 = kf.Values[1]
	}
	return dbox.And(dbox.Gte(kf.Field, v0), dbox.Lte(kf.Field, v1))
}
func (BetweenOperator) ToAggregationFilter(kf KendoFilter) toolkit.M {
    return nil // pass whatever if you dont want to implement
    // return DefaultOperator // or pass default
}
```

```go
// register it 
betOp := BetweenOperator{}
RegisterOperator("between", betOp)

// overwrite Default Operator
SetDefaultOperator(betOp)
```

## Contribute
Feel free to contribute, don't forget to mention if needed

## License
MIT License 

## Author and Contributor
Radityo