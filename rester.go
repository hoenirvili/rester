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
	notfound         http.HandlerFunc
	methodnotallowed http.HandlerFunc
	middleware       struct {
		global    []middleware
		validator middleware
	}
	resources map[string]Resource
	origins   []string
}

func (c *config) setValidator(m middleware) {
	c.middleware.validator = m
}

func (c *config) appendGlobal(m ...middleware) {
	c.middleware.global = append(c.middleware.global, m...)
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
			claims, err := r.options.validator.Verify(req)
			if err != nil {
				resp := response.Unauthorized(err.Error())
				resp.Render(w)
				return
			}

			p, ok := claims["permissions"]
			if !ok {
				resp := response.Unauthorized("No 'permissions' key found in the token")
				resp.Render(w)
				return
			}

			value, ok := p.(float64)
			if !ok {
				resp := response.Unauthorized("Invalid permission value")
				resp.Render(w)
				return
			}

			claims["permissions"] = permission.Permissions(value)
			ctx := req.Context()
			for key, value := range claims {
				ctx = context.WithValue(ctx, key, value)
			}
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}

	r.config.setValidator(middleware)
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

	// global middlwares holds a list of middlewares that will
	// be used in front of all routes
	globalMiddlewares middleware
}

// TokenValidator defines ways of interactions with the token
type TokenValidator interface {
	// Verify verifies if the request contains the desired token
	// This also can verify the expiration time or other claims
	// With success this will return the validated claims
	Verify(r *http.Request) (map[string]interface{}, error)
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
	// append into middleware stack
	r.config.notfound = httphandler(h, nil)
}

// UseGlobalMiddleware appends the list of middlewares into the global
// middleware stack to be put in front of every request emitted to the api
func (r *Rester) UseGlobalMiddleware(m ...middleware) {
	r.config.appendGlobal(m...)
}

// MethodNotAllowed defines a handler to respond whenever a method is
// not allowed on a route
func (r *Rester) MethodNotAllowed(h handler.Handler) {
	// append into middleware stack
	r.config.methodnotallowed = httphandler(h, nil)
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
		g.Route(r.options.version, func(router chi.Router) {
			router.NotFound(r.config.notfound)
			router.MethodNotAllowed(r.config.methodnotallowed)
			for _, middleware := range r.config.middleware.global {
				// global middlewares
				router.Use(middleware)
			}

			router.Use(r.corsMiddleware)

			for path, resource := range r.config.resources {
				r.resource(router, path, resource)
			}
		})
	})
}

func (r *Rester) corsMiddleware(next http.Handler) http.Handler {
	origins := strings.Join(r.config.origins, ",")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		switch origins {
		case "":
			w.Header().Set("Access-Control-Allow-Origin", "*")
		default:
			w.Header().Set("Access-Control-Allow-Origin", origins)
		}
		next.ServeHTTP(w, r)
	})
}

func allowAllRequests(permission.Permissions, request.Request) error { return nil }

func (r *Rester) decideWhichPermissionFunction(
	p permission.Permissions,
) func(permission.Permissions, request.Request) error {
	fn := allowAllRequests
	if p == permission.Anonymous {
		return fn
	}
	// if we did specify a token validation schema, proceed with checking
	// the permission return by the validation process in the context
	if r.options.validator != nil {
		fn = checkPermission
	}
	return fn
}

type makeHandlerConfig struct {
	isRequestAllowed func(permission.Permissions, request.Request) error
	route            route.Route
}

func makeHandler(c makeHandlerConfig) handler.Handler {
	return handler.Handler(func(req request.Request) resource.Response {
		if err := c.isRequestAllowed(c.route.Allow, req); err != nil {
			return response.Forbidden(err.Error())
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
		return c.route.Handler(req)
	})
}

func (r *Rester) resource(g chi.Router, base string, res Resource) {
	g.Route(base, func(router chi.Router) {
		routes := res.Routes()
		for _, route := range routes {
			r.validRoute(route)
			if route.Allow == 0 {
				route.Allow = permission.Anonymous
			}
			h := makeHandler(makeHandlerConfig{
				isRequestAllowed: r.decideWhichPermissionFunction(route.Allow),
				route:            route,
			})
			r.method(router, route, h)

			var handle http.Handler
			switch {
			case r.config.origins != nil:
				if route.Method == "OPTIONS" {
					continue
				}
				handle = enableCors(r.config.origins)
			case r.config.origins == nil:
				handle = http.HandlerFunc(enableAllCors)
			}
			router.Method(http.MethodOptions, route.URL, handle)
		}
	})
}

func enableAllCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func enableCors(origins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := r.Header.Get("Origin")
		if o == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, origin := range origins {
			if o == origin {
				header := w.Header()
				header.Set("Access-Control-Allow-Origin", o)
				header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
				header.Set("Access-Control-Allow-Headers", "content-type")
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		// TODO(hoenir): respond with 405?
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}

func (r *Rester) method(router chi.Router, route route.Route, h handler.Handler) {
	switch route.Allow {
	case permission.Anonymous:
		router.With(route.Middlewares...).MethodFunc(
			route.Method,
			route.URL,
			httphandler(h, route.QueryPairs),
		)
		return
	default:
		if r.options.validator != nil {
			router.With(r.config.middleware.validator).
				With(route.Middlewares...).
				MethodFunc(
					route.Method,
					route.URL,
					httphandler(h, route.QueryPairs),
				)
			return
		}
		router.With(route.Middlewares...).MethodFunc(
			route.Method,
			route.URL,
			httphandler(h, route.QueryPairs),
		)
	}
}

// Resource initializes a resource with the all available sub-routes of the resource
func (r *Rester) Resource(base string, resource Resource) {
	if _, ok := r.config.resources[base]; ok {
		panic("cannot append the same resource " + base + "twice")
	}
	r.config.resources[base] = resource
}

// Cors appends a list of origins that are allowed to make a request to the server
// By default this all preflight requests and other things are
// Access-Control-Allow-Origin *
func (r *Rester) Cors(origins ...string) {
	r.config.origins = origins
}

func httphandler(h handler.Handler, pairs query.Pairs) http.HandlerFunc {
	if h == nil {
		panic("no handler given for the route")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := request.New(r, pairs)
		response := h(req)
		response.Render(w)
	})
}
