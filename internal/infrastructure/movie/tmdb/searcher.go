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
		res.Movies = append(res.Movies, domain.Movie{
			ID:          m.ID,
			Title:       m.Title,
			Overview:    m.Overview,
			PosterURL:   poster,
			ReleaseDate: m.FirstAirDate,
		})
	}

	return res, nil
}
