package tmdb

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieSearcher struct{}

func NewMovieSearcher() *MovieSearcher {
	return &MovieSearcher{}
}

func (MovieSearcher) SearchMovies(ctx context.Context, query string, page int) (domain.SearchMoviesResult, error) {
	return domain.SearchMoviesResult{}, nil
}
