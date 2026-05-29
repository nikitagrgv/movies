package domain

import "context"

type WatchServiceProvider interface {
	GetServices(ctx context.Context) ([]WatchService, error)

	GetMovieWatchLink(ctx context.Context, serviceID, movieID int) (string, error)

	GetTvShowWatchLink(ctx context.Context, serviceID, tvID, season, episode int) (string, error)
}
