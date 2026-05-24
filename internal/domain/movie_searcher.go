package domain

import "context"

type SearchMovieResult struct {
	Movies      []Movie
	CurrentPage int
	TotalPages  int
}

type MovieSearcher interface {
	SearchMovie(ctx context.Context, query string, page int) (SearchMovieResult, error)
}
