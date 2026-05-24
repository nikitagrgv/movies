package domain

import "context"

type SearchMoviesResult struct {
	Movies      []Movie
	CurrentPage int
	TotalPages  int
}

type MoviesSearcher interface {
	SearchMovies(ctx context.Context, query string, page int) (SearchMoviesResult, error)
}
