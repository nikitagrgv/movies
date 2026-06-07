package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikitagrgv/movies/internal/grpc"
	"github.com/nikitagrgv/movies/internal/httpsrv"
	"github.com/nikitagrgv/movies/internal/logger"
	postgresLogRepo "github.com/nikitagrgv/movies/internal/logger/postgres"
	"github.com/nikitagrgv/movies/internal/media"
	mediaStub "github.com/nikitagrgv/movies/internal/media/stub"
	mediaTmdb "github.com/nikitagrgv/movies/internal/media/tmdb"
	"github.com/nikitagrgv/movies/internal/pkg/postgres"
	"github.com/nikitagrgv/movies/internal/watch"
	"github.com/nikitagrgv/movies/internal/web"
	"golang.org/x/sync/errgroup"

	"github.com/nikitagrgv/movies/internal/config"
	"github.com/nikitagrgv/movies/internal/media/tmdb"
	"github.com/nikitagrgv/movies/internal/watch/static"
)

func makePostgresConfig(cfg config.DbConfig, schema string) postgres.Config {
	return postgres.NewConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Db,
		schema,
	))
}

const (
	cacheVersion = 1 // Increment when static web files changes to invalidate browser caches
	tmdbApiURL   = "https://api.themoviedb.org/3"
	tmdbImageURL = "https://image.tmdb.org/t/p"
)

func main() {
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(signalCtx)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	loggerDbPool, err := postgres.Connect(gCtx, makePostgresConfig(cfg.Db, "logger"))
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	defer loggerDbPool.Close()

	visitRepo := postgresLogRepo.NewVisitRepository(loggerDbPool)
	loggerService := logger.NewService(visitRepo)

	noImageURL := web.ResolveStaticAssetPath(cacheVersion, "noimage.png")
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

	tmpl, err := template.ParseFS(web.Assets, "templates/*.html", "templates/partials/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	mux := http.NewServeMux()
	handler := web.NewHandler(tmpl, mediaService, watchService)
	handler.RegisterRoutes(mux, cacheVersion, loggerService)

	httpServer := httpsrv.NewServer(cfg.ListenPort, mux)
	g.Go(func() error {
		return httpServer.Run(gCtx)
	})

	grpcServer := grpc.NewServer(cfg.GRPCListenPort)
	g.Go(func() error {
		return grpcServer.Run(gCtx)
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("FATAL: %v", err)
	}

	log.Println("Server stopped")
}
