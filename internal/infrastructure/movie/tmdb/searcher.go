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

func (s *MovieSearcher) SearchMovies(ctx context.Context, query string, page int) (domain.SearchMoviesResult, error) {
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
		return domain.SearchMoviesResult{}, err
	}

	res := domain.SearchMoviesResult{
		CurrentPage: raw.Page,
		TotalPages:  raw.TotalPages,
	}

	for _, m := range raw.Results {
		poster := s.client.getImageURL(m.PosterPath)
		base := domain.MediaBase{
			ID:          m.ID,
			Title:       m.Title,
			Overview:    m.Overview,
			PosterURL:   poster,
			ReleaseDate: m.ReleaseDate,
		}
		res.Movies = append(res.Movies, domain.Movie{Base: base})
	}

	return res, nil
}

func (s *MovieSearcher) SearchTvShows(ctx context.Context, query string, page int) (domain.SearchTvShowsResult, error) {
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
		return domain.SearchTvShowsResult{}, err
	}

	res := domain.SearchTvShowsResult{
		CurrentPage: raw.Page,
		TotalPages:  raw.TotalPages,
	}

	for _, m := range raw.Results {
		poster := s.client.getImageURL(m.PosterPath)
		base := domain.MediaBase{
			ID:          m.ID,
			Title:       m.Name,
			Overview:    m.Overview,
			PosterURL:   poster,
			ReleaseDate: m.FirstAirDate,
		}
		res.TvShows = append(res.TvShows, domain.TvShow{Base: base})
	}

	return res, nil
}
