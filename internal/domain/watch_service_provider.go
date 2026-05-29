package domain

import "context"

type WatchServerProvider interface {
	GetServers(ctx context.Context) ([]WatchServer, error)

	GetMovieWatchLink(ctx context.Context, serverID, movieID int) (string, error)

	GetTvShowWatchLink(ctx context.Context, serverID, tvID, season, episode int) (string, error)
}
