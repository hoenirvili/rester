package request

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/hoenirvili/rester/permission"
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

func (r *Request) Permission() permission.Permissions {
	return r.Context().Value("permissions").(permission.Permissions)
}

func (r *Request) Query(key string) value.Value {
	input := ""
	v, ok := r.URL.Query()[key]
	if ok {
		input = v[0]
	}
	return value.Parse(input, r.pairs[key].Type)
}

func (r Request) URLParam(key string, t value.Type) value.Value {
	input := chi.URLParam(r.Request, key)
	return value.Parse(input, t)
}
