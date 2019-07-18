# rester
[![Build Status](https://travis-ci.com/hoenirvili/rester.svg?branch=master)](https://travis-ci.com/hoenirvili/rester) [![GoDoc](https://godoc.org/github.com/hoenirvili/rester?status.svg)](https://godoc.org/github.com/hoenirvili/rester) [![Go Report Card](https://goreportcard.com/badge/github.com/hoenirvili/rester)](https://goreportcard.com/report/github.com/hoenirvili/rester) [![Coverage Status](https://coveralls.io/repos/github/hoenirvili/rester/badge.svg?branch=master)](https://coveralls.io/github/hoenirvili/rester?branch=master) ![GitHub](https://img.shields.io/github/license/hoenirvili/rester.svg)

An Opinionated library for building REST APIs in Go.


## Simple rest resource with one method without checking the permissions

```go

package main

import (
	"net/http"

	"github.com/hoenirvili/rester"
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
	"github.com/hoenirvili/rester/response"
	"github.com/hoenirvili/rester/route"
)

type jsonResponse struct {
	Message string `json:"message"`
}

type root struct{}

func (r *root) index(req request.Request) resource.Response {
	return response.Payload(&jsonResponse{"Hello World !"})
}

func (r *root) Routes() route.Routes {
	return route.Routes{{
		URL:     "/",
		Method:  resource.Get,
		Handler: r.index,
	}}
}

func main() {
	rester := rester.New()
	rester.Resource("/", new(root))
	rester.Build()
	if err := http.ListenAndServe(":8080", rester); err != nil {
		panic(err)
	}
}

```


## Using the query api to query url parameters

```go

package main

import (
	"fmt"
	"net/http"

	"github.com/hoenirvili/rester"
	"github.com/hoenirvili/rester/query"
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
	"github.com/hoenirvili/rester/response"
	"github.com/hoenirvili/rester/route"
	"github.com/hoenirvili/rester/value"
)

type root struct {
}

type message struct {
	Str string `json:"message"`
}

func (r root) index(req request.Request) resource.Response {
	if err := req.Query("page").Error(); err != nil {
		return response.BadRequest("invalid page param")
	}

	page := req.Query("page").Int()
	str := fmt.Sprintf("get request with query param page=%d", page)
	return response.Payload(&message{str})
}

func (r root) Routes() route.Routes {
	return route.Routes{{
		URL:     "/",
		Method:  http.MethodGet,
		Handler: r.index,
		QueryPairs: query.Pairs{
			"page": query.Value{Type: value.Int},
		},
	}}
}
func main() {
	rester := rester.New()
	rester.Resource("/", new(root))
	rester.Build()
	if err := http.ListenAndServe(":8080", rester); err != nil {
		panic(err)
	}
}
```

### For more examples how to use the API, please visit the [examples](https://github.com/hoenirvili/rester/tree/master/examples) folder.
