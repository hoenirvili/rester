package route

import (
	"github.com/hoenirvili/rester/handler"
	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/query"
)

type Route struct {
	Allow   permission.Permissions
	Method  string
	URL     string
	Handler handler.Handler
	Queries query.Pairs
}

type Routes []Route
