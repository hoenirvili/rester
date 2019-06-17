package resource

import (
	"net/http"
)

type Response interface {
	Render(w http.ResponseWriter)
}

const (
	Get     = "GET"
	Head    = "HEAD"
	Post    = "POST"
	Put     = "PUT"
	Patch   = "PATCH"
	Delete  = "DELETE"
	Connect = "CONNECT"
	Options = "OPTIONS"
	Trace   = "TRACE"
)
