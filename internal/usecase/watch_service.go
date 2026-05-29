package usecase

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type WatchServerUsecase struct {
	provider domain.WatchServerProvider
}

func NewWatchServerUsecase(provider domain.WatchServerProvider) *WatchServerUsecase {
	return &WatchServerUsecase{provider: provider}
}

func (u *WatchServerUsecase) GetServers(ctx context.Context) ([]domain.WatchServer, error) {
	return u.provider.GetServers(ctx)
}

func (u *WatchServerUsecase) GetMovieWatchLink(ctx context.Context, serverID string, movieID int) (string, error) {
	return u.provider.GetMovieWatchLink(ctx, serverID, movieID)
}

func (u *WatchServerUsecase) GetTvShowWatchLink(ctx context.Context, serverID string, tvID, season, episode int) (string, error) {
	return u.provider.GetTvShowWatchLink(ctx, serverID, tvID, season, episode)
}
