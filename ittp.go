package ittp

import (
	"fmt"
	"net/http"
)

type (
	// ServeMux is a wrapper around the standard library ServeMux and is intended to work
	// in the exact same way except with middlewareFuncs and helper functions.
	ServeMux struct {
		stdMux      *http.ServeMux
		middlewares []middlewareFunc
	}
	middlewareFunc func(next http.Handler) http.Handler
)

// NewServeMux allocates and returns a new ittp.ServeMux
func NewServeMux() *ServeMux {
	return &ServeMux{
		http.NewServeMux(),
		[]middlewareFunc{},
	}
}

// ServeHTTP adds the middlewareFuncs to http.ServeMux in chronological order, then calls ServeHTTP.
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = m.stdMux
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		h = m.middlewares[i](h)
	}
	h.ServeHTTP(w, r)
}

// MethodHandleFunc uses go 1.22 routing with the method specified.
func (m *ServeMux) MethodHandleFunc(method string, path string, handler http.HandlerFunc) {
	m.stdMux.HandleFunc(fmt.Sprintf("%s %s", method, path), handler)
}

// MethodHandle uses go 1.22 routing with the method specified.
func (m *ServeMux) MethodHandle(method string, path string, handler http.Handler) {
	m.stdMux.Handle(fmt.Sprintf("%s %s", method, path), handler)
}

// HandleFunc just passes pattern and handler to the standard mux.
// Allows ittp to be used as a direct replacement to net/http package.
func (m *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	m.stdMux.HandleFunc(pattern, handler)
}

// Handle just passes pattern and handler to the standard mux.
// Allows ittp to be used as a direct replacement to net/http package.
func (m *ServeMux) Handle(pattern string, handler http.Handler) {
	m.stdMux.Handle(pattern, handler)
}

func (m *ServeMux) Get(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodGet, path, handler)
}

func (m *ServeMux) Head(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodHead, path, handler)
}

func (m *ServeMux) Post(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodPost, path, handler)
}

func (m *ServeMux) Put(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodPut, path, handler)
}

func (m *ServeMux) Patch(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodPatch, path, handler)
}

func (m *ServeMux) Delete(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodDelete, path, handler)
}

func (m *ServeMux) Connect(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodConnect, path, handler)
}

func (m *ServeMux) Options(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodOptions, path, handler)
}

func (m *ServeMux) Trace(path string, handler http.HandlerFunc) {
	m.MethodHandleFunc(http.MethodTrace, path, handler)
}

// AddMiddleware adds a ittp.middlewareFunc to ServeMux's middlewareFunc slice.
// middlewareFuncs are executed in chronological order.
func (m *ServeMux) AddMiddleware(middlware middlewareFunc) {
	m.middlewares = append(m.middlewares, middlware)
}

func (m *ServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	return m.stdMux.Handler(r)
}
