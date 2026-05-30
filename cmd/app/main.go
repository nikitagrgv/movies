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

	"github.com/nikitagrgv/movies/internal/config"
	deliveryHttp "github.com/nikitagrgv/movies/internal/delivery/http"
	"github.com/nikitagrgv/movies/internal/domain"
	"github.com/nikitagrgv/movies/internal/infrastructure/movie/stub"
	"github.com/nikitagrgv/movies/internal/infrastructure/movie/tmdb"
	"github.com/nikitagrgv/movies/internal/infrastructure/watch/static"
	"github.com/nikitagrgv/movies/internal/usecase"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	mux := http.NewServeMux()

	staticFs, err := fs.Sub(deliveryHttp.Assets, "static")
	if err != nil {
		log.Fatalf("Error loading static assets: %v", err)
	}

	staticHandler := http.FileServer(http.FS(staticFs))
	mux.Handle("/static/",
		http.StripPrefix("/static/",
			deliveryHttp.GzipMiddleware(staticHandler)))
	mux.Handle("/favicon.ico",
		deliveryHttp.GzipMiddleware(staticHandler),
	)

	tmpl, err := template.ParseFS(deliveryHttp.Assets, "templates/*.html", "templates/partials/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	noImageURL := "/static/noimage.png"
	tmdbApiURL := "https://api.themoviedb.org/3"
	tmdbImageURL := "https://image.tmdb.org/t/p"
	tmdbClient, err := tmdb.NewClient(tmdbApiURL, tmdbImageURL, cfg.TmdbToken)
	if err != nil {
		log.Fatalf("Error loading aggregator client: %v", err)
	}

	var searcher domain.MediaSearcher
	if cfg.IsStubUsed(config.SearchStub) {
		searcher = stub.NewMediaSearcher()
	} else {
		searcher = tmdb.NewMediaSearcher(tmdbClient)
	}
	search := usecase.NewSearchMediaUsecase(searcher, noImageURL)

	var getter domain.MediaGetter
	if cfg.IsStubUsed(config.SearchStub) {
		getter = stub.NewMediaGetter()
	} else {
		getter = tmdb.NewMediaGetter(tmdbClient)
	}
	get := usecase.NewGetMediaUsecase(getter, noImageURL)

	var servers []static.WatchServer
	for _, s := range cfg.WatchServers {
		servers = append(servers, static.WatchServer{
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
	watch := usecase.NewWatchServerUsecase(watchProvider)

	handler := deliveryHttp.NewHandler(tmpl, search, get, watch)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		handler.ShowMain(w, r)
	})

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleSearch(w, r)
	})

	mux.HandleFunc("GET /movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handler.HandleMovie(id, w, r)
	})

	mux.HandleFunc("GET /tv/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handler.HandleTvShow(id, w, r)
	})

	mux.HandleFunc("/", handler.ShowNotFound)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		Handler:      deliveryHttp.GzipMiddleware(mux),
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
