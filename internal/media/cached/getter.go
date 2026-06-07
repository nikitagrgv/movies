package cached

import (
	"context"

	"github.com/nikitagrgv/movies/internal/media"
	"github.com/redis/go-redis/v9"
)

type MediaGetter struct {
	client *redis.Client
	base   media.Getter
}

func NewMediaGetter(client *redis.Client, base media.Getter) *MediaGetter {
	return &MediaGetter{
		client: client,
		base:   base,
	}
}

func (m *MediaGetter) GetMovie(ctx context.Context, id int) (media.Movie, error) {
	return m.base.GetMovie(ctx, id)
}

func (m *MediaGetter) GetTvShow(ctx context.Context, id int) (media.TvShow, error) {
	return m.base.GetTvShow(ctx, id)
}

func (m *MediaGetter) GetTvShowSeason(ctx context.Context, id, season int) (media.Season, error) {
	return m.base.GetTvShowSeason(ctx, id, season)
}
