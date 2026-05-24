package tmdb

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type TMDBMovieSearcher struct{}

func NewTMDBMovieSearcher() *TMDBMovieSearcher {
	return &TMDBMovieSearcher{}
}

func (TMDBMovieSearcher) SearchMovie(ctx context.Context, query string, page int) (domain.SearchMovieResult, error) {
	return domain.SearchMovieResult{}, nil
}
