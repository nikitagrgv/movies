package web

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/nikitagrgv/movies/internal/server"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	staticFs, err := fs.Sub(Assets, "static")
	if err != nil {
		log.Fatalf("Error loading static assets: %v", err)
	}

	staticHandler := http.FileServer(http.FS(staticFs))

	baseMiddleware := server.NewMiddlewareBuilder().
		With(server.RecoveryMiddleware())

	mux.Handle("/static/", baseMiddleware.
		With(server.StripPrefixMiddleware("/static/")).
		With(server.GzipMiddleware).
		Build(staticHandler))

	mux.Handle("/favicon.ico", baseMiddleware.
		With(server.GzipMiddleware).
		Build(staticHandler))

	mux.Handle("GET /{$}", baseMiddleware.
		Build(http.HandlerFunc(h.showMain)))

	mux.Handle("GET /search", baseMiddleware.
		Build(http.HandlerFunc(h.handleSearch)))

	mux.Handle("GET /movie/{id}", baseMiddleware.
		Build(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			h.handleMovie(id, w, r)
		})))

	mux.Handle("GET /tv/{id}", baseMiddleware.
		Build(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			h.handleTvShow(id, w, r)
		})))

	mux.Handle("/", baseMiddleware.
		Build(http.HandlerFunc(h.showNotFound)))
}
