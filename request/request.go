package request

import "net/http"

type Request struct {
	*http.Request
}

func New(r *http.Request) Request {
	return Request{r}
}
