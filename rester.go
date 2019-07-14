// Package rester used to define and compose multiple HTTP REST resources
// The package offers abstractions in order to construct your REST API
// in a friendly and maintainable way, yet offering the flexibility
// for a well more complex solution
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

type config struct {
	notfound        http.HandlerFunc
	methodnotallowd http.HandlerFunc
	middlewares     []middleware
	resources       map[string]Resource
}

// Rester used for constructing and initializing the http router with routes
type Rester struct {
	root    chi.Router
	options Options
	config  config
}

type middleware func(http.Handler) http.Handler

// New returns a new Rester http.Handler compatible that's ready
// to serve incoming  http rest request
func New(opts ...Option) *Rester {
	options := Options{}
	for _, setter := range opts {
		setter(&options)
	}
	r := &Rester{
		root:    chi.NewRouter(),
		options: options,
		config: config{
			resources: make(map[string]Resource),
		},
	}
	r.appendTokenMiddleware()
	return r
}

func (r *Rester) appendTokenMiddleware() {
	if r.options.validator == nil {
		return
	}

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if err := r.options.validator.Verify(req); err != nil {
				resp := response.Unauthorized(err.Error())
				resp.Render(w)
				return
			}

			permissions, err := r.options.validator.Extract()
			if err != nil {
				resp := response.Unauthorized(err.Error())
				resp.Render(w)
				return
			}

			ctx := context.WithValue(req.Context(), "permissions", permissions)
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}

	r.config.middlewares = append(r.config.middlewares, middleware)
}

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
	// version adds the the api version as the base route
	version string
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

// WithVersioning appends to the path route the prefix "/version/"
func WithVersioning(version string) Option {
	return func(opts *Options) { opts.version = "/" + version }
}

// NotFound defines a handler to respond whenever a route could not be found
func (r *Rester) NotFound(h handler.Handler) {
	//append into middleware stack
	r.config.notfound = httphandler(h, nil)
}

// MethodNotAllowed defines a handler to respond whenever a method is
// not allowed on a route
func (r *Rester) MethodNotAllowed(h handler.Handler) {
	// append into middleware stack
	r.config.methodnotallowd = httphandler(h, nil)
}

// ServeHTTP based on the incoming request route it to the available resource handler
func (r *Rester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.root.ServeHTTP(w, req)
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

func serveFiles(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("serving files does not permit URL parameters")
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

// Static serve static files from the location pointed by dir
// This has limited support so don't expect much customization
func (r *Rester) Static(dir string) {
	r.root.Group(func(g chi.Router) {
		serveFiles(g, "/", http.Dir(dir))
	})
}

// Build builds the internal state of the router making it ready for
// dispatching requests
func (r *Rester) Build() {
	r.root.Group(func(g chi.Router) {
		if r.options.version == "" {
			r.options.version = "/"
		}
		router := chi.NewRouter()
		g.Mount(r.options.version, router)
		router.NotFound(r.config.notfound)
		router.MethodNotAllowed(r.config.methodnotallowd)
		for _, middleware := range r.config.middlewares {
			// global middlewares
			router.Use(middleware)
		}
		for path, resource := range r.config.resources {
			r.resource(router, path, resource)
		}
	})
}

func (r *Rester) resource(g chi.Router, path string, res Resource) {
	g.Route(path, func(router chi.Router) {
		isRequestAllowed := func(permission.Permissions, request.Request) error {
			return nil
		}
		// if we did specify a token validation schema, proceed with checking
		// the permission return by the validation process in the context
		if r.options.validator != nil {
			isRequestAllowed = checkPermission
		}
		for _, route := range res.Routes() {
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
			router.Method(route.Method, route.URL, httphandler(h, route.QueryPairs))
		}
	})
}

// Resource initializes a resource with the all available sub-routes of the resource
func (r *Rester) Resource(base string, resource Resource) {
	if _, ok := r.config.resources[base]; ok {
		panic("cannot append the same resource " + base + "twice")
	}
	r.config.resources[base] = resource
}

func httphandler(h handler.Handler, pairs query.Pairs) http.HandlerFunc {
	if h == nil {
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := request.New(r, pairs)
		response := h(request)
		response.Render(w)
	})
}
