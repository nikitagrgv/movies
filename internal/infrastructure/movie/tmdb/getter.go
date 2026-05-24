package tmdb

import (
	"context"
	"strconv"

	"github.com/nikitagrgv/movies/internal/domain"
)

type MovieGetter struct {
	client *Client
}

func NewMovieGetter(client *Client) *MovieGetter {
	return &MovieGetter{client: client}
}

func (g MovieGetter) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	var raw GetMovieResponse
	err := g.client.get(
		ctx,
		"/search/movie/"+strconv.Itoa(id),
		nil,
		&raw,
	)
	if err != nil {
		return domain.Movie{}, err
	}

	poster := g.client.getImageURL(raw.PosterPath)
	res := domain.Movie{
		ID:          raw.ID,
		Title:       raw.Title,
		Overview:    raw.Overview,
		PosterURL:   poster,
		ReleaseDate: raw.ReleaseDate,
	}

	return res, nil
}
