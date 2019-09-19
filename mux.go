// Package mux implements simple HTTP with greedy patterns
package mux

import (
	"context"
	"net/http"
	"strings"
)

// Mux is a simple HTTP muxer
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
	pattern []string
	handler http.Handler
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers := m.routes[r.Method]
	path := strings.Split(strings.Trim(r.URL.EscapedPath(), " /"), "/")
	for _, route := range handlers {
		if matches, ok := match(route.pattern, path); ok {
			ctx := context.WithValue(r.Context(), "params", matches)
			route.handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}
	if m.NotFound != nil {
		m.NotFound.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not found."))
	}
}

// Add registers handler for specified method and matched pattern
func (m *Mux) Add(method string, pattern string, handler http.Handler) {
	m.routes[method] = append(m.routes[method], route{
		pattern: strings.Split(strings.Trim(pattern, " /"), "/"),
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

//GetParams extracts route parameters from request
func GetParams(r *http.Request) map[string]string {
	params, ok := r.Context().Value("params").(map[string]string)
	if !ok {
		return nil
	}
	return params
}

func match(pattern []string, test []string) (map[string]string, bool) {
	if len(test) < len(pattern) {
		return nil, false
	}
	matches := make(map[string]string, 0)
	acc := make([]string, 0)
	key := ""
	cp := ""
	ct := ""
	for {
		cp, pattern = pop(pattern)
		if cp == "" && len(pattern) == 0 {
			return matches, true
		}
		// current pattern - parameter key
		if cp[:1] == ":" {
			key = cp[1:]
			acc = make([]string, 0)
			// it is last parameter - so we just put rest test substrings to it value
			if len(pattern) == 0 {
				matches[key] = strings.Join(test, "/")
				return matches, true
			}
			continue
		}
		if key == "" {
			// we not in parameter - just check that next test substring equals current pattern
			ct, test = pop(test)
			if ct == "" && len(test) == 0 {
				return nil, false
			}
			if ct == cp {
				continue
			}
			return nil, false
		}
		// we in parameter - pushing test substrings to it value, while test substring not equals current pattern
		for {
			ct, test = pop(test)
			if ct == "" && len(test) == 0 {
				return nil, false
			}
			if ct == cp {
				matches[key] = strings.Join(acc, "/")
				key = ""
				break
			}
			acc = append(acc, ct)
		}
	}
}

func pop(arr []string) (string, []string) {
	if len(arr) == 0 {
		return "", []string{}
	}
	head, tail := arr[0], arr[1:]
	if head == "" && len(tail) > 0 {
		return pop(tail)
	}
	return strings.ToLower(head), tail
}
