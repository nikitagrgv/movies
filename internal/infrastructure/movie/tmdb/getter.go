package tmdb

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieGetter struct {
	client *Client
}

func NewMovieGetter(client *Client) *MovieGetter {
	return &MovieGetter{client: client}
}

func (g MovieGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	return domain.Movie{}, nil
}
