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
		server.StripPrefix("/static/"),
		server.GzipMiddleware,
	))
	mux.Handle("/favicon.ico", server.Chain(
		staticHandler,
		server.GzipMiddleware,
	))

	mux.Handle("GET /{$}", server.Chain(
		http.HandlerFunc(h.ShowMain),
	))

	mux.Handle("GET /search", server.Chain(
		http.HandlerFunc(h.HandleSearch),
	))

	mux.Handle("GET /movie/{id}", server.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			h.HandleMovie(id, w, r)
		}),
	))

	mux.Handle("GET /tv/{id}", server.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			h.HandleTvShow(id, w, r)
		}),
	))

	mux.Handle("/", server.Chain(
		http.HandlerFunc(h.ShowNotFound),
	))
}
