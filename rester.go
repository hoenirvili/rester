// Package rester used to define and compose http rest resources
// The package offers tiny handy abstractions to construct your rest api in a friendly and
// super mantainable way, yet to offer the flexibility required for a well more complex
// solution
package rester

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/hoenirvili/rester/handler"
	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/query"
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

// New returns a new Rester http.Handler compatible that's ready
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
	// if we do not have a custom validator set
	// do not include the token middleware
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

// TokenValidator defines ways of interactions with the token
type TokenValidator interface {
	// Verify verifies if the request contains the desired token
	// This also can verify the expiration time or other claims
	Verify(r *http.Request) error

	// Extract extracts the permission type from the token
	// to verify in the request chain what kind of callee we are dealing with
	Extract() (permission.Permissions, error)
}

// WithTokenValidator sets the underlying token validation implementation
// to use to validate and extract token meta-data information to authorize and
// authenticate the
func WithTokenValidator(t TokenValidator) Option {
	return func(opts *Options) { opts.validator = t }
}

// NotFound defines a handler to respond whenever a route could
// not be found
func (r *Rester) NotFound(h handler.Handler) {
	r.r.NotFound(httphandler(h, nil))
}

// MethodNotAllowed defines a handler to respond whenever a method is
// not allowed
func (r *Rester) MethodNotAllowed(h handler.Handler) {
	r.r.MethodNotAllowed(httphandler(h, nil))
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

// checkPermission checks if the value with the key permission exists and if
// it passes the guard check
func checkPermission(allow permission.Permissions, req request.Request) error {
	in := req.Request.Context().Value("permissions").(permission.Permissions)
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

type ResourceInliner interface {
	ResourceInline() []Resource
}

func serveFiles(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		},
	))
}

func (r *Rester) Static(path string, dir http.Dir) {
	r.r.Group(func(groupRouter chi.Router) {
		serveFiles(r.r, path, dir)
	})
}

func (r *Rester) ResourceInline(base string, i ResourceInliner) {
	r.r.Group(func(groupRouter chi.Router) {
		for _, resource := range i.ResourceInline() {
			r.resource(groupRouter, base, resource)
		}
	})
}

func (r *Rester) resource(groupRouter chi.Router, base string, router Resource) {
	isRequestAllowed := func(permission.Permissions, request.Request) error {
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

		h := func(req request.Request) resource.Response {
			if err := isRequestAllowed(route.Allow, req); err != nil {
				return response.Unauthorized(err.Error())
			}

			values := req.URL.Query()
			pairs := req.Pairs()
			for key := range pairs {
				if pairs[key].Required {
					if err := pairs.Parse(key, values); err != nil {
						return response.BadRequest(err.Error())
					}
				}
			}
			return route.Handler(req)
		}
		groupRouter.Method(route.Method, route.URL, httphandler(h, route.QueryPairs))
	}
}

// Resource initializes a resource with the all available sub-routes of the resource
func (r *Rester) Resource(base string, router Resource) {
	r.resource(r.r, base, router)
}

func httphandler(h handler.Handler, pairs query.Pairs) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := request.New(r, pairs)
		response := h(request)
		response.Render(w)
	})
}
