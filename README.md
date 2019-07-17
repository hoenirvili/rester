# rester
[![Build Status](https://travis-ci.com/hoenirvili/rester.svg?branch=master)](https://travis-ci.com/hoenirvili/rester) [![GoDoc](https://godoc.org/github.com/hoenirvili/rester?status.svg)](https://godoc.org/github.com/hoenirvili/rester) [![Go Report Card](https://goreportcard.com/badge/github.com/hoenirvili/rester)](https://goreportcard.com/report/github.com/hoenirvili/rester) [![Coverage Status](https://coveralls.io/repos/github/hoenirvili/rester/badge.svg?branch=master)](https://coveralls.io/github/hoenirvili/rester?branch=master)

An Opinionated library for building REST APIs in Go.


# Simple rest resource with one method without checking the permissions

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
