// Package rester used to define and compose http rest resources
// The package offers tiny handy abstractions to construct your rest api in a friendly and
// super mantainable way, yet to offer the flexibility required for a well more complex
// solution
package rester

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hoenirvili/rester/handler"
	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
	"github.com/hoenirvili/rester/response"
	"github.com/hoenirvili/rester/route"
)

// Rester used for constructing and initializing
// the http router with routes
type Rester struct {
	// r is the main underlying routes to dispatch and to assemble
	// all http rest resource handlers
	r chi.Router

	// o options type holding different kind of rest options
	o Options
}

// New returns a new Rester http.Handler comptabile that's ready
// to serve incoming  http rest request
func New(opts ...Option) *Rester {
	o := Options{}
	for _, setter := range opts {
		setter(&o)
	}
	r := &Rester{chi.NewRouter(), o}
	r.appendMiddlewares()
	return r
}

func (r *Rester) appendTokenMiddleware() {
	if r.o.validator == nil {
		return
	}

	r.r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if err := r.o.validator.Verify(req); err != nil {
				resp := response.Unauthorized(err.Error())
				resp.Render(w)
				return
			}

			permissions, err := r.o.validator.Extract()
			if err != nil {
				resp := response.Unauthorized(err.Error())
				resp.Render(w)
				return
			}

			ctx := context.WithValue(req.Context(), "permissions", permissions)
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	})
}

func (r *Rester) appendMiddlewares() {
	r.appendTokenMiddleware()
}

// guard return true if the permission on is set in
func guard(in permission.Permissions, on permission.Permissions) bool {
	return in&on != 0
}

// Option defines a setter callback type to set
// an underlying option value
type Option func(opt *Options)

// Options type holding all underlying options for the
// rester api
type Options struct {
	// validator used for token validation and extraction
	validator TokenValidator
}

// TokenValidator defines ways of interactinv with the token
type TokenValidator interface {
	// Verify verifies if the request contains the desired token
	// This also can verify the expiration time or other claims
	Verify(r *http.Request) error

	// Extract extracts the permission type from the token
	// to verify in the request chain what kind of callee we are dealing with
	Extract() (permission.Permissions, error)
}

// WithTokenValidator sets the underlying token validation implementation
// to use to validate and extract token metadata information to authorize and
// authenticate the
func WithTokenValidator(t TokenValidator) Option {
	return func(opts *Options) { opts.validator = t }
}

// NotFound defines a handler to respond whenever a route could
// not be found
func (r *Rester) NotFound(h handler.Handler) {
	r.r.NotFound(handler.HttpHandlerFunc(h))
}

// MethodNotAllowed defines a handler to respond whenever a method is
// not allowed
func (r *Rester) MethodNotAllowed(h handler.Handler) {
	r.r.MethodNotAllowed(handler.HttpHandlerFunc(h))
}

// ServeHTTP based on the incoming request route it to the available resource handler
func (r *Rester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

// Resource defines a way of composing routes into a resource
type Resource interface {
	// Routes returns a list of sub-routes of the given resource
	Routes() route.Routes
}

func checkPermission(allow permission.Permissions, req request.Request) error {
	in := req.Request.Context().Value("permission").(permission.Permissions)
	if !guard(in, allow) {
		return errors.New("You don't have permission to access this resource")
	}
	return nil
}

func (r *Rester) validRoute(route route.Route) {
	if route.Handler == nil {
		panic("Cannot use a nil handler")
	}

	if route.URL == "" {
		panic("Cannot use an empty URL route")
	}
}

// Resource initializes a resource with the all available sub-routes of the resource
func (r *Rester) Resource(base string, router Resource) {
	isRequestAllowed := func(allow permission.Permissions, req request.Request) error {
		return nil
	}

	if r.o.validator != nil {
		isRequestAllowed = checkPermission
	}

	for _, route := range router.Routes() {
		r.validRoute(route)

		if route.Allow == 0 {
			route.Allow = permission.Anonymous
		}

		handler := handler.HttpHandlerFunc(func(req request.Request) resource.Response {
			if err := isRequestAllowed(route.Allow, req); err != nil {
				return response.Unauthorized(err.Error())
			}
			return route.Handler(req)
		})

		r.r.Method(route.Method, route.URL, handler)
	}
}
