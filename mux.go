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
	path := strings.Trim(r.URL.EscapedPath(), " /")
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
	pattern = strings.ToLower(strings.Trim(pattern, " /"))
	p := make([]string, 0)
	inparam := false
	part := ""
	for _, e := range pattern {
		if inparam {
			if !isAlpha(e) && !isNum(e) {
				inparam = false
				if part != "" {
					p = append(p, part)
					part = ""
				}
			}
		} else {
			if e == ':' {
				inparam = true
				if part != "" {
					p = append(p, part)
					part = ""
				}
			}
		}
		part += string(e)
	}
	if part != "" {
		p = append(p, part)
	}
	m.routes[method] = append(m.routes[method], route{
		pattern: p,
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

func match(pattern []string, test string) (map[string]string, bool) {
	matches := make(map[string]string, 0)
	matched := ""
	for idx, p := range pattern {
		if p[:1] == ":" {
			if idx == len(pattern)-1 {
				matches[p[1:]] = test
				return matches, true
			}
			matched, test = splitStr(test, pattern[idx+1])
			matches[p[1:]] = matched
		} else {
			if len(p) <= len(test) && p == test[:len(p)] {
				test = test[len(p):]
				continue
			}
			return nil, false
		}
	}
	return matches, true
}

func splitStr(str string, m string) (string, string) {
	for i := 0; i <= len(str)-len(m); i++ {
		if str[i:(i+len(m))] == m {
			return str[:i], str[i:]
		}
	}
	return str, ""
}

func isAlpha(e int32) bool {
	return e >= 'a' && e <= 'z'
}

func isNum(e int32) bool {
	return e >= '0' && e <= '9'
}
