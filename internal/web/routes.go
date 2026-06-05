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

	mux.Handle("/static/", server.Chain(
		staticHandler,
		server.RecoveryMiddleware(),
		server.StripPrefixMiddleware("/static/"),
		server.GzipMiddleware,
	))
	mux.Handle("/favicon.ico", server.Chain(
		staticHandler,
		server.RecoveryMiddleware(),
		server.GzipMiddleware,
	))

	mux.Handle("GET /{$}", server.Chain(
		http.HandlerFunc(h.ShowMain),
		server.RecoveryMiddleware(),
	))

	mux.Handle("GET /search", server.Chain(
		http.HandlerFunc(h.HandleSearch),
		server.RecoveryMiddleware(),
	))

	mux.Handle("GET /movie/{id}", server.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			h.HandleMovie(id, w, r)
		}),
		server.RecoveryMiddleware(),
	))

	mux.Handle("GET /tv/{id}", server.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			h.HandleTvShow(id, w, r)
		}),
		server.RecoveryMiddleware(),
	))

	mux.Handle("/", server.Chain(
		http.HandlerFunc(h.ShowNotFound),
		server.RecoveryMiddleware(),
	))
}
