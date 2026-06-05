package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikitagrgv/movies/internal/media"
	mediaStub "github.com/nikitagrgv/movies/internal/media/stub"
	mediaTmdb "github.com/nikitagrgv/movies/internal/media/tmdb"
	"github.com/nikitagrgv/movies/internal/server"
	"github.com/nikitagrgv/movies/internal/watch"
	web "github.com/nikitagrgv/movies/internal/web"

	"github.com/nikitagrgv/movies/internal/config"
	"github.com/nikitagrgv/movies/internal/media/tmdb"
	"github.com/nikitagrgv/movies/internal/watch/static"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	mux := http.NewServeMux()

	staticFs, err := fs.Sub(web.Assets, "static")
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

	tmpl, err := template.ParseFS(web.Assets, "templates/*.html", "templates/partials/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	noImageURL := "/static/noimage.png"
	tmdbApiURL := "https://api.themoviedb.org/3"
	tmdbImageURL := "https://image.tmdb.org/t/p"

	var mediaService *media.Service
	if cfg.IsStubUsed(config.MediaStub) {
		mediaService = media.NewService(
			mediaStub.NewMediaGetter(),
			mediaStub.NewMediaSearcher(),
			noImageURL,
		)
	} else {
		client, err := tmdb.NewClient(tmdbApiURL, tmdbImageURL, cfg.TmdbToken)
		if err != nil {
			log.Fatalf("Error loading tmdb client: %v", err)
		}

		mediaService = media.NewService(
			mediaTmdb.NewMediaGetter(client),
			mediaTmdb.NewMediaSearcher(client),
			noImageURL,
		)
	}

	var servers []static.WatchServerDescription
	for _, s := range cfg.WatchServers {
		servers = append(servers, static.WatchServerDescription{
			ID:                s.ID,
			Name:              s.Name,
			MovieURLTemplate:  s.MovieURLTemplate,
			TvShowURLTemplate: s.TvShowURLTemplate,
		})
	}

	watchProvider, err := static.NewWatchServerProvider(servers)
	if err != nil {
		log.Fatalf("Error loading watch servers: %v", err)
	}
	watchService := watch.NewService(watchProvider)

	handler := web.NewHandler(tmpl, mediaService, watchService)

	mux.Handle("GET /{$}", server.Chain(
		http.HandlerFunc(handler.ShowMain),
	))

	mux.Handle("GET /search", server.Chain(
		http.HandlerFunc(handler.HandleSearch),
	))

	mux.Handle("GET /movie/{id}", server.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			handler.HandleMovie(id, w, r)
		}),
	))

	mux.Handle("GET /tv/{id}", server.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			handler.HandleTvShow(id, w, r)
		}),
	))

	mux.Handle("/", server.Chain(
		http.HandlerFunc(handler.ShowNotFound),
	))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("Listening on port %d\n", cfg.ListenPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-stop
	fmt.Println("\nShutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
