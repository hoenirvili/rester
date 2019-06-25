package request

import (
	"net/http"

	"github.com/hoenirvili/rester/query"
)

type Request struct {
	*http.Request
	pairs map[string]query.Type
}

func New(r *http.Request, pairs map[string]query.Type) Request {
	return Request{r, pairs}
}

func (r Request) Query(key string) query.Value {
	t := r.pairs[key]
	return query.NewValue(t, r.URL.Query().Get(key))
}
