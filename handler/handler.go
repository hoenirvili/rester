package handler

import (
	"net/http"

	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
)

type Handler func(request.Request) resource.Response

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := request.New(r)
	response := h(req)
	response.Render(w)
}

func HandlerFunc(handler Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}
