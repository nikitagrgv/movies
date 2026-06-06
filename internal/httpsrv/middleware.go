package httpsrv

import (
	"log"
	"net/http"
	"runtime/debug"
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

// Generic Middlewares

func StripPrefixMiddleware(prefix string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	}
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC RECOVERED] %v\nStack Trace:\n%s", err, debug.Stack())
				w.Header().Set("Connection", "close")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
