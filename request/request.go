package request

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/hoenirvili/rester/query"
)

type Request struct {
	*http.Request
	pairs map[string]query.Value
}

func (r Request) Pairs() map[string]query.Value {
	return r.pairs
}

func New(r *http.Request, pairs map[string]query.Value) Request {
	if pairs == nil {
		pairs = make(map[string]query.Value)
	}
	return Request{r, pairs}
}

func (r Request) Query(key string) query.Value {
	return r.pairs[key]
}

func (r Request) URLParam(key string) string {
	return chi.URLParam(r.Request, key)
}
