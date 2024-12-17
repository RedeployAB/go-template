package server

import (
	"net/http"
	"strings"
)

// NewRouter creates a new router.
func NewRouter() *router {
	return &router{
		ServeMux: http.NewServeMux(),
	}
}

// router is a custom router that implements the http.Handler interface.
type router struct {
	*http.ServeMux
}

// ServeHTTP wraps the http.ServeMux ServeHTTP method.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.ServeMux.ServeHTTP(w, req)
}

// Handle wraps the http.ServeMux Handle method that adds handlers for patterns with
// and without trailing slashes. This to support adding sub routers (handlers) once for
// both patterns.
func (r *router) Handle(pattern string, handler http.Handler) {
	r.ServeMux.Handle(pattern, handler)
	if pattern == "/{$}" {
		return
	}
	if strings.HasSuffix(pattern, "/") {
		trimmedPattern := strings.TrimSuffix(pattern, "/")
		if len(trimmedPattern) > 0 {
			r.ServeMux.Handle(trimmedPattern, handler)
		}
	} else {
		r.ServeMux.Handle(pattern+"/", handler)
	}
}
