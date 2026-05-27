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

func (u *SearchMediaUsecase) SearchMovies(ctx context.Context, query string, page int) (domain.SearchResult, error) {
	result, err := u.searcher.SearchMovies(ctx, query, page)
	if err != nil {
		return domain.SearchResult{}, err
	}

	for i := range result.Items {
		result.Items[i] = normalizeMedia(result.Items[i], u.noImageURL)
	}

	return result, nil
}

func (u *SearchMediaUsecase) SearchTvShows(ctx context.Context, query string, page int) (domain.SearchResult, error) {
	result, err := u.searcher.SearchTvShows(ctx, query, page)
	if err != nil {
		return domain.SearchResult{}, err
	}

	for i := range result.Items {
		result.Items[i] = normalizeMedia(result.Items[i], u.noImageURL)
	}

	return result, nil
}
