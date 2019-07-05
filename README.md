# Go Kendo Parser
Golang Kendo parser, parsing Kendo data source request to golang struct immediately to eaciit/dbox filter or aggregate filter. Available parser to
- [eaciit/dbox Filter](https://github.com/eaciit/dbox)
- [eaciit/dbox Pipe or Aggregation type](https://github.com/eaciit/dbox)
- eaciit/dbflux Filter (Coming Soon)

Easy to use like title said
## To eaciit/dbox filter
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

## To eaciit/dbox aggregation filter (return eaciit/toolkit/M)
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
## To eaciit/dbflux filter (Coming soon)
Same like previously one
```go
```
