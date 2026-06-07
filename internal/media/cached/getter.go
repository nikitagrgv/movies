package cached

import (
	"context"
	"strconv"
	"time"

	"github.com/nikitagrgv/movies/internal/media"
	"github.com/nikitagrgv/movies/internal/pkg/cache"
	"github.com/redis/go-redis/v9"
)

const getterTTL = 2 * time.Hour

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
	cacheKey := "media:movie:" + strconv.Itoa(id)

	return cache.GetOrSet(ctx, m.client, cacheKey, getterTTL, func() (media.Movie, error) {
		return m.base.GetMovie(ctx, id)
	})
}

func (m *MediaGetter) GetTvShow(ctx context.Context, id int) (media.TvShow, error) {
	cacheKey := "media:tv:" + strconv.Itoa(id)

	return cache.GetOrSet(ctx, m.client, cacheKey, getterTTL, func() (media.TvShow, error) {
		return m.base.GetTvShow(ctx, id)
	})
}

func (m *MediaGetter) GetTvShowSeason(ctx context.Context, id, season int) (media.Season, error) {
	cacheKey := "media:tv:" + strconv.Itoa(id) + ":season:" + strconv.Itoa(season)

	return cache.GetOrSet(ctx, m.client, cacheKey, getterTTL, func() (media.Season, error) {
		return m.base.GetTvShowSeason(ctx, id, season)
	})

}
