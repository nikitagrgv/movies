package domain

import "context"

type SearchResult struct {
	Items       []Media
	CurrentPage int
	TotalPages  int
}

type MediaSearcher interface {
	SearchMovies(ctx context.Context, query string, page int) (SearchResult, error)
	SearchTvShows(ctx context.Context, query string, page int) (SearchResult, error)
}
