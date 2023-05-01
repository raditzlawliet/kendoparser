> **Warning**
> Project will be archived as read-only. This project still can be use, But there will be no update until unknown date announcement. 
 
 
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/raditzlawliet/kendoparser.svg?label=release)
![Codecov](https://img.shields.io/codecov/c/github/raditzlawliet/kendoparser.svg)
[![](https://godoc.org/github.com/raditzlawliet/kendoparser?status.svg)](http://godoc.org/github.com/raditzlawliet/kendoparser)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/raditzlawliet/kendoparser/issues)
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/raditzlawliet/kendoparser)](https://goreportcard.com/report/github.com/raditzlawliet/kendoparser) -->

# Golang Kendo Parser (v 1.*)
Parsing Kendo DataSource filter & sort request to any kind of data in go that can immediately serve to mongo query or aggregate and other ORMs. Already used in EACIIT's environment and can be scale-able for generic purpose. Currently there are few parser available:
- [mongo-go-driver Filter](https://github.com/mongodb/mongo-go-driver/) (mainly used for primary test)
- [eaciit/dbox Filter](https://github.com/eaciit/dbox)
- [eaciit/dbox Pipe or Aggregation type](https://github.com/eaciit/dbox)
- sebar/dbflex Filter
- [xorm/builder Filter](https://gitea.com/xorm/builder/) This builder have nice ToSQL function, so you can transform into SQL 

## Features
- Parse Kendo DataSource request into go struct
- Parse Kendo DataSource request into filter & sort database / orm-ish
- Transforming filter + Chaining
- Customizable filter operator
- Additional Pre-parse filter
- Adaptor based, customizable parse into other generic purpose
- Kendo Request Wrapper
- Go mod support

## Getting Started
- Install
```sh
go get -v raditzlawliet/kendoparser
```

- Import
```go
import (
    kp "github.com/raditzlawliet/kendoparser"
    kpmongo "github.com/raditzlawliet/kendoparser/parser/kpmongo"
)
```
- Use 
```go
request := struct {
    Filter kp.Filter,
    Sort kp.Sort,
}{}

filter := request.Filter.Parse(kpmongo.FilterParser)
sort := request.Sort.Parse(kpmongo.SortParser)
```
#### How to use Parse Filter
- JSON Sample
```json
{
    "data": {
        "filter": {
            "field": "id",
            "operator": "eq",
            "value": "val"
        }
    }
}
```

- GO Implementation 
```go
import kpmongo "github.com/raditzlawliet/kendoparser/parser/kpmongo"
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
            Field: "_id", Operator: "eq", Value: "val"
        },
    },
}

filter := payload.Data.Filter.Parse(kpmongo.FilterParser)
sort := payload.Data.Sort.Parse(kpmongo.SortParser)
```

- More Example
```json
{
    "data": {
        "filters": [
            {
                "filter": {
                    "field": "id",
                    "operator": "eq",
                    "value": "val"
                }
            },
            {
                "filter": {
                    "field": "abc",
                    "operator": "in",
                    "values": ["a", "b"]
                }
            }
        ],
        "logic": "and"
    }
}
```
```go
import kpmongo "github.com/raditzlawliet/kendoparser/parser/kpmongo"
// just for information
// you can use gorilla mux or knot
var k = k *knot.WebContext

// retrieve payload
payload := &KendoRequest{}
k.GetPayload(payload)

// Usually payload the struct will be like this
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

filter := payload.Data.Filter.Parse(kpmongo.FilterParser)
sort := payload.Data.Sort.Parse(kpmongo.SortParser)
```
#### Available Parser & Operator Manager - Filters
- [kpmongo - mongo-go-driver](./parser/mongo)
    - equal, notequal, equaldate, notequaldate, contain, notcontain, in, notin, gt, lt, gte, lte, gtdate, ltdate, gtedate, ltedate, exists, dateexists, between
- [kpdbox - eaciit/dbox](./parser/dbox)
    - equal, notequal, contain, notcontain, in, gt, lt, gte, lte, gtedate, ltedate, exists, between
- [kpdboxpipe - eaciit/dbox](./parser/dboxpipe)
    - equal, notequal, contain, notcontain, in, gt, lt, gte, lte, gtedate, ltedate, exists, between
- [kpdbflex - sebar/dbflex](./parser/dbflex)
    - equal, notequal, contain, notcontain, in, gt, lt, gte, lte, gtedate, ltedate, exists, between
- [kpxorm - xorm/builder](./parser/xorm)
    - equal, notequal, contain, notcontain, in, gt, lt, gte, lte, gtedate, ltedate, exists, between

#### Operator Manager
_Operator Manager_ mainly used for filter parser to parsing the data based on operator. In available parser they also have each _Operator Manager_. 
You also can create your own _Operator Manager_ and you can also set scoped _Operator Manager_ for spesific payload

- Extending Existing Operator Manager

```go
import kpmongo "github.com/raditzlawliet/kendoparser/parser/kpmongo"

// Create the func
func IsOne(kf kendoparser.Filter) interface{} {
    return bson.M{kf.Field: bson.M{"$eq": 1}}
}

// Assign into Exisiting Operator Manager
kpmongo.OperatorManager.RegisterOperator(IsOne, "isone")

// Example kendo json request
"filter": {
    "field": "id",
    "operator": "isone"
}

// Parse
request := struct {
    Filter kp.Filter,
}{}
filter := request.Filter.Parse(kpmongo.FilterParser)

// Result
{"id" : {"$eq": 1}}
``` 

- Creating Your own Operator Manager

    Usually when you try to create your own Operator Manager, you should create your own Parser. if the parser already existed, i recommended you use the existing Operator Manager.
    
```go
// creating global package variable
var (
	// OperatorManager of Mongo Parser
	OperatorManager = new(kendoparser.OperatorManager)
	// Operator bundle of Mongo Parser
	Operator = OperatorBundle{}
)

// Wrapper Bundle, so you can put all the func filter into this struct
type OperatorBundle struct{}

// Register when import
func init() {
	RegisterOperator()
}

// Registration Process
func RegisterOperator() {
	OperatorManager.SetDefaultOperator(Operator.Equal)
	OperatorManager.RegisterOperator(Operator.Equal, "eq", "equal")
}

// Some Filter 
func (o *OperatorBundle) Equal(kf kendoparser.Filter) interface{} {
	if value, ok := kf.Value.(string); ok && kf.IgnoreCase {
		value := regexp.QuoteMeta(value)
		return bson.M{kf.Field: primitive.Regex{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}}
	}
	return bson.M{kf.Field: bson.M{"$eq": kf.Value}}
}
```

#### Transforming filter
Need modify your field? lowercase all field before processing? don't worry, you can use Transform to modify and apply to your all field. See [kendo_test.go](kendo_test.go)) for more uses
```go
import kpmongo "github.com/raditzlawliet/kendoparser/parser/kpmongo"

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
kendoFilter.TransformFieldAll(strings.ToLower).TransformAll(transformIDMongo).Parse(kpmongo.FilterParser)

```

#### Additional pre-parse filter
You can also add additional single handler before building filter by registered operator. This approach you can add custom direct filter within loop filter. pre-filter it self act same like parser, so you can chain the filter with other filter. But if your pre-filter already return value, the last parse only act return the value
```go
import kpmongo "github.com/raditzlawliet/kendoparser/parser/kpmongo"

beforeParser := func (kf *KendoFilter) interface{} {
    return nil
}

// dbox filter
resultFilter := kendoFilter.TransformAllField(strings.ToLower).
    TransformAll(func(kf *KendoFilter) {
        if kf.Field == "id" {
            kf.Field = "_id"
        }
    }).
    AddParse(beforeParser).
    Parse(kpmongo.FilterParser)

// reset if needed another
kendoFilter.ResetAdditionalParsers()
kendoFilter.ResetAllAdditionalParsers() // Recursive to all filter childs

// dbox pipe
resultFilterPipe := kendoFilter.TransformAllField(strings.ToLower).
    TransformAll(func(kf *KendoFilter) {
        if kf.Field == "id" {
            kf.Field = "_id"
        }
    }).
    AddParse(beforeParser).
    Parse(kpmongo.FilterParser)
```

#### Sort
do you need sort? You can do it easly.
```go
import kpmongo "github.com/raditzlawliet/kendoparser/parser/kpmongo"

sort := kData.Sort.Parse(kpmongo.SortParser)
```

## FAQ
- go mod tidy keep error `git.eaciitapp.com/sebar/dbflex: no matching versions for query "latest"`

    This error appear because I use some private database driver in this parser, but you can exclude this using go.mod.
    
    Modify go.mod file, if you dont have the package, just use this dummy in your go.mod, AND If you have the package and using it, replace it into your correct path package:
    
    ```replace git.eaciitapp.com/sebar/dbflex => ./pkg/git.eaciitapp.com/sebar/dbflex```


## Contribute
Feel free to contribute, don't forget to mention if needed

## License
MIT License 

## Author and Contributor
Radityo 
