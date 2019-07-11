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
- ~~Transform filter + Chaining~~
- ~~Extendable & hookable operator handler~~
- ~~local scope operator~~
- Parser Sort
- Custom filter on-the-fly
- Custom Sort on-the-fly
- Extendable & hookable database driver (MySQL, SQL, Oracle)

## Current Limitation!
- Only working for eaciit/dbox package
- Only working for MongoDB

## Getting Started
Easy to use like title said

### To eaciit/dbox filter
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

resultFilter := payload.Data.ToDboxFilter() 
// dbox.Eq("_id", "val")
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

resultFilter := payload.Data.ToDboxFilter() 
// dbox.And(
//     dbox.Eq("_id", "val"),
//     dbox.In("abc", []interface{}{"a", "b"}...),
// )
```

### To eaciit/dbox aggregation filter (return eaciit/toolkit/M)
Same like previously one

```go
resultFilter := payload.Data.ToDboxPipe() 
// tk.M{"$and": []tk.M{tk.M{"_id": tk.M{"$eq": "val"}}}}

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
	ToDboxPipe(KendoFilter) toolkit.M
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
func (BetweenOperator) ToDboxPipe(kf KendoFilter) toolkit.M {
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

#### Local scope Operator
You also can defined local scope only operator, just use the struct and register like global scope. It support nested or not
```go
betOp := BetweenOperator{}
kendoFilter := KendoFilter{}

// register operator in it struct only
kendoFilter.RegisterOperator("between", betOp)

// register operator in it struct and all child filters
kendoFilter.RegisterOperatorAll("between", betOp)
```

### Transforming filter
Need modify your field? lowercase all field before processing? don't worry, you can use Transform to modify and apply to your all field. See [kendo_test.go](kendo_test.go)) for more uses
```go
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
kendoFilter.TransformFieldAll(strings.ToLower).TransformAll(transformIDMongo).ToDboxFilter()
```

## Contribute
Feel free to contribute, don't forget to mention if needed

## License
MIT License 

## Author and Contributor
Radityo

