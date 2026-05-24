package tmdb

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieSearcher struct{}

func NewMovieSearcher() *MovieSearcher {
	return &MovieSearcher{}
}

func (MovieSearcher) SearchMovie(ctx context.Context, query string, page int) (domain.SearchMovieResult, error) {

	return domain.SearchMovieResult{}, nil
}
