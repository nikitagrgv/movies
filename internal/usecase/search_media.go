package usecase

import (
	"context"

	"github.com/nikitagrgv/movies/internal/domain"
)

type SearchMediaUsecase struct {
	searcher   domain.MediaSearcher
	noImageURL string
}

func NewSearchMediaUsecase(searcher domain.MediaSearcher, noImageURL string) *SearchMediaUsecase {
	return &SearchMediaUsecase{searcher: searcher, noImageURL: noImageURL}
}

func (u *SearchMediaUsecase) SearchMovies(ctx context.Context, query string, page int) (domain.SearchMoviesResult, error) {
	result, err := u.searcher.SearchMovies(ctx, query, page)
	if err != nil {
		return domain.SearchMoviesResult{}, err
	}

	for i := range result.Movies {
		result.Movies[i].Media = normalizeMedia(result.Movies[i].Media, u.noImageURL)
	}

	return result, nil
}

func (u *SearchMediaUsecase) SearchTvShows(ctx context.Context, query string, page int) (domain.SearchTvShowsResult, error) {
	result, err := u.searcher.SearchTvShows(ctx, query, page)
	if err != nil {
		return domain.SearchTvShowsResult{}, err
	}

	for i := range result.TvShows {
		result.TvShows[i].Media = normalizeMedia(result.TvShows[i].Media, u.noImageURL)
	}

	return result, nil
}
