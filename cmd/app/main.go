package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	const noImageURL = "/static/noimage.png"
	var mediaService *media.Service
	if cfg.IsStubUsed(config.MediaStub) {
		mediaService = media.NewService(
			mediaStub.NewMediaGetter(),
			mediaStub.NewMediaSearcher(),
			noImageURL,
		)
	} else {
		const tmdbApiURL = "https://api.themoviedb.org/3"
		const tmdbImageURL = "https://image.tmdb.org/t/p"

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

	tmpl, err := template.ParseFS(web.Assets, "templates/*.html", "templates/partials/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	mux := http.NewServeMux()
	handler := web.NewHandler(tmpl, mediaService, watchService)
	handler.RegisterRoutes(mux)

	srv := server.NewServer(cfg.ListenPort, mux)
	srv.Run(ctx)

	log.Println("Server cleanly stopped")
}
