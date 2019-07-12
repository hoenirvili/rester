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
