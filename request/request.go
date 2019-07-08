package request

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/hoenirvili/rester/query"
)

type Request struct {
	*http.Request
	pairs query.Pairs
}

func (r Request) Pairs() query.Pairs {
	return r.pairs
}

func New(r *http.Request, pairs query.Pairs) Request {
	if pairs == nil {
		pairs = make(query.Pairs)
	}
	return Request{r, pairs}
}

func (r *Request) Query(key string) *query.Value {
	value := r.pairs[key]
	if !value.Parsed() {
		r.pairs.Parse(key, r.URL.Query())
	}
	return value
}

func (r Request) URLParam(key string) string {
	return chi.URLParam(r.Request, key)
}
