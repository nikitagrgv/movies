package domain

import "context"

type SearchMoviesResult struct {
	Movies      []Movie
	CurrentPage int
	TotalPages  int
}

type SearchTvShowsResult struct {
	TvShows     []TvShow
	CurrentPage int
	TotalPages  int
}

type MediaSearcher interface {
	SearchMovies(ctx context.Context, query string, page int) (SearchMoviesResult, error)
	SearchTvShows(ctx context.Context, query string, page int) (SearchTvShowsResult, error)
}
