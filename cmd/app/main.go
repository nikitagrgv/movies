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
	"github.com/nikitagrgv/movies/internal/infrastructure/movie_searcher/stub"
	"github.com/nikitagrgv/movies/internal/usecase"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	mux := http.NewServeMux()

	staticFs, err := fs.Sub(deliveryHttp.Assets, "static")
	if err != nil {
		log.Fatalf("Error loading static assets: %v", err)
	}

	staticHandler := http.FileServer(http.FS(staticFs))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	mux.Handle("/favicon.ico", staticHandler)

	tmpl, err := template.ParseFS(deliveryHttp.Assets, "templates/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	noImageURL := "/static/noimage.png"

	//searcher := tmdb.NewMovieSearcher()
	searcher := stub.NewMovieSearcher()
	search := usecase.NewSearchMoviesUsecase(searcher, noImageURL)

	//getter := tmdb.NewMovieGetter()
	getter := stub.NewMovieGetter()
	get := usecase.NewGetMovieUsecase(getter, noImageURL)

	handler := deliveryHttp.NewHandler(tmpl, search, get)

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

	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ShowNotFound(w, r)
	})

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.ListenPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h, pattern := mux.Handler(r)
			if pattern == "" {
				notFound.ServeHTTP(w, r)
				return
			}
			h.ServeHTTP(w, r)
		}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("Listening on port %d\n", cfg.ListenPort)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
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
