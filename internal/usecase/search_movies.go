package usecase

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type SearchMoviesUsecase struct {
	searcher   domain.MoviesSearcher
	noImageURL string
}

func NewSearchMoviesUsecase(searcher domain.MoviesSearcher, noImageUrl string) *SearchMoviesUsecase {
	return &SearchMoviesUsecase{searcher: searcher, noImageURL: noImageUrl}
}

func (u *SearchMoviesUsecase) SearchMovies(ctx context.Context, query string, page int) (domain.SearchMoviesResult, error) {
	result, err := u.searcher.SearchMovies(ctx, query, page)
	if err != nil {
		return domain.SearchMoviesResult{}, err
	}

	result.Movies = u.normalizeMovies(result.Movies)
	return result, nil
}

func (u *SearchMoviesUsecase) normalizeMovies(movies []domain.Movie) []domain.Movie {
	out := make([]domain.Movie, len(movies))

	for i, m := range movies {
		if m.PosterURL == "" {
			m.PosterURL = u.noImageURL
		}
		out[i] = m
	}

	return out
}
