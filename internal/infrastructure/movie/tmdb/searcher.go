package tmdb

import (
	"context"
	"net/url"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MediaSearcher struct {
	client *Client
}

func NewMediaSearcher(client *Client) *MediaSearcher {
	return &MediaSearcher{client: client}
}

func (s *MediaSearcher) SearchMovies(ctx context.Context, query string, page int) (domain.SearchResult, error) {
	var raw SearchMovieResponse
	err := s.client.get(
		ctx,
		"/search/movie",
		url.Values{
			"query": {query},
			"page":  {strconv.Itoa(page)},
		},
		&raw,
	)
	if err != nil {
		return domain.SearchResult{}, err
	}

	res := domain.SearchResult{
		CurrentPage: raw.Page,
		TotalPages:  raw.TotalPages,
	}

	for _, m := range raw.Results {
		poster := s.client.getImageURL(m.PosterPath)
		media := domain.Media{
			ID:          m.ID,
			Title:       m.Title,
			Overview:    m.Overview,
			PosterURL:   poster,
			ReleaseYear: parseYear(m.ReleaseDate),
		}
		res.Items = append(res.Items, media)
	}

	return res, nil
}

func (s *MediaSearcher) SearchTvShows(ctx context.Context, query string, page int) (domain.SearchResult, error) {
	var raw SearchTvShowResponse
	err := s.client.get(
		ctx,
		"/search/tv",
		url.Values{
			"query": {query},
			"page":  {strconv.Itoa(page)},
		},
		&raw,
	)
	if err != nil {
		return domain.SearchResult{}, err
	}

	res := domain.SearchResult{
		CurrentPage: raw.Page,
		TotalPages:  raw.TotalPages,
	}

	for _, m := range raw.Results {
		poster := s.client.getImageURL(m.PosterPath)
		media := domain.Media{
			ID:          m.ID,
			Title:       m.Name,
			Overview:    m.Overview,
			PosterURL:   poster,
			ReleaseYear: parseYear(m.FirstAirDate),
		}
		res.Items = append(res.Items, media)
	}

	return res, nil
}
