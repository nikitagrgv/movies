package media

import (
	"context"
)

type SearchResult struct {
	Items       []Media
	CurrentPage int
	TotalPages  int
}

type Searcher interface {
	SearchMovies(ctx context.Context, query string, page int) (SearchResult, error)
	SearchTvShows(ctx context.Context, query string, page int) (SearchResult, error)
}

type Getter interface {
	GetMovie(ctx context.Context, id int) (Movie, error)
	GetTvShow(ctx context.Context, id int) (TvShow, error)
	GetTvShowSeason(ctx context.Context, id, season int) (Season, error)
}

type Service struct {
	getter     Getter
	searcher   Searcher
	noImageURL string
}

func NewService(getter Getter, searcher Searcher, noImageURL string) *Service {
	return &Service{getter: getter, searcher: searcher, noImageURL: noImageURL}
}

func (u *Service) GetMovie(ctx context.Context, id int) (Movie, error) {
	media, err := u.getter.GetMovie(ctx, id)
	if err != nil {
		return Movie{}, err
	}

	media.Media = normalizeMedia(media.Media, u.noImageURL)
	return media, nil
}

func (u *Service) GetTvShow(ctx context.Context, id int) (TvShow, error) {
	media, err := u.getter.GetTvShow(ctx, id)
	if err != nil {
		return TvShow{}, err
	}

	media.Media = normalizeMedia(media.Media, u.noImageURL)
	return media, nil
}

func (u *Service) GetTvShowSeason(ctx context.Context, id, season int) (Season, error) {
	s, err := u.getter.GetTvShowSeason(ctx, id, season)
	if err != nil {
		return Season{}, err
	}

	return s, nil
}

func (u *Service) SearchMovies(ctx context.Context, query string, page int) (SearchResult, error) {
	result, err := u.searcher.SearchMovies(ctx, query, page)
	if err != nil {
		return SearchResult{}, err
	}

	for i := range result.Items {
		result.Items[i] = normalizeMedia(result.Items[i], u.noImageURL)
	}

	return result, nil
}

func (u *Service) SearchTvShows(ctx context.Context, query string, page int) (SearchResult, error) {
	result, err := u.searcher.SearchTvShows(ctx, query, page)
	if err != nil {
		return SearchResult{}, err
	}

	for i := range result.Items {
		result.Items[i] = normalizeMedia(result.Items[i], u.noImageURL)
	}

	return result, nil
}

// TODO: Move out of here?
func normalizeMedia(media Media, noImageURL string) Media {
	if media.PosterURL == "" {
		media.PosterURL = noImageURL
	}

	return media
}
