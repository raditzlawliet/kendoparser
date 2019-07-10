![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/raditzlawliet/gokendoparser.svg?label=release)
![Codecov](https://img.shields.io/codecov/c/github/raditzlawliet/gokendoparser.svg)
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/raditzlawliet/gokendoparser)](https://goreportcard.com/report/github.com/raditzlawliet/gokendoparser) -->

[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/raditzlawliet/gokendoparser/issues)

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
- Extendable & hookable operator handler
- Extendable & hookable database driver (MySQL, SQL, Oracle)

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

resultFilter := payload.Data.ToDboxFilter(func(s *string, filter *KendoFilter) *dbox.Filter {
    // extra function if you wish to modify the field it will loop each filter
    // Nullable / you can pass nil
    if *s == "id" {
        *s = "_id"
    }
    return nil
})
```

### To eaciit/dbox aggregation filter (return eaciit/toolkit/M)
Same like previously one

```go
resultFilter := payload.Data.ToAggregationFilter(func(s *string, filter *KendoFilter) tk.M {
    // extra function if you wish to modify the field it will loop each filter
    // Nullable / you can pass nil
    if *s == "id" {
        *s = "_id"
    }
    return nil
})

```
### To eaciit/dbflux filter (Coming soon)
Same like previously one
```go
```

## Contribute
Feel free to contribute, don't forget to mention if needed

## License
MIT License 

## Author and Contributor
Radit