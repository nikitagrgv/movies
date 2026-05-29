package usecase

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type WatchServiceUsecase struct {
	provider domain.WatchServiceProvider
}

func NewWatchServiceUsecase(provider domain.WatchServiceProvider) *WatchServiceUsecase {
	return &WatchServiceUsecase{provider: provider}
}

func (u *WatchServiceUsecase) GetServices(ctx context.Context) ([]domain.WatchService, error) {
	return u.provider.GetServices(ctx)
}

func (u *WatchServiceUsecase) GetMovieWatchLink(ctx context.Context, serviceID, movieID int) (string, error) {
	return u.provider.GetMovieWatchLink(ctx, serviceID, movieID)
}

func (u *WatchServiceUsecase) GetTvShowWatchLink(ctx context.Context, serviceID, tvID, season, episode int) (string, error) {
	return u.provider.GetTvShowWatchLink(ctx, serviceID, tvID, season, episode)
}
