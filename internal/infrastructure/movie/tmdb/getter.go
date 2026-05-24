package tmdb

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieGetter struct{}

func NewMovieGetter() *MovieGetter {
	return &MovieGetter{}
}

func (g MovieGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	return domain.Movie{}, nil
}
