package route

import (
	"net/http"

	"github.com/hoenirvili/rester/handler"
	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/query"
)

// Route defines how a route should be treated by an http rest client
// This holds various ways of interaction with a rest resource
type Route struct {
	// Allow defines a permission bit flag that can be set
	// To specify what users can access this resource
	Allow permission.Permissions
	// Method is the main http method that the route will respond to
	Method string
	// URL holds the relative URL of the resource
	URL string
	// Handler main handler that will be called in a separate go routine
	// by the main router to handle the client's request
	Handler handler.Handler
	// QueryPairs holds a list of query parameters key and value used for
	// retrieving different types of values
	QueryPairs query.Pairs
	// Middlewares list of middlewares that will be executed first one by one
	// like a chain before executing the main Handler
	Middlewares []func(http.Handler) http.Handler
}

// Routes defines a slice of routes for a resource
type Routes []Route
