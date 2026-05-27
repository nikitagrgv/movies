package tmdb

import (
	"context"
	"net/url"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieSearcher struct {
	client *Client
}

func NewMovieSearcher(client *Client) *MovieSearcher {
	return &MovieSearcher{client: client}
}

func (s *MovieSearcher) SearchMovies(ctx context.Context, query string, page int) (domain.SearchResult, error) {
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

func (s *MovieSearcher) SearchTvShows(ctx context.Context, query string, page int) (domain.SearchResult, error) {
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
