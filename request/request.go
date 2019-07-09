package request

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/hoenirvili/rester/query"
	"github.com/hoenirvili/rester/value"
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

func (r *Request) Query(key string) value.Value {
	input := r.URLParam(key)
	return value.Parse(input, r.pairs[key].Type)
}

func (r Request) URLParam(key string) string {
	return chi.URLParam(r.Request, key)
}
