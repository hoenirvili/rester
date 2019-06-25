package request

import (
	"net/http"

	"github.com/hoenirvili/rester/query"
)

type Request struct {
	*http.Request
	pairs map[string]query.Value
}

func (r Request) Pairs() map[string]query.Value {
	return r.pairs
}

func New(r *http.Request) Request {
	return Request{Request: r}
}

func (r Request) Query(key string) query.Value {
	return r.pairs[key]
}
