package route

import (
	"net/http"

	"github.com/hoenirvili/rester/handler"
	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/query"
)

type Route struct {
	Allow       permission.Permissions
	Method      string
	URL         string
	Handler     handler.Handler
	QueryPairs  query.Pairs
	Middlewares []func(http.Handler) http.Handler
}

type Routes []Route
