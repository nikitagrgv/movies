package tmdb

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieSearcher struct {
	client *Client
}

func NewMovieSearcher(client *Client) *MovieSearcher {
	return &MovieSearcher{client: client}
}

func (MovieSearcher) SearchMovies(ctx context.Context, query string, page int) (domain.SearchMoviesResult, error) {
	return domain.SearchMoviesResult{}, nil
}
