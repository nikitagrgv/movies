package web

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/nikitagrgv/movies/internal/httpsrv"
	"github.com/nikitagrgv/movies/internal/logger"
)

func ResolveStaticAssetPath(cacheVersion int, relPath string) string {
	versionedPath := fmt.Sprintf("/static/v%d/%s", cacheVersion, relPath)
	return versionedPath
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, cacheVersion int, logger *logger.Service) {
	const cacheControlTime = time.Hour * 24 * 365

	staticFs, err := fs.Sub(Assets, "static")
	if err != nil {
		log.Fatalf("Error loading static assets: %v", err)
	}

	staticHandler := http.FileServer(http.FS(staticFs))

	baseMiddleware := httpsrv.NewMiddlewareBuilder().
		With(httpsrv.RecoveryMiddleware).
		With(LoggerMiddleware(logger))

	staticMiddleware := baseMiddleware.
		With(httpsrv.CacheControlMiddleware(cacheControlTime)).
		With(httpsrv.GzipMiddleware)

	versionedPath := ResolveStaticAssetPath(cacheVersion, "")
	mux.Handle("GET "+versionedPath, staticMiddleware.
		With(httpsrv.StripPrefixMiddleware(versionedPath)).
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
