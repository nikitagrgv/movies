package usecase

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type SearchMoviesParams struct {
	Query string
	Page  int
}

type SearchMovieUsecase struct {
	searcher domain.MovieSearcher
}

func NewSearchMovieUsecase(searcher domain.MovieSearcher) *SearchMovieUsecase {
	return &SearchMovieUsecase{searcher: searcher}
}

func (u *SearchMovieUsecase) SearchMovie(ctx context.Context, query string, page int) (domain.SearchMovieResult, error) {
	return u.searcher.SearchMovie(ctx, query, page)
}
