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
	Message string `json:"message,omitempty"`
}

type root struct{}

func (r *root) index(req request.Request) resource.Response {
	payload := &jsonResponse{"Hello world!"}
	return response.Payload(payload)
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
