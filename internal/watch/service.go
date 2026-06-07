package watch

import "context"

type provider interface {
	GetServers(ctx context.Context) ([]WatchServer, error)
	GetMovieWatchLink(ctx context.Context, serverID string, movieID int) (string, error)
	GetTvShowWatchLink(ctx context.Context, serverID string, tvID, season, episode int) (string, error)
}

type Service struct {
	provider provider
}

func NewService(provider provider) *Service {
	return &Service{provider: provider}
}

func (u *Service) GetServers(ctx context.Context) ([]WatchServer, error) {
	return u.provider.GetServers(ctx)
}

func (u *Service) GetMovieWatchLink(ctx context.Context, serverID string, movieID int) (string, error) {
	return u.provider.GetMovieWatchLink(ctx, serverID, movieID)
}

func (u *Service) GetTvShowWatchLink(ctx context.Context, serverID string, tvID, season, episode int) (string, error) {
	return u.provider.GetTvShowWatchLink(ctx, serverID, tvID, season, episode)
}
