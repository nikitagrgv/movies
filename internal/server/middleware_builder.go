package server

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// MiddlewareBuilder is an immutable builder for chaining HTTP middlewares.
type MiddlewareBuilder struct {
	middlewares []Middleware
}

func NewMiddlewareBuilder() MiddlewareBuilder {
	return MiddlewareBuilder{}
}

func (b MiddlewareBuilder) With(middleware Middleware) MiddlewareBuilder {
	oldLen := len(b.middlewares)
	cloned := make([]Middleware, oldLen, oldLen+1)
	copy(cloned, b.middlewares)
	cloned = append(cloned, middleware)
	return MiddlewareBuilder{middlewares: cloned}
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
