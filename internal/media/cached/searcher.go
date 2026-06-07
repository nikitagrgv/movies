package cached

import (
	"context"

	"github.com/nikitagrgv/movies/internal/media"
	"github.com/redis/go-redis/v9"
)

type MediaSearcher struct {
	client *redis.Client
	base   media.Searcher
}

func NewMediaSearcher(client *redis.Client, base media.Searcher) *MediaSearcher {
	return &MediaSearcher{
		client: client,
		base:   base,
	}
}

func (m *MediaSearcher) SearchMovies(ctx context.Context, query string, page int) (media.SearchResult, error) {
	return m.base.SearchMovies(ctx, query, page)
}

func (m *MediaSearcher) SearchTvShows(ctx context.Context, query string, page int) (media.SearchResult, error) {
	return m.base.SearchTvShows(ctx, query, page)
}
