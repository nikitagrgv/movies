package main

import (
	"context"
	"errors"
	"fmt"
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
	mediaCache "github.com/nikitagrgv/movies/internal/media/cached"
	mediaStub "github.com/nikitagrgv/movies/internal/media/stub"
	mediaTmdb "github.com/nikitagrgv/movies/internal/media/tmdb"
	"github.com/nikitagrgv/movies/internal/pkg/cache"
	"github.com/nikitagrgv/movies/internal/pkg/postgres"
	"github.com/nikitagrgv/movies/internal/watch"
	"github.com/nikitagrgv/movies/internal/web"
	"golang.org/x/sync/errgroup"

	"github.com/nikitagrgv/movies/internal/config"
	"github.com/nikitagrgv/movies/internal/media/tmdb"
	"github.com/nikitagrgv/movies/internal/watch/static"
)

const (
	cacheVersion    = 2 // Increment when static web files changes to invalidate browser caches
	staticFilesHash = "a696737ccbd1ee5325c118b119e13a416b390f1e0bb2b34ad0822da271b9c66f"

	tmdbApiURL   = "https://api.themoviedb.org/3"
	tmdbImageURL = "https://image.tmdb.org/t/p"
)

func main() {
	staticHash, err := web.GetStaticFilesHash()
	if err != nil {
		log.Fatalf("Failed to get static files hash: %v", err)
	}

	if staticHash != staticFilesHash {
		log.Fatalf("Static files hash does not match. Current hash: %s", staticHash)
	}

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

	redisClient, err := cache.NewRedisClient(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
	}
	defer redisClient.Close()

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

		var getter media.Getter = mediaTmdb.NewMediaGetter(client)
		getter = mediaCache.NewMediaGetter(redisClient, getter)

		var searcher media.Searcher = mediaTmdb.NewMediaSearcher(client)

		mediaService = media.NewService(
			getter,
			searcher,
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

	tmpl, err := web.LoadTemplates(cacheVersion)
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
