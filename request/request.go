package request

import (
	"net/http"

	"github.com/hoenirvili/rester/query"
)

type Request struct {
	*http.Request
	query.Query
}

func New(r *http.Request, q query.Query) Request {
	return Request{r, q}
}
