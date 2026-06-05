package server

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// MiddlewareBuilder
// Immutable builder
type MiddlewareBuilder struct {
	middlewares []Middleware
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (b MiddlewareBuilder) With(middleware Middleware) MiddlewareBuilder {
	b.middlewares = append(b.middlewares, middleware)
	return b
}

func (b MiddlewareBuilder) Build(handler http.Handler) http.Handler {
	return chain(handler, b.middlewares)
}

func chain(handler http.Handler, middlewares []Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
