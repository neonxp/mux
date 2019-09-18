// Package mux implements simple regexp based http muxer
package mux

import (
	"context"
	"net/http"
	"regexp"
)

// Mux is a simple HTTP muxer. All path patterns must be valid go regexp (https://golang.org/pkg/regexp/syntax/).
type Mux struct {
	routes   map[string][]route
	NotFound http.Handler
}

// New returns new Mux instance
func New() *Mux {
	return &Mux{
		routes: map[string][]route{
			http.MethodGet:     {},
			http.MethodPost:    {},
			http.MethodPatch:   {},
			http.MethodPut:     {},
			http.MethodHead:    {},
			http.MethodDelete:  {},
			http.MethodConnect: {},
			http.MethodOptions: {},
			http.MethodTrace:   {},
		},
	}
}

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers := m.routes[r.Method]
	path := r.URL.EscapedPath()
	for _, route := range handlers {
		if route.pattern.MatchString(path) {
			matches := route.pattern.FindStringSubmatch(path)
			ctx := context.WithValue(r.Context(), "params", matches)
			route.handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}
	if m.NotFound != nil {
		m.NotFound.ServeHTTP(w, r)
	}
}

// Add registers handler for specified method and matched pattern
func (m *Mux) Add(method string, pattern string, handler http.Handler) {
	m.routes[method] = append(m.routes[method], route{
		pattern: regexp.MustCompile("(?i)" + pattern),
		handler: handler,
	})
}

// Get registers handler for GET method
func (m *Mux) Get(pattern string, handler http.Handler) {
	m.Add(http.MethodGet, pattern, handler)
}

// Post registers handler for POST method
func (m *Mux) Post(pattern string, handler http.Handler) {
	m.Add(http.MethodPost, pattern, handler)
}

// Put registers handler for PUT method
func (m *Mux) Put(pattern string, handler http.Handler) {
	m.Add(http.MethodPut, pattern, handler)
}

// Patch registers handler for PATCH method
func (m *Mux) Patch(pattern string, handler http.Handler) {
	m.Add(http.MethodPatch, pattern, handler)
}

// Options registers handler for OPTIONS method
func (m *Mux) Options(pattern string, handler http.Handler) {
	m.Add(http.MethodOptions, pattern, handler)
}

// Del registers handler for DELETE method
func (m *Mux) Del(pattern string, handler http.Handler) {
	m.Add(http.MethodDelete, pattern, handler)
}

// Head registers handler for HEAD method
func (m *Mux) Head(pattern string, handler http.Handler) {
	m.Add(http.MethodHead, pattern, handler)
}
