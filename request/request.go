package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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

func (r Request) Permission() permission.Permissions {
	value := r.Context().Value("permissions")
	if value == nil {
		return permission.NoPermission
	}
	return value.(permission.Permissions)
}

func (r Request) Query(key string) value.Value {
	input := ""
	v, ok := r.URL.Query()[key]
	if ok {
		input = v[0]
	}
	return value.Parse(input, r.pairs[key].Type)
}

func (r Request) JSON(p interface{}) error {
	defer r.Request.Body.Close()
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return errors.New("Invalid content type header, we only support application/json")
	}
	return json.NewDecoder(r.Request.Body).Decode(p)
}

func (r Request) URLParam(key string, t value.Type) value.Value {
	input := chi.URLParam(r.Request, key)
	return value.Parse(input, t)
}

func (r Request) ID() (uint64, error) {
	value := r.URLParam("id", value.Uint64)
	if err := value.Error(); err != nil {
		return 0, err
	}
	return value.Uint64(), nil
}
